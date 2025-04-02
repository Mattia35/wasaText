package api

import (
	"net/http"
)

func InternalServerError(w http.ResponseWriter, err error, message string) {
	http.Error(w, message+err.Error(), http.StatusInternalServerError)
}

func BadRequest(w http.ResponseWriter, err error, message string) {
	http.Error(w, message+err.Error(), http.StatusBadRequest)
}

func Forbidden(w http.ResponseWriter, err error, message string) {
	http.Error(w, message+err.Error(), http.StatusForbidden)
}
