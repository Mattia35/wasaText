package api

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
)

func (rt *_router) UncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	if _, err := rt.db.GetUserByConv(convId, UserId); err != nil {
		BadRequest(w, err, "User isn't in the conversation: ")
		return
	}

	// Get the message id
	messId, err := strconv.Atoi(ps.ByName("mess_id"))
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the message is in the conversation
	_, err = rt.db.GetMessageById(messId, convId)
	if err != nil {
		BadRequest(w, err, "Message isn't in the conversation: ")
		return
	}
	// Get the comment id
	commId, err := strconv.Atoi(ps.ByName("comm_id"))
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the comment is in the conversation
	check, err := rt.db.IsCommentInConv(commId, messId, convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	if !check {
		http.Error(w, "Comment isn't in the conversation", http.StatusBadRequest)
		return
	}

	// Get the comment
	comm, err := rt.db.GetCommentById(commId, messId, convId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		BadRequest(w, err, "Comment isn't in the conversation: ")
		return
	}

	// Check if the user is the sender of the comment
	if comm.SenderId != UserId {
		http.Error(w, "User isn't the sender of the comment, so he can't delete it", http.StatusBadRequest)
		return
	}

	// Remove the message
	err = rt.db.RemoveComment(commId, messId, convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}

	w.WriteHeader(http.StatusOK)

	// Send the response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode("comment has been successfully deleted from list comments!"); err != nil {
		InternalServerError(w, err, "Can't encode the response: ")
		return
	}
}
