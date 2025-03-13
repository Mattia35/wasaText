package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"progetto.wasa/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/structions"
)

func (rt *_router) GroupNameModify (w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Get the groupId
	GroupId, err := strconv.Atoi(ps.ByName("group_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the group exists
	if _, err := rt.db.GetGroupByGroupId(GroupId); err != nil {
		http.Error(w, "Group doesn't exist"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user could modify the groupname
	check, err := rt.db.UserControlByGroup(UserId, GroupId); 
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	
	if check == false {
		http.Error(w, "User can't modify the groupname, because isn't in the group", http.StatusBadRequest)
		return
	}
	type RequestBody struct {
		Groupname string `json:"groupname"`
	}
	var request RequestBody

	var group structions.Group
	// Check if the user makes a bad request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}
	group.Username = request.Groupname

	// Check if the group respects the regex, so there is a bad request
	if !group.IsValid() {
		http.Error(w, "Invalid groupname", http.StatusBadRequest)
		return
	}
	group.GroupId = GroupId
	// Try to modify the groupname. If it fails, it gives an error
	if err := rt.db.GroupnameModify(GroupId, group.Username); err!= nil {
		http.Error(w, "Groupname modify failed. Retry!", http.StatusBadRequest)
		return
	}

	// Groupname changed, response 200
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "application/json")
	if err := json.NewEncoder(w).Encode(group); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}