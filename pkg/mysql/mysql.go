package mysql

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/arunvm/twitter-clone"
	_ "github.com/go-sql-driver/mysql"
)

var _ twitterclone.TwitterCloneDB = &MySQL{}

type MySQL struct {
	Con *sql.DB
}

func (db *MySQL) AddUser(user *twitterclone.User) error {
	stmt, err := db.Con.Prepare("INSERT INTO users (Name,UserName,Password) values (?,?,?)")
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not prepare sql statement: %v", err)
	}

	_, err = stmt.Exec(&user.Name, &user.UserName, &user.Password)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not execute sql statement: %v", err)
	}

	return nil
}

func (db *MySQL) GetUser(username, password string) (*twitterclone.User, error) {
	var user twitterclone.User
	stmt, err := db.Con.Prepare("SELECT * FROM users where UserName = ? AND Password = ?")
	if err != nil {
		return nil, err
	}

	res := stmt.QueryRow(username, password)
	err = res.Scan(&user.Name, &user.UserName, &user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (db *MySQL) AddPost(userName string, text string) (*twitterclone.Post, error) {

	stmt, err := db.Con.Prepare("INSERT INTO posts (PostID,Text,UserName,TimeStamp) values (?,?,?,?)")
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Could not prepare sql statement: %v", err)
	}

	rand.Seed(time.Now().UTC().UnixNano())
	post := &twitterclone.Post{
		PostID:    rand.Intn(10000),
		Text:      text,
		UserName:  userName,
		TimeStamp: time.Now().UTC(),
	}

	_, err = stmt.Exec(&post.PostID, &post.Text, &post.UserName, &post.TimeStamp)
	if err != nil {
		log.Println(err)
		return nil, fmt.Errorf("Could not execute sql statement: %v", err)
	}

	return post, nil
}

func (db *MySQL) GetAllPostsOfUser(userName string) (*[]twitterclone.Post, error) {
	var post twitterclone.Post
	var posts []twitterclone.Post

	results, err := db.Con.Query("SELECT * FROM posts where UserName = ?", userName)
	if err != nil {
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&post.PostID, &post.Text, &post.UserName, &post.TimeStamp)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Error when reading from db: %v", err)
		}

		posts = append(posts, post)
	}

	return &posts, nil
}

func (db *MySQL) GetAllUsers(userName string) (*[]twitterclone.User, error) {
	var user twitterclone.User
	var users []twitterclone.User

	results, err := db.Con.Query("SELECT Name,UserName FROM users where UserName != ?", userName)
	if err != nil {
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&user.Name, &user.UserName)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Error when reading from db: %v", err)
		}

		users = append(users, user)
	}

	return &users, nil
}

func (db *MySQL) FollowUser(userName, userNameToFollow string) error {
	stmt, err := db.Con.Prepare("INSERT INTO follows (UserName,FollowsName) values (?,?)")
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not prepare sql statement: %v", err)
	}

	_, err = stmt.Exec(userName, userNameToFollow)
	if err != nil {
		log.Println(err)
		return fmt.Errorf("Could not execute sql statement: %v", err)
	}

	return nil
}

func (db *MySQL) GetPostsFeed(userName string) (*[]twitterclone.Post, error) {
	var post twitterclone.Post
	var posts []twitterclone.Post

	results, err := db.Con.Query("select * from posts where UserName in (select FollowsName from follows where UserName = ?) order by TimeStamp;", userName)
	if err != nil {
		return nil, err
	}

	for results.Next() {
		err = results.Scan(&post.PostID, &post.Text, &post.UserName, &post.TimeStamp)
		if err != nil {
			log.Println(err)
			return nil, fmt.Errorf("Error when reading from db: %v", err)
		}

		posts = append(posts, post)
	}

	return &posts, nil
}
