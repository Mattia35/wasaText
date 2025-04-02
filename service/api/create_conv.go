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
		BadRequest(w, err, "Bad Request: ")
		return
	}

	userID := ctx.UserID

	// Check if the user is authorized
	if UserId != userID {
		Forbidden(w, err, "Forbidden: ")
		return
	}

	type RequestBody struct {
		User string `json:"user"`
		Text string `json:"text"`
	}
	var request RequestBody

	// Take the message to sent from the Request
	var message structions.Message
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the destination user exists
	Dest, err := rt.db.GetUserByName(request.User)
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Get the destination user id
	DestId := Dest.UserId

	// Check if the text in the request is empty
	if request.Text == "" {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the conversation already exists
	convNumb, err := rt.db.GetConvByUsers(UserId, DestId)
	if convNumb != 1 || err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Create the conversation
	var conversation structions.Conversation
	conversation.GroupId = 0
	conversation, err = rt.db.CreateConversation(conversation)
	if err != nil {
		BadRequest(w, err, "Can't create the conversation of the private chat: ")
		return
	}

	message.Text = request.Text
	message.SenderId = UserId
	message.Status = false
	message.ConvId = conversation.ConvId

	// Create the message
	message, err = rt.db.CreateMessage(message)
	if err != nil {
		BadRequest(w, err, "Can't create the welcome message: ")
		return
	}

	// Update last message of a conversation
	err = rt.db.AddMessageToConv(message.MessageId, conversation.ConvId)
	if err != nil {
		InternalServerError(w, err, "Can't update the last message of the conversation: ")
		return
	}
	conversation.LastMessage = message.MessageId
	// Add the user to the conversation
	err = rt.db.AddUserToConv(UserId, conversation.ConvId)
	if err != nil {
		InternalServerError(w, err, "You can't add a user to the conversation: ")
		return
	}
	// Add the destination user to the conversation
	err = rt.db.AddUserToConv(DestId, conversation.ConvId)
	if err != nil {
		InternalServerError(w, err, "You can't add a user to the conversation: ")
		return
	}

	// Add the destination user to the list of readers of the message
	err = rt.db.AddUserToListOfReadersOfMess(message.MessageId, DestId, conversation.ConvId)
	if err != nil {
		InternalServerError(w, err, "Error adding the user to the list of readers of the message: ")
		return
	}

	w.WriteHeader(http.StatusCreated)
	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(conversation); err != nil {
		InternalServerError(w, err, "Error encoding the conversation: ")
		return
	}

}
