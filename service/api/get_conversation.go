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
func (rt *_router) GetConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	convId, err := strconv.Atoi(ps.ByName("conv_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is in the conversation
	check, err := rt.db.IsUserInConv(UserId, convId)
	if err != nil {
		http.Error(w, "Internal server error A"+err.Error(), http.StatusInternalServerError)
		return
	}
	if !check {
		http.Error(w, "user ins't in this conversation", http.StatusBadRequest)
		return
	}

	// Get messages
	messages, err := rt.db.GetMessagesByConvId(convId)
	if err != nil {
		http.Error(w, "Internal server error B"+err.Error(), http.StatusInternalServerError)
		return
	}

	// sort messages by message id
	for i := 0; i < len(messages); i++ {
		for j := 0; j < len(messages); j++ {
			if messages[i].MessageId < messages[j].MessageId {
				messages[i], messages[j] = messages[j], messages[i]
			}
		}
	}

	// Update the list of users that have read the messages
	for i := 0; i < len(messages); i++ {
		if !messages[i].Status {
			// Control if the user has already read the message
			check, err := rt.db.CheckIfUserHasReadMess(messages[i].MessageId, UserId)
			if err != nil {
				http.Error(w, "Internal server error C"+err.Error(), http.StatusInternalServerError)
				return
			}
			if check {
				continue
			}
			// Add the user to the list of users that have read the message
			err = rt.db.AddUserToListOfAlreadyReadersOfMess(messages[i].MessageId, UserId, convId)
			if err != nil {
				http.Error(w, "Internal server error D"+err.Error(), http.StatusInternalServerError)
				return
			}
			// Check if all the users have read the message
			check, err = rt.db.CheckAllUsersHaveReadMess(messages[i].MessageId)
			if err != nil {
				http.Error(w, "Internal server error E"+err.Error(), http.StatusInternalServerError)
				return
			}
			if check {
				// Update the message status
				err = rt.db.UpdateMessageStatus(messages[i].MessageId)
				if err != nil {
					http.Error(w, "Internal server error F"+err.Error(), http.StatusInternalServerError)
					return
				}
				// Get the updated message
				messages[i].Status = true
			}
		}
	}

	// Get only the last 50 messages
	if len(messages) > 50 {
		messages = messages[len(messages)-50:]
	}

	type MessageData struct {
		Message structions.Message `json:"message"`
		Sender structions.User    `json:"sender"`
		DateTime string           `json:"dateTime"`
		Comments []structions.Comment `json:"comments"`
	}

	var response []MessageData
	for i := 0; i < len(messages); i++ {
		// Get the sender
		sender, err := rt.db.GetUserById(messages[i].SenderId)
		if err != nil {
			http.Error(w, "Internal server error G"+err.Error(), http.StatusInternalServerError)
			return
		}
		// Get the comments
		comments, err := rt.db.GetCommentsByMessId(messages[i].MessageId, convId)
		if err != nil {
			http.Error(w, "Internal server error H"+err.Error(), http.StatusInternalServerError)
			return
		}
		// Get the date and time
		dateTime := messages[i].DateTime.Format("15:04 - 02/01/2006")

		response = append(response, MessageData{messages[i], sender, dateTime, comments})
	}

	// Sort the messages by message dateTime
	sort.Slice(response, func(i, j int) bool {
		return response[i].Message.DateTime.After(response[j].Message.DateTime)
	})

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Error in encoding the response")
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}