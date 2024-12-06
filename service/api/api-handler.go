package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))

	// Login
	rt.router.POST("/session", rt.wrap(rt.doLoginHandler))

	// Users
	rt.router.GET("/users/", rt.authMiddleware(rt.wrap(rt.searchUsers)))
	rt.router.PUT("/users/me/name", rt.authMiddleware(rt.wrap(rt.setMyUserName)))
	rt.router.PUT("/users/me/photo", rt.authMiddleware(rt.wrap(rt.setMyPhoto)))

	// Conversations
	rt.router.GET("/conversations/", rt.authMiddleware(rt.wrap(rt.getMyConversations)))
	rt.router.POST("/conversations/", rt.authMiddleware(rt.wrap(rt.createConversation)))
	rt.router.GET("/conversations/:conversationID", rt.authMiddleware(rt.wrap(rt.getConversation)))
	rt.router.POST("/conversations/:conversationID", rt.authMiddleware(rt.wrap(rt.sendMessage)))
	rt.router.DELETE("/conversations/:conversationID/messages/:messageID", rt.authMiddleware(rt.wrap(rt.deleteMessage)))
	rt.router.POST("/conversations/:conversationID/messages/:messageID", rt.authMiddleware(rt.wrap(rt.forwardMessage)))

	// Group
	rt.router.POST("/conversations/:conversationID/participants", rt.authMiddleware(rt.wrap(rt.addToGroup)))
	rt.router.DELETE("/conversations/:conversationID/participants", rt.authMiddleware(rt.wrap(rt.leaveGroup)))
	rt.router.PUT("/conversations/:conversationID/name", rt.authMiddleware(rt.wrap(rt.setGroupName)))
	rt.router.PUT("/conversations/:conversationID/photo", rt.authMiddleware(rt.wrap(rt.setGroupPhoto)))

	// Comments
	rt.router.POST("/conversations/:conversationID/messages/:messageID/comments", rt.authMiddleware(rt.wrap(rt.commentMessage)))
	rt.router.GET("/conversations/:conversationID/messages/:messageID/comments", rt.authMiddleware(rt.wrap(rt.getComments)))
	rt.router.DELETE("/conversations/:conversationID/messages/:messageID/comments/:commentID", rt.authMiddleware(rt.wrap(rt.uncommentMessage)))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
