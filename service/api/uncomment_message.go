package api

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
)

func (rt *_router) UncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	_, err = rt.db.GetMessageById(messId, convId)
	if err != nil {
		http.Error(w, "Message isn't in the conversation"+err.Error(), http.StatusBadRequest)
		return
	}
	// Get the comment id
	commId, err := strconv.Atoi(ps.ByName("comm_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the comment is in the conversation
	check, err := rt.db.IsCommentInConv(commId, messId, convId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	if !check {
		http.Error(w, "Comment isn't in the conversation", http.StatusBadRequest)
		return
	}

	// Get the comment
	comm, err := rt.db.GetCommentById(commId, messId, convId)
	if err != nil {
		if err == sql.ErrNoRows {
			http.Error(w, "Comment not found", http.StatusNotFound)
			return
		}
		http.Error(w, "Comment isn't in the conversation"+err.Error(), http.StatusBadRequest)
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
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)

	// Send the response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode("comment has been successfully deleted from list comments!"); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
