package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.POST("/session", rt.wrap(rt.doLogin, false))
	rt.router.GET("/users", rt.wrap(rt.SearchUsers, true))
	rt.router.PUT("/users/:user/photo", rt.wrap(rt.SetMyPhoto, true))
	rt.router.PUT("/users/:user/username", rt.wrap(rt.UsernameModify, true))
	rt.router.POST("/users/:user/groups", rt.wrap(rt.CreateGroup, true))
	rt.router.DELETE("/users/:user/groups/:group_id", rt.wrap(rt.LeaveGroup, true))
	rt.router.PUT("/users/:user/groups/:group_id", rt.wrap(rt.AddToGroup, true))
	rt.router.PUT("/users/:user/groups/:group_id/name", rt.wrap(rt.GroupNameModify, true))
	rt.router.PUT("/users/:user/groups/:group_id/photo", rt.wrap(rt.SetGroupPhoto, true))
	rt.router.GET("/users/:user/conversations", rt.wrap(rt.GetConversations, true))
	rt.router.PUT("/users/:user/conversations", rt.wrap(rt.CreateConv, true))
	rt.router.GET("/users/:user/conversations/:conv_id", rt.wrap(rt.GetConversation, true))
	rt.router.POST("/users/:user/conversations/:conv_id/messages", rt.wrap(rt.SendMessage, true))
	rt.router.DELETE("/users/:user/conversations/:conv_id/messages/:mess_id", rt.wrap(rt.DeleteMessage, true))
	rt.router.POST("/users/:user/conversations/:conv_id/messages/:mess_id", rt.wrap(rt.ForwardMessage, true))
	rt.router.PUT("/users/:user/conversations/:conv_id/messages/:mess_id/comments", rt.wrap(rt.CommentMessage, true))
	rt.router.DELETE("/users/:user/conversations/:conv_id/messages/:mess_id/comments/:comm_id", rt.wrap(rt.UncommentMessage, true))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
