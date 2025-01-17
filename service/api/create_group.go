package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"progetto.wasa/service/api/structions"
)

func (rt *_router) CreateGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	var group Group
	type RequestBody struct {
		groupname string `json:"groupname"`
		users     []User `json:"users"`
	}
	var request RequestBody
	// Check if the user makes a bad request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	group.Username = request.groupname
	// Check if the group respects the regex, so there is a bad request
	if !group.IsValid() {
		http.Error(w, "Invalid groupname", http.StatusBadRequest)
		return
	}
	var conversation structions.Conversation
	conversation.group = true

	groupDB, err := rt.db.CreateGroup(group.GroupToDatabase(), UserId, conversation.convId)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't create the group")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)

}
