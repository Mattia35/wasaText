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

func (rt *_router) GetConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

			message := structions.Message{}
			var dateTime string
			senderUser := structions.User{}
			if conv.LastMessage != 0 {
				// Get last message
				message, err = rt.db.GetMessageById(conv.LastMessage, conv.ConvId)
				if err != nil {
					BadRequest(w, err, "Bad Request: ")
					return
				}
				// Get the dateTime of the last message
				dateTime = message.DateTime.Format("15:04 - 02/01/2006")
				// Get the sender of the last message
				senderUser, err = rt.db.GetUserById(message.SenderId)
				if err != nil {
					BadRequest(w, err, "Bad Request: ")
					return
				}
			}

			response[idx] = ConvObject{
				Conversation: conv,
				User:         user,
				Message:      message,
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
			message := structions.Message{}
			var dateTime string
			senderUser := structions.User{}
			if conv.LastMessage != 0 {
				// Get last message
				message, err = rt.db.GetMessageById(conv.LastMessage, conv.ConvId)
				if err != nil {
					BadRequest(w, err, "Bad Request: ")
					return
				}
				// Get the dateTime of the last message
				dateTime = message.DateTime.Format("15:04 - 02/01/2006")
				// Get the sender of the last message
				senderUser, err = rt.db.GetUserById(message.SenderId)
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
				Message:      message,
				SenderUser:   senderUser,
				DateTime:     dateTime,
			}
		}
	}
	// Sort the conversations by the last message dateTime
	sort.Slice(response, func(i, j int) bool {
		return response[i].Message.DateTime.After(response[j].Message.DateTime)
	})

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		InternalServerError(w, err, "Error encoding response: ")
		return
	}

}
