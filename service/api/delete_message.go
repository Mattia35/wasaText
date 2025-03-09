package api

import (
	"net/http"
	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"strconv"
	"encoding/json"
)

func (rt *_router) DeleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Check if the user is the sender of the message
	check, err := rt.db.CheckMessageSender(messId, UserId); 
	if check==false {
		http.Error(w, "User isn't the sender of the message", http.StatusBadRequest)
		return
	}
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}
	// Get the max message id
	maxMessId, err := rt.db.GetMaxMessageId(convId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Remove the message
	err = rt.db.RemoveMessage(messId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Get the new max message id
	newMaxMessId, err := rt.db.GetMaxMessageId(convId)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}

	// If the message is the last of the conversation, update the new last message
	if maxMessId == messId {
		// Update the last message of the conversation
		err = rt.db.AddMessageToConv(newMaxMessId, convId)
		if err != nil {
			http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
	
	// Send the response
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode("message has been successfully deleted from conversation!"); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}