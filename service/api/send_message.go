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
		BadRequest(w, err, "Bad Request: ")
		return
	}

	userID := ctx.UserID

	// Check if the user is authorized
	if UserId != userID {
		Forbidden(w, err, "Forbidden: ")
		return
	}

	// Get the conversation id
	convId, err := strconv.Atoi(ps.ByName("conv_id"))
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the user is in the conversation
	if _, err := rt.db.GetUserByConv(convId, userID); err != nil {
		BadRequest(w, err, "User isn't in the conversation: ")
		return
	}

	// Get the conversation from the id
	conv, err := rt.db.GetConvById(convId)
	if err != nil {
		http.Error(w, "Conversation not found: "+err.Error(), http.StatusNotFound)
		return
	}

	var mess structions.Message

	// Check the weight of the message
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		BadRequest(w, err, "The image is too big: ")
		return
	}
	// Get the text of the message
	mess.Text = r.FormValue("text")

	messIdToReplyTo := 0
	// Try to get the message to reply to. If it fails, it means that the message is not a reply, so continue
	if r.FormValue("messToReplyTo") != "" {
		messToReplyTo, err := strconv.Atoi(r.FormValue("messToReplyTo"))
		if err != nil {
			BadRequest(w, err, "Error taking the message id to reply to: ")
			return
		} else {
			messIdToReplyTo = messToReplyTo
		}
	}
	// Get the image
	file, _, err := r.FormFile("image")
	if err == nil && file != nil {
		buf := make([]byte, 1)
		n, _ := file.Read(buf)
		if n == 0 {
			file = nil
		}
	}
	// Get the gif
	fileGif, _, err1 := r.FormFile("gif")
	if err1 == nil && fileGif != nil {
		buf := make([]byte, 1)
		n, _ := fileGif.Read(buf)
		if n == 0 {
			fileGif = nil
		}
	}

	// Check if the message is empty
	if mess.Text == "" && file == nil && fileGif == nil {
		BadRequest(w, err, "The message is empty: ")
		return
	}

	// Check if the message has a gif, a text, a text with an image or an image
	thereIsText := false
	thereIsGif := false
	thereIsImage := false
	if mess.Text != "" {
		thereIsText = true
	}
	if fileGif != nil {
		thereIsGif = true
	}
	if file != nil {
		thereIsImage = true
	}
	if thereIsGif && (thereIsText || thereIsImage) {
		http.Error(w, "You can't send a gif with a text or an image", http.StatusBadRequest)
		return
	}
	// Check if the request have a file, and if it has, encode it
	if thereIsImage {
		if _, err := file.Seek(0, io.SeekStart); err != nil {
			InternalServerError(w, err, "Error seeking the image file: ")
			return
		}

		// Read the file
		data, err := io.ReadAll(file) // In data we have the image file taked in the request
		if err != nil {
			InternalServerError(w, err, "Error reading the image file: ")
			return
		}

		// Check if the file is a jpeg
		fileType := http.DetectContentType(data)
		if fileType != "image/jpeg" {
			http.Error(w, "Bad Request, wrong file type", http.StatusBadRequest)
			return
		}
		defer func() { err = file.Close() }()

		// Encode the image in base64
		mess.Photo = base64.StdEncoding.EncodeToString(data)
	}

	// Check if the request have a gif, and if it has, encode it
	if thereIsGif {
		if _, err := fileGif.Seek(0, io.SeekStart); err != nil {
			InternalServerError(w, err, "Error seeking the gif file: ")
			return
		}
		// Read the file
		dataGif, err := io.ReadAll(fileGif) // In data we have the gif file taked in the request
		if err != nil {
			InternalServerError(w, err, "Error reading the gif file: ")
			return
		}

		// Check if the file is a gif
		fileType := http.DetectContentType(dataGif)
		if fileType != "image/gif" {
			http.Error(w, "Bad Request, wrong file type", http.StatusBadRequest)
			return
		}
		defer func() { err = fileGif.Close() }()

		// Encode the gif in base64
		mess.Gif = base64.StdEncoding.EncodeToString(dataGif)
	}

	// Set the id of the conversation
	mess.ConvId = conv.ConvId
	mess.SenderId = UserId
	mess.Status = false

	// query message
	type Response struct {
		MessToreplyTo structions.Message `json:"messToReplyTo"`
		MessSended    structions.Message `json:"messSended"`
	}
	var response Response

	if messIdToReplyTo != 0 {
		// Get the message by the id
		MessToreplyTo, err := rt.db.GetMessageById(messIdToReplyTo, conv.ConvId)
		if err != nil {
			BadRequest(w, err, "Error taking the message by the id: ")
			return
		}
		// Set the message query in the response
		response.MessToreplyTo = MessToreplyTo

		// Set the reply id in the message
		mess.ReplyId = messIdToReplyTo
	}

	// Insert the message in the db
	mess, err = rt.db.CreateMessage(mess)
	if err != nil {
		BadRequest(w, err, "Error insert the message in the database: ")
		return
	}

	// Update the last message
	err = rt.db.AddMessageToConv(mess.MessageId, conv.ConvId)
	if err != nil {
		BadRequest(w, err, "Error updating last message id: ")
		return
	}

	// get users of the conversation
	users, err := rt.db.GetUsersByConvId(conv.ConvId)
	if err != nil {
		BadRequest(w, err, "Error taking the users of the conversation: ")
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
			BadRequest(w, err, "Error adding the user to the list of readers of the message: ")
			return
		}
	}

	// Set the message sended in the response
	response.MessSended = mess

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		InternalServerError(w, err, "Error encoding response: ")
		return
	}

}
