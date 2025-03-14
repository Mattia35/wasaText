package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
)

func (rt *_router) LeaveGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Get the group id
	groupId, err := strconv.Atoi(ps.ByName("group_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is in the group
	check, err := rt.db.IsUserInGroup(UserId, groupId)

	if err != nil {
		http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
		return
	}
	if !check {
		http.Error(w, "User isn't in the group", http.StatusBadRequest)
		return
	}

	// Remove the user from the group
	err = rt.db.RemoveUserFromGroup(UserId, groupId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	// Select the conversation of the group
	conv, err := rt.db.GetConvByGroupId(groupId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	convId := conv.ConvId
	// Remove the user from the conversation
	err = rt.db.RemoveUserFromConv(UserId, convId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}

	// control if the user is the last user in the group
	users, err := rt.db.GetUsersByGroupId(groupId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	if len(users) == 0 {
		// remove the group
		err = rt.db.RemoveGroup(groupId)
		if err != nil {
			http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
			return
		}
		// remove the conversation
		err = rt.db.RemoveConv(convId)
		if err != nil {
			http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	// user has been removed from group, response 200
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode("user has been successfully removed from group!"); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
