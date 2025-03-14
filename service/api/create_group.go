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

	var group structions.Group
	type RequestBody struct {
		Groupname string            `json:"groupname"`
		Users     []structions.User `json:"users"`
	}
	var request RequestBody
	// Check if the user makes a bad request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	group.Username = request.Groupname
	// Check if the group respects the regex, so there is a bad request
	if !group.IsValid() {
		http.Error(w, "Invalid groupname", http.StatusBadRequest)
		return
	}
	var conversation structions.Conversation

	// Create the group
	group, err = rt.db.CreateGroup(group, UserId)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't create the group")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	conversation.GroupId = group.GroupId
	// Create the group conversation
	conversation, err = rt.db.CreateConversation(conversation)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't create the conversation of the group")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Add the user who is creating the group to the group and conversation
	err = rt.db.AddUserToConv(UserId, conversation.ConvId)
	if err != nil {
		ctx.Logger.WithError(err).Error("you can't add a user to the group")
		w.WriteHeader(http.StatusInternalServerError)
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
		err = rt.db.AddUserToGroup(user.UserId, group.GroupId)
		if err != nil {
			ctx.Logger.WithError(err).Error("you can't add a user to the group")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		err = rt.db.AddUserToConv(user.UserId, conversation.ConvId)
		if err != nil {
			ctx.Logger.WithError(err).Error("you can't add a user to the conversation of the group")
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	type Response struct {
		Group          structions.Group `json:"group"`
		ConversationId int              `json:"conversationId"`
	}

	var response Response
	response.Group = group
	response.ConversationId = conversation.ConvId

	message := structions.Message{
		SenderId: UserId,
		ConvId:   conversation.ConvId,
		Text:     "You are now part of the group " + group.Username,
	}

	// Create the welcome message
	message, err = rt.db.CreateMessage(message)
	if err != nil {
		ctx.Logger.WithError(err).Error("server now can't create the welcome message")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Update last message of a conversation
	err = rt.db.AddMessageToConv(message.MessageId, conversation.ConvId)
	if err != nil {
		ctx.Logger.WithError(err).Error("server now can't update the last message of the conversation")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// get users of the group
	users, err := rt.db.GetUsersByGroupId(group.GroupId)
	if err != nil {
		http.Error(w, "Error taking the users of the group"+err.Error(), http.StatusBadRequest)
		return
	}

	// Set the users that have read the message: all the users of the group, unless the sender
	newUsers := make([]structions.User, 0)
	for i := 0; i < len(users); i++ {
		if users[i].UserId != UserId {
			newUsers = append(newUsers, users[i])
		}
	}

	// Add the users to the list of readers of the message
	users = newUsers
	for i := 0; i < len(users); i++ {
		err = rt.db.AddUserToListOfReadersOfMess(message.MessageId, users[i].UserId, conversation.ConvId)
		if err != nil {
			http.Error(w, "Error adding the user to the list of readers of the message"+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Error in encoding the response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
