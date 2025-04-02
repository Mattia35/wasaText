package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"sort"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
)

func (rt *_router) SearchUsers(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Get the search query
	query := r.URL.Query().Get("query")
	if query == "" {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}
	// control if the query is valid
	validUser := regexp.MustCompile(`^[a-z0-9]{1,15}$`)
	check := validUser.MatchString(query)
	if !check {
		BadRequest(w, nil, "Bad Request: ")
		return
	}
	// Get the users
	users, err := rt.db.SearchUsers(query)
	if err != nil {
		InternalServerError(w, err, "Internal Server Error: ")
		return
	}
	if len(users) == 0 {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}
	// sort the users
	sort.Slice(users, func(i, j int) bool {
		return users[i].Username < users[j].Username
	})
	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		InternalServerError(w, err, "Error in encoding the response: ")
		return
	}
}
