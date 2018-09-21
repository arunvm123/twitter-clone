package operations

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/arunvm/cleverhires/API/globals"
	"github.com/arunvm/twitter-clone"

	"github.com/arunvm/twitter-clone/pkg/auth"
	"github.com/arunvm/twitter-clone/pkg/mysql"
)

func userSignup(db *mysql.MySQL) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var user twitterclone.User

			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				log.Println(err)
				respondWithError(w, http.StatusInternalServerError, "Error decoding request data")
				return
			}

			err = db.AddUser(&user)
			if err != nil {
				log.Println(err)
				respondWithError(w, http.StatusInternalServerError, "Error adding user to db")
				return
			}

			respondWithJSON(w, http.StatusOK, "Successful Signup")
			return
		},
	)
}

func userLogin(db *mysql.MySQL) http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			var user twitterclone.User

			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				log.Println(err)
				respondWithError(w, http.StatusInternalServerError, "Error decoding request data")
				return
			}

			_, err = db.GetUser(user.UserName, user.Password)
			if err != nil {
				log.Println(err)
				respondWithError(w, http.StatusInternalServerError, "Error fetching user from db")
				return
			}

			t, err := auth.JWTTokenGeneration(user.UserName)
			if err != nil {
				log.Println(err)
				respondWithError(w, http.StatusInternalServerError, "Error generating token")
				return
			}

			respondWithJSON(w, http.StatusOK, t)
			return

		},
	)
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, err := json.Marshal(payload)
	if err != nil {
		log.Println(err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
	return
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, globals.ErrorResponse{
		Error: message,
	})
	return
}
