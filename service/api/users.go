package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	"github.com/antoniomonaco/WasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUSersHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	username := r.URL.Query().Get("name") // Se ho passato un nome lo imposta, altrimenti ritorna una stringa vuota
	rows, err := rt.db.RetrieveUsers(username)
	if err != nil {
		http.Error(w, "Errore durante il recupero degli utenti : 1", http.StatusInternalServerError)
		return
	}

	var users []User
	var userID int
	var photoUrl sql.NullString // La stringa potrebbe essere nulla quindi uso questo tipo di sql per non avere errori nello scan

	for rows.Next() {
		err := rows.Scan(&userID, &username, &photoUrl)
		if err != nil {
			http.Error(w, "Errore durante il recupero degli utenti : 2", http.StatusInternalServerError)
			return
		}
		users = append(users, composeUser(userID, username, photoUrl.String))

	}
	defer rows.Close()

	if len(users) == 0 {
		http.Error(w, "Nessun utente trovato", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func composeUser(userID int, username string, photoUrl string) User {
	user := User{
		ID:       userID,
		Username: username,
		PhotoUrl: photoUrl,
	}

	return user
}
