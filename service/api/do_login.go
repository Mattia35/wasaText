package api

import (
	"encoding/json"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
)

func (rt *_router) doLogin(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	var user User

	// Read the request body with json
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Bad request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the username is valid
	if !user.IsValid() {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}

	/*
		Check if the user exists
		If the user does not exist, create a new user
		else return the user object
	*/
	nameExistance, err := rt.db.NameControl(user.Username)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't check if the user exists")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if !nameExistance {
		user, err = rt.CreateUser(user)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't create the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	} else {
		userInDb, err := rt.db.GetUserByName(user.Username)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't load the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = user.FromDatabase(userInDb)
		if err != nil {
			ctx.Logger.WithError(err).Error("can't convert the user")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}

	// This struct contain the User object and the authorization token.
	type UserAuthentication struct {
		User  User `json:"user"`
		Token int  `json:"token"`
	}

	/*
		Create the AuthUser object and set the user object and the authorization token.
		The authorization token is the userID.
	*/
	authUser := UserAuthentication{user, user.UserId}

	// Encode the AuthUser object in JSON and send it to the client.
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(authUser); err != nil {
		ctx.Logger.WithError(err).Error("can't encode the response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
