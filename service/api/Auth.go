package api
import (
	"net/http"
	"strconv"
)
func isAuthorized(header http.Header) int{
	authToken, err := strconv.Atoi(header.Get("Authorization"))
	if err != nil {return 0}
	return authToken
}