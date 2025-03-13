package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"progetto.wasa/service/api/structions"
)

func (rt *_router) ForwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext){
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

	// Get the conversation id
	convId, err := strconv.Atoi(ps.ByName("conv_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the message id
	messId, err := strconv.Atoi(ps.ByName("mess_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is in the conversation
	if _, err := rt.db.GetUserByConv(convId, UserId); err != nil {
		http.Error(w, "User isn't in the conversation"+err.Error(), http.StatusBadRequest)
		return
	}

	type Destination struct {
		DestGroup int `json:"group"`
		DestUser  int `json:"user"`
	}
	type RequestBody struct {
		DestConvId []Destination `json:"destination"`
	}
	var request RequestBody

	// Get request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	var requestToDelete RequestBody
	// Check if destinations are all different
	for i := 0; i < len(request.DestConvId); i++ {
		check := false
		for j := i + 1; j < len(request.DestConvId); j++ {
			if request.DestConvId[i].DestGroup != 0 {
				if request.DestConvId[i].DestGroup == request.DestConvId[j].DestGroup{
					check = true
					break
				}
			} else if request.DestConvId[i].DestUser != 0 {
				if request.DestConvId[i].DestUser == request.DestConvId[j].DestUser{
					check = true
					break
				}
			}
		}
		if check == true {
			requestToDelete.DestConvId = append(requestToDelete.DestConvId, request.DestConvId[i])
		}
	}
	var newRequest RequestBody

	// Delete all duplicated request
	for i := 0; i < len(request.DestConvId); i++ {
		check := false
		for j := i + 1; j < len(requestToDelete.DestConvId); j++ {
			if request.DestConvId[i].DestGroup != 0{
				if request.DestConvId[i].DestGroup == requestToDelete.DestConvId[j].DestGroup {
					check = true
					break
				}
			} else if request.DestConvId[i].DestUser != 0{
				if request.DestConvId[i].DestUser != requestToDelete.DestConvId[j].DestUser {
					check = true
					break
				}
			}
		}
		if check == false {
			newRequest.DestConvId = append(newRequest.DestConvId, request.DestConvId[i])
		}
	}
	request = newRequest
	Conver:=0;
	var messages []structions.Message

	// For each destination, check if it is a group or a user
	for i := 0; i < len(request.DestConvId); i++ {
		if request.DestConvId[i].DestGroup != 0 {
			// Get conv by group id
			AllConv, err := rt.db.GetConvByGroupId(request.DestConvId[i].DestGroup )
			if err != nil {
				http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			// Check if the conversation exists
			Conver = AllConv.ConvId
			check, err := rt.db.DoConversationExist(Conver)
			if err != nil {
				http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			// If the conversation doesn't exist, return an error	
			if check == false {
				http.Error(w, "The group doesn't exist", http.StatusBadRequest)
				return
			}
		} else {
			// Get destination by username
			Dest, err := rt.db.GetUserById(request.DestConvId[i].DestUser)
			if err != nil {
				http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			// Get the conversation by destination id
			Conver, err = rt.db.GetConversationByUsers(UserId, Dest.UserId)
			if err != nil && err.Error() != "sql: no rows in result set" {
				http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
				return
			}
			if err.Error() == "sql: no rows in result set" {
				var conversation structions.Conversation
				conversation.GroupId = 0
				conversation, err = rt.db.CreateConversation(conversation)
				if err != nil {
					ctx.Logger.WithError(err).Error("can't create the conversation of the private chat")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// Add the user to the conversation
				err = rt.db.AddUserToConv(UserId, conversation.ConvId)
				if err != nil {
					ctx.Logger.WithError(err).Error("you can't add a user to the conversation")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				// Add the destination user to the conversation
				err = rt.db.AddUserToConv(Dest.UserId, conversation.ConvId)
				if err != nil {
					ctx.Logger.WithError(err).Error("you can't add a user to the conversation")
					w.WriteHeader(http.StatusInternalServerError)
					return
				}
				Conver = conversation.ConvId
			}
		}
		// Check if the user is in the conversation/destionation
		if _, err := rt.db.GetUserByConv(Conver, UserId); err != nil {
			http.Error(w, "User isn't in the conversation"+err.Error(), http.StatusBadRequest)
			return
		}

		// Get the message
		message, err := rt.db.GetMessageById(messId, convId)
		if err != nil {
			http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
			return
		}

		var mess structions.Message
		mess.Text = message.Text
		mess.Status = false
		mess.SenderId = UserId
		mess.ConvId = Conver
		mess.Photo = message.Photo

		// Send the message
		mess, err = rt.db.CreateMessage(mess)
		if err != nil {
			http.Error(w, "Error insert the message in the database"+err.Error(), http.StatusBadRequest)
			return
		}

		// Update the last message
		err = rt.db.AddMessageToConv(mess.MessageId, Conver)
		if err != nil {
			http.Error(w, "Error updating last message id"+err.Error(), http.StatusBadRequest)
			return
		}

		// get users of the conversation
		users, err := rt.db.GetUsersByConvId(Conver)
		if err != nil {
			http.Error(w, "Error taking the users of the conversation"+err.Error(), http.StatusBadRequest)
			return
		}
		
		// Set the users that have read the message: all the users of the group, unless the sender
		newUsers := make([]structions.User, 0)
		for i := 0; i < len(users); i++ {
			if users[i].UserId != UserId {
				newUsers = append(newUsers, users[i])
			}
		}
		users = newUsers

		// Set the users that have read the message
		for i := 0; i < len(users); i++ {
			err = rt.db.AddUserToListOfReadersOfMess(mess.MessageId, users[i].UserId, mess.ConvId)
			if err != nil {
				http.Error(w, "Error adding the user to the list of readers of the message"+err.Error(), http.StatusBadRequest)
				return
			}
		}
		messages = append(messages, mess)
	}

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(messages); err != nil {
		http.Error(w, "Error encoding response"+err.Error(), http.StatusInternalServerError)
		return
	}
	
}