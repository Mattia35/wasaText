package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"progetto.wasa/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) UsernameModify (w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Check if the user request is valid
	UserId, err := strconv.Atoi(ps.ByName("user"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	userID := ctx.UserID

	// Check if the user is authorized
	if UserId != userID {
		http.Error(w, "Forbidden", http.StatusForbidden)
		return
	}

	var user User
	// Check if the user makes a bad request
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user respects the regex, so there is a bad request
	if !user.IsValid() {
		http.Error(w, "Invalid username", http.StatusBadRequest)
		return
	}
	user.UserId = UserId
	// Check if the server has internal problems
	if err := user.FromDatabase(user.ToDatabase()); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	// Try to modify the username. If username is already taken, it gives an error
	if err := rt.db.UsernameModify(userID, user.Username); err!= nil {
		http.Error(w, "Username already taken. Retry!", http.StatusBadRequest)
		return
	}

	// Username changed, resposne 200
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}