package api

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (rt *_router) authMiddleware(next httprouter.Handle) httprouter.Handle {
	return func(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
		// Recupero il token dall'header Authorization
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" || len(authHeader) < 7 || authHeader[:7] != "Bearer " {
			http.Error(w, "Token mancante o non valido", http.StatusUnauthorized)
			return
		}

		tokenString := authHeader[7:] // Rimuove "Bearer "

		// Valida il token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method")
			}
			return secretKey, nil
		})

		if err != nil || !token.Valid {
			http.Error(w, "Token non valido", http.StatusUnauthorized)
			return
		}

		// Recupera l'userID dai claims del token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			http.Error(w, "Token non valido", http.StatusUnauthorized)
			return
		}

		userID, ok := claims["userID"].(float64)
		if !ok {
			http.Error(w, "Token non valido", http.StatusUnauthorized)
			return
		}

		// Passa l'userID nel contesto (o in un altro modo)
		ctx := reqcontext.WithUserID(r.Context(), int(userID))
		next(w, r.WithContext(ctx), ps)
	}
}
