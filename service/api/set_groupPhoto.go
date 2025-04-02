package api

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"image/jpeg"
	"io"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"github.com/nfnt/resize"
	"progetto.wasa/service/api/reqcontext"
)

func (rt *_router) SetGroupPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
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

	// Get the groupId
	groupId, err := strconv.Atoi(ps.ByName("group_id"))
	if err != nil {
		BadRequest(w, err, "Bad Request: ")
		return
	}

	// Check if the user is in the group
	check, err := rt.db.IsUserInGroup(UserId, groupId)

	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
		return
	}
	if !check {
		http.Error(w, "User isn't in the group", http.StatusBadRequest)
		return
	}

	// Check the weight of the photo
	err = r.ParseMultipartForm(5 << 20)
	if err != nil {
		BadRequest(w, err, "The image is too big: ")
		return
	}

	// Get the file
	file, _, err := r.FormFile("image")
	if err != nil {
		BadRequest(w, err, "Error getting the image file: ")
		return
	}
	defer file.Close()

	// Check if there is a photo in the request
	if file == nil {
		BadRequest(w, err, "The photo isn't in the request: ")
		return
	}

	// Encode the file in base64
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

	// Decode the file in image
	img, err := jpeg.Decode(bytes.NewReader(data))
	if err != nil {
		InternalServerError(w, err, "Error decoding the file: ")
		return
	}

	// Resize the image
	newImg := resize.Resize(250, 250, img, resize.Lanczos3)

	// Encode the resized image to a buffer
	var buf bytes.Buffer
	err = jpeg.Encode(&buf, newImg, nil)
	if err != nil {
		InternalServerError(w, err, "Error encoding resized image: ")
		return
	}

	// Encode the resized image to Base64
	response := base64.StdEncoding.EncodeToString(buf.Bytes())

	// Set the photo in the record
	err = rt.db.SetGroupPhoto(groupId, response)
	if err != nil {
		InternalServerError(w, err, "Internal server error: ")
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
		InternalServerError(w, err, "Error encoding the response: ")
		return
	}
}
