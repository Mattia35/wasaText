package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"progetto.wasa/service/api/structions"
)

func (rt *_router) CreateConv(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Check if the destination user request is valid
	DestId, err := strconv.Atoi(ps.ByName("conv_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the destination user exists
	_, err = rt.db.GetUserById(DestId)
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the conversation already exists
	convNumb, err := rt.db.GetConvByUsers(UserId, DestId)
	if convNumb != 1 || err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Create the conversation
	var conversation structions.Conversation
	conversation.GroupId = 0
	conversation, err = rt.db.CreateConversation(conversation)
	if err != nil {
		ctx.Logger.WithError(err).Error("can't create the conversation of the private chat")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	type RequestBody struct {
		Text string `json:"text"`
	}
	var request RequestBody
	// Take the message to sent from the Request
	var message structions.Message
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	message.Text = request.Text
	message.SenderId = UserId
	message.Status = false
	message.ConvId = conversation.ConvId

	// Create the message
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
	conversation.LastMessage = message.MessageId
	// Add the user to the conversation
	err = rt.db.AddUserToConv(UserId, conversation.ConvId)
	if err != nil {
		ctx.Logger.WithError(err).Error("you can't add a user to the conversation")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// Add the destination user to the conversation
	err = rt.db.AddUserToConv(DestId, conversation.ConvId)
	if err != nil {
		ctx.Logger.WithError(err).Error("you can't add a user to the conversation")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		ctx.Logger.WithError(err).Error("Error in encoding the conversation")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}
