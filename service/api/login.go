package api

import (
	"encoding/json"
	"net/http"

	"github.com/antoniomonaco/WasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// doLoginHandler gestisce la richiesta di login.
// Se l'utente non esiste, viene creato.
// Restituisce un identificativo univoco per l'utente.
func (rt *_router) doLoginHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Decodifico la richiesta e leggo il nome utente.
	var request struct {
		Name string `json:"name"`
	}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica della richiesta", http.StatusBadRequest)
		return
	}

	// 2. Verifico se l'utente esiste già
	userID, err := rt.db.RetrieveUser(request.Name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// 4. Creo la risposta con l'userID come identificativo univoco
	response := struct {
		Identifier int `json:"identifier"`
	}{
		Identifier: userID,
	}

	// 5. Imposto l'header Content-Type e codifico la risposta in JSON.
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}
