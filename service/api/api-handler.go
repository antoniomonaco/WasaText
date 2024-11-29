package api

import (
	"net/http"
)

// Handler returns an instance of httprouter.Router that handle APIs registered here
func (rt *_router) Handler() http.Handler {
	// Register routes
	rt.router.GET("/", rt.getHelloWorld)
	rt.router.GET("/context", rt.wrap(rt.getContextReply))
	rt.router.POST("/session", rt.wrap(rt.doLoginHandler))
	rt.router.GET("/users/", rt.authMiddleware(rt.wrap(rt.getUSersHandler)))
	rt.router.PUT("/users/me/name", rt.authMiddleware(rt.wrap(rt.setMyUserName)))
	rt.router.GET("/conversations/", rt.authMiddleware(rt.wrap(rt.getMyConversationsHandler)))
	rt.router.POST("/conversations/", rt.authMiddleware(rt.wrap(rt.createConversationHandler)))
	rt.router.GET("/conversations/:conversationID", rt.authMiddleware(rt.wrap(rt.getConversationHandler)))

	// Special routes
	rt.router.GET("/liveness", rt.liveness)

	return rt.router
}
