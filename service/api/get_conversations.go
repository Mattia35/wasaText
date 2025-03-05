package api
import (
	"encoding/json"
	"net/http"
	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"strconv"
	"progetto.wasa/service/api/structions"

)

func (rt *_router) GetConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Get the conversations of the user
	conversations, err := rt.db.GetConversationsByUserId(UserId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}





	// Struct used for the response
	type ConvObject struct {
		Conversation structions.Conversation 		`json:"conversation"`
		User         structions.User                `json:"user"`
		Group        structions.Group               `json:"group"`
		GroupUsers   []structions.User              `json:"groupUsers"`
		Message      structions.Message      		`json:"message"`
		SenderUser   structions.User                `json:"senderUser"`
	}

	// Response
	response := make([]ConvObject, len(conversations))

	// Fornire anche informazioni riguardo gruppi o utenti
	for idx, conv := range conversations {
		if conv.GroupId == 0 {
			// Get the user from the conversation
			userID, err := rt.db.GetUserByConv(conv.ConvId, UserId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}
			user, err := rt.db.GetUserById(userID.UserId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}
			var user User
			err = user.ConvertUserFromDB(userDB)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			// Get last message
			message, err := rt.db.GetMessageById(conv.LastMessageId, conv.ConversationId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			SenderUserDB, err := rt.db.GetUserById(message.SenderUserId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			var senderUser structions.User
			err = senderUser.ConvertUserFromDB(SenderUserDB)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			response[idx] = ConvObject{
				Conversation: conv,
				User:         user,
				Message:      message,
				SenderUser:   senderUser,
			}
		} else {
			// Get the group from the conversation
			groupDB, err := rt.db.GetGroupById(conv.GroupId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}
			var group structions.Group
			err = group.ConvertGroupFromDB(groupDB)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			// Get last message
			message, err := rt.db.GetMessageById(conv.LastMessageId, conv.ConversationId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			SenderUserDB, err := rt.db.GetUserById(message.SenderUserId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			var senderUser structions.User
			err = senderUser.ConvertUserFromDB(SenderUserDB)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			users, err := rt.db.GetMembers(conv.GroupId)
			if err != nil {
				http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
			}

			var groupUsers []structions.User
			for i := 0; i < len(users); i++ {
				userDB, err := rt.db.GetUserById(users[i].UserId)
				if err != nil {
					http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
				}
				var user structions.User
				err = user.ConvertUserFromDB(userDB)
				if err != nil {
					http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
				return
				}
				groupUsers = append(groupUsers, user)
			}

			response[idx] = ConvObject{
				Conversation: conv,
				Group:        group,
				GroupUsers:   groupUsers,
				Message:      message,
				SenderUser:   senderUser,
			}
		}
	}

	// Write the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(response); err != nil {
		ctx.Logger.WithError(err).Error("Error encoding response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}
	
}