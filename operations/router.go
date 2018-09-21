package operations

import (
	"github.com/arunvm/twitter-clone/pkg/mysql"
	"github.com/gorilla/mux"
)

func NewRoutes(db *mysql.MySQL) *mux.Router {
	m := mux.NewRouter()
	m.Handle("/signup", userSignup(db)).Methods("POST")
	m.Handle("/login", userLogin(db)).Methods("POST")
	m.Handle("/graphql", graphqlHandler(db))

	return m
}
