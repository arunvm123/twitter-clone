package main

import (
	"database/sql"
	"log"
	"net/http"

	interpose "github.com/carbocation/interpose/middleware"

	"github.com/arunvm/twitter-clone/operations"
	"github.com/arunvm/twitter-clone/pkg/mysql"
)

func main() {
	var db mysql.MySQL

	// Setting up MySQL 
	dbconn, err := sql.Open("mysql", "root:root@tcp(db:3306)/twitterClone?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		log.Fatal(err.Error())
	}
	defer dbconn.Close()

	db.Con = dbconn

	routes := operations.NewRoutes(&db)

	log.Println("Server starting at port :8080")
	logViaLogrus := interpose.NegroniLogrus()
	log.Fatal(http.ListenAndServe(":8080", logViaLogrus(routes)))
}
