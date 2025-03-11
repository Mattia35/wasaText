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

	type RequestBody struct {
		DestConvId int `json:"destination"`
	}
	var request RequestBody

	// Get request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if destination exist
	check, err := rt.db.DoConversationExist(request.DestConvId)

	if err != nil {
		http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
		return
	}
	if check == false {
		http.Error(w, "The conversation doesn't exist", http.StatusBadRequest)
		return
	}

	// Check if the user is in the conversation/destionation
	if _, err := rt.db.GetUserByConv(request.DestConvId, UserId); err != nil {
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
	mess.ConvId = request.DestConvId
	mess.Photo = message.Photo

	// Send the message
	mess, err = rt.db.CreateMessage(mess)
	if err != nil {
		http.Error(w, "Error insert the message in the database"+err.Error(), http.StatusBadRequest)
		return
	}

	// Update the last message
	err = rt.db.AddMessageToConv(mess.MessageId, mess.ConvId)
	if err != nil {
		http.Error(w, "Error updating last message id"+err.Error(), http.StatusBadRequest)
		return
	}

	// get users of the conversation
	users, err := rt.db.GetUsersByConvId(mess.ConvId)
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

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(mess); err != nil {
		http.Error(w, "Error encoding response"+err.Error(), http.StatusInternalServerError)
		return
	}
	
}