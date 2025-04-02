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
		BadRequest(w, err, "Bad Request: ")
		return
	}

	userID := ctx.UserID

	// Check if the user is authorized
	if UserId != userID {
		Forbidden(w, err, "Forbidden: ")
		return
	}

	// Get the group id
	groupId, err := strconv.Atoi(ps.ByName("group_id"))
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the user is in the group
	check, err := rt.db.IsUserInGroup(UserId, groupId)

	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	if !check {
		http.Error(w, "User isn't in the group", http.StatusBadRequest)
		return
	}

	// Remove the user from the group
	err = rt.db.RemoveUserFromGroup(UserId, groupId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	// Select the conversation of the group
	_, conv, err := rt.db.GetConvByGroupId(groupId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	convId := conv.ConvId
	// Remove the user from the conversation
	err = rt.db.RemoveUserFromConv(UserId, convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}

	// control if the user is the last user in the group
	users, err := rt.db.GetUsersByGroupId(groupId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	if len(users) == 0 {
		// remove the group
		err = rt.db.RemoveGroup(groupId)
		if err != nil {
			InternalServerError(w, err, "Internal server error: ")
			return
		}
		// remove the conversation
		err = rt.db.RemoveConv(convId)
		if err != nil {
			InternalServerError(w, err, "Internal server error: ")
			return
		}
	}

	// user has been removed from group, response 200
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode("user has been successfully removed from group!"); err != nil {
		InternalServerError(w, err, "Error encoding response: ")
		return
	}
}
