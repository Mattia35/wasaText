package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))
	rt.router.PUT("/users/:user/username", rt.wrap(rt.UsernameModify, true))
	rt.router.PUT("/users/:user/groups/:group_id/name", rt.wrap(rt.GroupNameModify, true))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
