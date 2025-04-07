package api

import (
	"encoding/json"
	"net/http"
	"sort"
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

	// Get the conversations of the user
	conversations, err := rt.db.GetConversationsByUserId(UserId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}

	// Struct used for the response
	type ConvObject struct {
		Conversation structions.Conversation `json:"conversation"`
		User         structions.User         `json:"user"`
		Group        structions.Group        `json:"group"`
		GroupUsers   []structions.User       `json:"groupUsers"`
		Message      structions.Message      `json:"message"`
		SenderUser   structions.User         `json:"senderUser"`
		DateTime     string                  `json:"dateTime"`
	}

	// Response
	response := make([]ConvObject, len(conversations))

	// Get information about the conversations
	for idx, conv := range conversations {
		if conv.GroupId == 0 {
			// Get the user from the conversation
			destUser, err := rt.db.GetOtherUserByConv(conv.ConvId, UserId)
			if err != nil {
				BadRequest(w, err, "Bad Request: ")
				return
			}
			// Get the user
			user, err := rt.db.GetUserById(destUser.UserId)
			if err != nil {
				BadRequest(w, err, "Bad Request: ")
				return
			}

			lastMessage := structions.Message{}
			var dateTime string
			senderUser := structions.User{}
			if conv.LastMessage != 0 {
				// Get last message
				lastMessage, err = rt.db.GetMessageById(conv.LastMessage, conv.ConvId)
				if err != nil {
					BadRequest(w, err, "Bad Request: ")
					return
				}
				// Get the dateTime of the last message
				dateTime = lastMessage.DateTime.Format("15:04 - 02/01/2006")
				// Get the sender of the last message
				senderUser, err = rt.db.GetUserById(lastMessage.SenderId)
				if err != nil {
					BadRequest(w, err, "Bad Request: ")
					return
				}
			}

			response[idx] = ConvObject{
				Conversation: conv,
				User:         user,
				Message:      lastMessage,
				SenderUser:   senderUser,
				DateTime:     dateTime,
			}
		} else {
			// Get the group from the conversation
			group, err := rt.db.GetGroupByGroupId(conv.GroupId)
			if err != nil {
				BadRequest(w, err, "Bad Request: ")
				return
			}
			lastMessage := structions.Message{}
			var dateTime string
			senderUser := structions.User{}
			if conv.LastMessage != 0 {
				// Get last message
				lastMessage, err = rt.db.GetMessageById(conv.LastMessage, conv.ConvId)
				if err != nil {
					BadRequest(w, err, "Bad Request: ")
					return
				}
				// Get the dateTime of the last message
				dateTime = lastMessage.DateTime.Format("15:04 - 02/01/2006")
				// Get the sender of the last message
				senderUser, err = rt.db.GetUserById(lastMessage.SenderId)
				if err != nil {
					BadRequest(w, err, "Bad Request: ")
					return
				}
			}
			// Get the users of the group
			users, err := rt.db.GetUsersByGroupId(conv.GroupId)
			if err != nil {
				BadRequest(w, err, "Bad Request: ")
				return
			}
			// Delete yourself from the list of users
			for i, user := range users {
				if user.UserId == UserId {
					users = append(users[:i], users[i+1:]...)
					break
				}
			}

			response[idx] = ConvObject{
				Conversation: conv,
				Group:        group,
				GroupUsers:   users,
				Message:      lastMessage,
				SenderUser:   senderUser,
				DateTime:     dateTime,
			}
		}
	}
	// Sort the conversations by the last message dateTime
	sort.Slice(response, func(i, j int) bool {
		return response[i].Message.DateTime.After(response[j].Message.DateTime)
	})

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerError(w, err, "Error encoding the conversation: ")
		return
	}

}
