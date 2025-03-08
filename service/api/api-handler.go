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
	rt.router.POST("/users/:user/groups", rt.wrap(rt.CreateGroup, true))
	rt.router.PUT("/users/:user/conversations/:receiver_id", rt.wrap(rt.CreateConv, true))
	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
