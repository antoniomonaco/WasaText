package api

import (
	"net/http"
	"strconv"

	"github.com/antoniomonaco/WasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// Middleware per gestire il Bearer token
func (rt *_router) authMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Estraggo il Bearer token dall'header Authorization
		token := r.Header.Get("Authorization")
		if len(token) < 7 || token[:7] != "Bearer " {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		// Estraggo l'userID dal token
		userIDStr := token[7:] // Rimuovo la scritta "Bearer "
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Creo un nuovo contesto con l'userID
		ctx := reqcontext.WithUserID(r.Context(), userID)

		// Passo il nuovo contesto alla richiesta
		next(w, r.WithContext(ctx), ps)
	}
}
