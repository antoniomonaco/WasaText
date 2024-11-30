package api

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"

	"github.com/antoniomonaco/WasaText/service/api/reqcontext"
	"github.com/antoniomonaco/WasaText/service/database"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) getUSersHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	IDFromContext := reqcontext.UserIDFromContext(r.Context()) // L'utente non può cercare se stesso, quindi passo l'id alla funzione nel db

	username := r.URL.Query().Get("name") // Se ho passato un nome lo imposta, altrimenti ritorna una stringa vuota
	rows, err := rt.db.RetrieveUsers(username, IDFromContext)
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
		http.Error(w, "Nessun utente trovato", http.StatusNotFound) // 404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(users)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError) // 500
		return
	}
}

func (rt *_router) setMyUserName(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())

	var request struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica della richiesta", http.StatusBadRequest)
		return
	}

	pattern := "^[a-zA-Z0-9]{3,16}$"
	re := regexp.MustCompile(pattern) //regex per il nome utente

	// Controllo che il nome passato rispetta il pattern
	if !re.MatchString(request.Name) {
		http.Error(w, "Nome utente non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.UpdateUserName(request.Name, IDFromContext)
	if err != nil {
		if err.Error() == fmt.Sprintf("il nome utente '%s' è già in uso", request.Name) {
			http.Error(w, "Nome utente già in uso", http.StatusConflict) // 409 Conflict
		} else {
			http.Error(w, "Errore nella modifica del nome utente", http.StatusInternalServerError)
		}
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204

}

func (rt *_router) setMyPhoto(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())

	var request struct {
		PhotoUrl string `json:"photoUrl"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica della richiesta", http.StatusBadRequest)
		return
	}

	err = rt.db.SetUserPhoto(request.PhotoUrl, IDFromContext)

	if err != nil {
		http.Error(w, "Impossibile impostare la foto profilo", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204

}

func composeUser(userID int, username string, photoUrl string) User {
	user := User{
		ID:       userID,
		Username: username,
		PhotoUrl: photoUrl,
	}

	return user
}

func composeUserFromID(userID int, db database.AppDatabase) (User, error) {

	var username string
	var photoUrl string
	var user User

	rows, err := db.RetrieveUserFromID(userID)
	if err != nil {

		return user, fmt.Errorf("impossibile trovare l'utente: %w", err)
	}

	rows.Scan(&username, &photoUrl)
	user = User{
		ID:       userID,
		Username: username,
		PhotoUrl: photoUrl,
	}

	return user, nil
}
