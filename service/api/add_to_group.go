package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"progetto.wasa/service/api/structions"
)

func (rt *_router) AddToGroup(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
		BadRequest(w, err, "User isn't in the group: ")
		return
	}

	type RequestBody struct {
		Users []structions.User `json:"users"`
	}
	var request RequestBody

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}
	// Add users to the group and conversation
	for i := 0; i < len(request.Users); i++ {
		user, err := rt.db.UserControlByUsername(request.Users[i].Username)
		if err != nil {
			BadRequest(w, err, "You can't add a user to the group, because it doesn't exist: ")
			return
		}
		// Check if the user isn't in the group
		check, err = rt.db.IsUserInGroup(user.UserId, groupId)

		if err != nil {
			InternalServerError(w, err, "Internal server error: ")
			return
		}
		if check {
			BadRequest(w, err, "You can't add a user to the group, because it already is in: ")
			return
		}
		// Add the user to the group
		err = rt.db.AddUserToGroup(user.UserId, groupId)
		if err != nil {
			BadRequest(w, err, "You can't add a user to the group: ")
			return
		}
		// Select the conversation of the group
		_, conversation, err := rt.db.GetConvByGroupId(groupId)
		if err != nil {
			BadRequest(w, err, "You can't add a user to the conversation of the group: ")
			return
		}
		// Add the user to the conversation
		err = rt.db.AddUserToConv(user.UserId, conversation.ConvId)
		if err != nil {
			BadRequest(w, err, "You can't add a user to the conversation of the group: ")
			return
		}
	}

	// Response (users of the group)
	membersList, err := rt.db.GetUsersByGroupId(groupId)
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Delete the user from the members list
	for i := 0; i < len(membersList); i++ {
		if membersList[i].UserId == UserId {
			membersList = append(membersList[:i], membersList[i+1:]...)
			break
		}
	}

	// users has been added to group, response 200
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(membersList); err != nil {
		InternalServerError(w, err, "Error encoding response: ")
		return
	}
}
