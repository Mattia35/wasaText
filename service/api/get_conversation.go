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
	convId, err := strconv.Atoi(ps.ByName("conv_id"))
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the user is in the conversation
	check, err := rt.db.IsUserInConv(UserId, convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	if !check {
		http.Error(w, "user ins't in this conversation", http.StatusBadRequest)
		return
	}

	// Get messages
	messages, err := rt.db.GetMessagesByConvId(convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
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
				InternalServerError(w, err, "Internal server error: ")
				return
			}
			if check {
				continue
			}
			// Add the user to the list of users that have read the message
			err = rt.db.AddUserToListOfAlreadyReadersOfMess(messages[i].MessageId, UserId, convId)
			if err != nil {
				InternalServerError(w, err, "Internal server error: ")
				return
			}
			// Check if all the users have read the message
			check, err = rt.db.CheckAllUsersHaveReadMess(messages[i].MessageId)
			if err != nil {
				InternalServerError(w, err, "Internal server error: ")
				return
			}
			if check {
				// Update the message status
				err = rt.db.UpdateMessageStatus(messages[i].MessageId)
				if err != nil {
					InternalServerError(w, err, "Internal server error: ")
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

	type CommentData struct {
		CommentId int             `json:"commentId"`
		MessageId int             `json:"messageId"`
		Content   string          `json:"content"`
		Sender    structions.User `json:"sender"`
		ConvId    int             `json:"convId"`
	}

	type MessageData struct {
		Message  structions.Message `json:"message"`
		Sender   structions.User    `json:"sender"`
		DateTime string             `json:"dateTime"`
		Comments []CommentData      `json:"comments"`
	}

	var response []MessageData
	for i := 0; i < len(messages); i++ {
		// Get the sender
		sender, err := rt.db.GetUserById(messages[i].SenderId)
		if err != nil {
			InternalServerError(w, err, "Internal server error: ")
			return
		}

		var commentsData []CommentData

		// Get the comments
		comments, err := rt.db.GetCommentsByMessId(messages[i].MessageId, convId)
		if err != nil {
			InternalServerError(w, err, "Internal server error: ")
			return
		}

		// Save the commments in commentsData
		for j := 0; j < len(comments); j++ {
			// Get the sender
			senderComment, err := rt.db.GetUserById(comments[j].SenderId)
			if err != nil {
				InternalServerError(w, err, "Internal server error: ")
				return
			}
			commentsData = append(commentsData, CommentData{comments[j].CommentId, comments[j].MessageId, comments[j].Content, senderComment, comments[j].ConvId})
		}

		// Get the date and time
		dateTime := messages[i].DateTime.Format("15:04 - 02/01/2006")

		response = append(response, MessageData{messages[i], sender, dateTime, commentsData})
	}

	// Sort the messages by message dateTime
	sort.Slice(response, func(i, j int) bool {
		return response[i].Message.DateTime.After(response[j].Message.DateTime)
	})

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerError(w, err, "Error encoding response: ")
		return
	}
}
