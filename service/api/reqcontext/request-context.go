/*
Package reqcontext contains the request context. Each request will have its own instance of RequestContext filled by the
middleware code in the api-context-wrapper.go (parent package).

Each value here should be assumed valid only per request only, with some exceptions like the logger.
*/
package reqcontext

import (
	"context"

	"github.com/gofrs/uuid"
	"github.com/sirupsen/logrus"
)

type KeyType string

const UserIDKey KeyType = "userID"

// RequestContext is the context of the request, for request-dependent parameters
type RequestContext struct {
	// ReqUUID is the request unique ID
	ReqUUID uuid.UUID

	// Logger is a custom field logger for the request
	Logger logrus.FieldLogger

	UserID int
}

// Serve ad aggiungere l'userID al contesto
func WithUserID(ctx context.Context, userID int) context.Context {
	return context.WithValue(ctx, UserIDKey, userID)
}

// Serve a recuperare l'userID dal contesto
func UserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserIDKey).(int)
	if !ok {
		return 0 // 0 se non l'ha trovato
	}
	return userID
}
