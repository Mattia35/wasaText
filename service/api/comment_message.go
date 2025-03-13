package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"progetto.wasa/service/api/structions"
	"unicode"
)

func (rt *_router) CommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	if _, err := rt.db.GetUserByConv(convId, UserId); err != nil {
		http.Error(w, "User isn't in the conversation"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the message id
	messId, err := strconv.Atoi(ps.ByName("mess_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the message is in the conversation
	mes, err := rt.db.GetMessageById(messId, convId)
	if err != nil {
		// stampa convId e messId
		http.Error(w, "Message isn't in the conversation"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is the sender of the message
	if mes.SenderId == UserId {
		http.Error(w, "User is the sender of the message, so he can't comment it", http.StatusBadRequest)
		return
	}
	type RequestBody struct {
		Emoji string `json:"emoji"`
	}
	var request RequestBody

	// Get the text of the comment
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	
	// Check if the string is an emoji
	runes := []rune(request.Emoji) 
	if len(runes) != 1 {
		http.Error(w, "The string has more than 1 character", http.StatusBadRequest)
		return
	}
	runee := runes[0]
	check := unicode.Is(unicode.So, runee)
	if !check {
		http.Error(w, "The string isn't a emoji", http.StatusBadRequest)
		return
	}

	var comment structions.Comment
	comment.SenderId = UserId
	comment.MessageId = messId
	comment.ConvId = convId
	comment.Content = request.Emoji

	// Check if the user has already commented the message
	check, err = rt.db.CheckIfUserHasAlreadyCommented(messId, UserId, convId);
	if err != nil {
		http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
		return
	}

	if check {
		// Update the comment in the db
		err = rt.db.UpdateComment(UserId, messId, convId, request.Emoji)
		if err != nil {
			http.Error(w, "Error update the comment in the database"+err.Error(), http.StatusBadRequest)
			return
		}
	} else{
		// Insert the comment in the db
		comment, err = rt.db.CreateComment(comment)
		if err != nil {
			http.Error(w, "Error insert the comment in the database"+err.Error(), http.StatusBadRequest)
			return
		}
	}

	type Response struct {
		NewComment  structions.Comment `json:"newComment"`
	}
	var response Response
	response.NewComment = comment

	// Response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Error encoding response"+err.Error(), http.StatusInternalServerError)
		return
	}

}