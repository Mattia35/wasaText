package api


import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"encoding/base64"
	"github.com/julienschmidt/httprouter"
	"progetto.wasa/service/api/reqcontext"
	"image/jpeg"
	"github.com/nfnt/resize"
	"bytes"
)

func (rt *_router) SetGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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
	groupId, err := strconv.Atoi(ps.ByName("group_id"))
	if err != nil {
		http.Error(w, "Bad Request"+err.Error(), http.StatusBadRequest)
		return
	}

	// Check if the user is in the group
	check, err := rt.db.IsUserInGroup(UserId, groupId)

	if err != nil {
		http.Error(w, "Internal server error"+err.Error(), http.StatusInternalServerError)
		return
	}
	if !check {
		http.Error(w, "User isn't in the group", http.StatusBadRequest)
		return
	}

	// Check the weight of the photo
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		http.Error(w, "The image is too big"+err.Error(), http.StatusBadRequest)
		return
	}

	// Get the file
	file, _, err := r.FormFile("image")
	if err != nil {
		http.Error(w, "Error getting the image file"+err.Error(), http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Check if there is a photo in the request
	if file == nil {
		http.Error(w, "The photo isn't in the request!", http.StatusBadRequest)
		return
	}

	// Encode the file in base64
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

	// Decode the file in image
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		http.Error(w, "Error decoding the file"+err.Error(), http.StatusInternalServerError)
		return
	}

	// Resize the image
	newImg := resize.Resize(250, 250, img, resize.Lanczos3)

	// Encode the resized image to a buffer
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, newImg, nil)
	if err != nil {
		http.Error(w, "Error encoding resized image: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Encode the resized image to Base64
	response := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Set the photo in the record
	err = rt.db.SetGroupPhoto(groupId, response)
	if err != nil {
		http.Error(w, "Internal Server Error"+err.Error(), http.StatusInternalServerError)
		return
	}

	type Response struct {
		Photo string `json:"photo"`
	}

	var res Response
	// Set the response
	res.Photo = response

	// Send the response
	w.WriteHeader(http.StatusOK)
	w.Header().Set("content-type", "plain/text")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		http.Error(w, "Error encoding the response"+err.Error(), http.StatusInternalServerError)
		return
	}
}