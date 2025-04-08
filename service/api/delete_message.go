package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
)

func (rt *_router) DeleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Get the message id
	messId, err := strconv.Atoi(ps.ByName("mess_id"))
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the user is in the conversation
	if _, err := rt.db.GetUserByConv(convId, UserId); err != nil {
		BadRequest(w, err, "User isn't in the conversation: ")
		return
	}

	// Check if the user is the sender of the message
	check, err := rt.db.CheckMessageSender(messId, UserId, convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	if !check {
		BadRequest(w, err, "User isn't the sender of the message: ")
		return
	}
	// Get the previus max message id
	maxMessId, err := rt.db.GetMaxMessageId(convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}

	// Remove the message
	err = rt.db.RemoveMessage(messId, convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}

	// Get the new max message id
	newMaxMessId, err := rt.db.GetMaxMessageId(convId)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}

	// If the message is the last of the conversation, update the new last message
	if maxMessId == messId {
		// Update the last message of the conversation
		err = rt.db.AddMessageToConv(newMaxMessId, convId)
		if err != nil {
			InternalServerError(w, err, "Internal Server Error: ")
			return
		}
	}

	w.WriteHeader(http.StatusOK)

	// Send the response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode("message has been successfully deleted from conversation!"); err != nil {
		InternalServerError(w, err, "Error encoding response: ")
		return
	}
}
