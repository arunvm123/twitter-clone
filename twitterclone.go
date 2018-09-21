package twitterclone

import (
	"time"
)

type TwitterCloneDB interface {
	AddUser(user *User) error
	GetUser(username string, password string) (*User, error)
	AddPost(userName string, text string) (*Post, error)
	GetAllPostsOfUser(username string) (*[]Post, error)
	GetAllUsers(username string) (*[]User, error)
	FollowUser(username string, usernameToFollow string) error
	GetPostsFeed(username string) (*[]Post, error)
}

type User struct {
	Name     string `json:"name"`
	UserName string `json:"user_name"`
	Password string `json:"password"`
}

type Post struct {
	PostID    int       `json:"post_id"`
	Text      string    `json:"text"`
	UserName  string    `json:"user_name"`
	TimeStamp time.Time `json:"time_stamp"`
}
