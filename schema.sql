create database twitterClone;
use twitterClone;
create table users (
    Name varchar(50) NOT NULL,
    UserName varchar(50) NOT NULL,
    Password varchar(50) NOT NULL,
    PRIMARY KEY(UserName)
);
CREATE TABLE posts (
    PostID int(11) NOT NULL,
    Text varchar(200) NOT NULL,
    UserName varchar(50) NOT NULL,
    TimeStamp timestamp NOT NULL,
    PRIMARY KEY(PostID),
    FOREIGN KEY (UserName) REFERENCES users(UserName)
);
CREATE TABLE follows (
    UserName varchar(50) NOT NULL,
    FollowsName varchar(50) NOT NULL,
    PRIMARY KEY (UserName,FollowsName),
    FOREIGN KEY (UserName) REFERENCES users(UserName),
    FOREIGN KEY (FollowsName) REFERENCES users(UserName)
);