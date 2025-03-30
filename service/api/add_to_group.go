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
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	if !check {
		http.Error(w, "User already is in the group", http.StatusBadRequest)
		return
	}

	type RequestBody struct {
		Users []structions.User `json:"users"`
	}
	var request RequestBody

	// Decode the request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	// Add users to the group and conversation
	for i := 0; i < len(request.Users); i++ {
		user, err := rt.db.UserControlByUsername(request.Users[i].Username)
		if err != nil {
			ctx.Logger.WithError(err).Error("you can't add a user to the group, because it doesn't exist")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		// Check if the user isn't in the group
		check, err = rt.db.IsUserInGroup(user.UserId, groupId)

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
		if check {
			http.Error(w, "you can't add a user to the group, because it already is in", http.StatusBadRequest)
			return
		}
		// Add the user to the group
		err = rt.db.AddUserToGroup(user.UserId, groupId)
		if err != nil {
			ctx.Logger.WithError(err).Error("you can't add a user to the group")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Select the conversation of the group
		conversation, err := rt.db.GetConvByGroupId(groupId)
		if err != nil {
			ctx.Logger.WithError(err).Error("you can't add a user to the conversation of the group")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		// Add the user to the conversation
		err = rt.db.AddUserToConv(user.UserId, conversation.ConvId)
		if err != nil {
			ctx.Logger.WithError(err).Error("you can't add a user to the conversation of the group")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	// Response (users of the group)
	membersList, err := rt.db.GetUsersByGroupId(groupId)
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
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
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
}
