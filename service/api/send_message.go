package api

import (
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"progetto.wasa/service/api/structions"
)

func (rt *_router) SendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Check if the user is in the conversation
	if _, err := rt.db.GetUserByConv(convId, userID); err != nil {
		http.Error(w, "User isn't in the conversation"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the conversation from the id
	conv, err := rt.db.GetConvById(convId)
	if err != nil {
		http.Error(w, "Conversation not found"+err.Error(), http.StatusNotFound)
		return
	}
	
	var mess structions.Message

	// Check the weight of the message
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "The image is too big"+err.Error(), http.StatusBadRequest)
		return
	}
	// Get the text of the message
	mess.Text = r.FormValue("text")

	messIdToReplyTo := 0
	// Try to get the message to reply to. If it fails, it means that the message is not a reply, so continue
	if r.FormValue("messToReplyTo") != "" {
		messToReplyTo, err := strconv.Atoi(r.FormValue("messToReplyTo"))
		if err != nil {
			http.Error(w, "Error taking the message id to reply to"+err.Error(), http.StatusBadRequest)
			return
		} else {
			messIdToReplyTo = messToReplyTo
		}
	}
	

	// Get the file
	file, _, err := r.FormFile("image")

	// Check if the message is empty
	if mess.Text == "null" && file == nil {
		http.Error(w, "The message is empty!"+err.Error(), http.StatusBadRequest)
		return
	}


	// Check if the request have a file, and if it has, encode it 
	if err == nil {
		// Read the file
		data, err := io.ReadAll(file) // In data we have the image file taked in the request
		if err != nil {
			http.Error(w, "Error reading the image file"+err.Error(), http.StatusInternalServerError)
			return
		}

		// Check if the file is a jpeg
		fileType := http.DetectContentType(data)
		if fileType != "image/jpeg" {
			http.Error(w, "Bad Request, wrong file type", http.StatusBadRequest)
			return
		}
		defer func() { err = file.Close() }()

		mess.Photo = base64.StdEncoding.EncodeToString(data)
	}
	// Set the id of the conversation
	mess.ConvId = conv.ConvId
	mess.SenderId = UserId
	mess.Status = false

	// query message
	type Response struct {
		MessToreplyTo  structions.Message `json:"messToReplyTo"`
		MessSended structions.Message `json:"messSended"`
	}
	var response Response
	
	if messIdToReplyTo != 0 {
		// Get the message by the id
		MessToreplyTo, err := rt.db.GetMessageById(messIdToReplyTo, conv.ConvId)
		if err != nil {
			http.Error(w, "Error taking the message by the id"+err.Error(), http.StatusBadRequest)
			return
		}
		// Set the message query in the response
		response.MessToreplyTo = MessToreplyTo
	}

	// Insert the message in the db
	mess, err = rt.db.CreateMessage(mess)
	if err != nil {
		http.Error(w, "Error insert the message in the database"+err.Error(), http.StatusBadRequest)
		return
	}

	// Update the last message
	err = rt.db.AddMessageToConv(mess.MessageId, conv.ConvId)
	if err != nil {
		http.Error(w, "Error updating last message id"+err.Error(), http.StatusBadRequest)
		return
	}

	// get users of the conversation
	users, err := rt.db.GetUsersByConvId(conv.ConvId)
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
		err = rt.db.AddUserToListOfReadersOfMess(mess.MessageId, users[i].UserId, conv.ConvId)
		if err != nil {
			http.Error(w, "Error adding the user to the list of readers of the message"+err.Error(), http.StatusBadRequest)
			return
		}
	}

	// Set the message sended in the response
	response.MessSended = mess

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response"+err.Error(), http.StatusInternalServerError)
		return
	}

}