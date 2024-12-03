package api

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"

	"github.com/antoniomonaco/WasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) addToGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())
	conversationID, error := strconv.Atoi(ps.ByName("conversationID"))
	if error != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}

	var request struct {
		UserID int `json:"id"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica del messaggio", http.StatusBadRequest)
		return
	}

	isParticipant, err = rt.db.IsUserParticipantOfConversation(conversationID, request.UserID)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if isParticipant {
		http.Error(w, "L'utente fa già parte del gruppo", http.StatusForbidden)
		return
	}

	isGroup, err := rt.db.IsGroup(conversationID)
	if err != nil {
		http.Error(w, "Errore durante la verifica della conversazione", http.StatusInternalServerError)
		return
	}
	if !isGroup {
		http.Error(w, "non puoi aggiungere un utente ad una chat privata", http.StatusForbidden)
		return
	}

	err = rt.db.AddParticipant(conversationID, request.UserID)
	if err != nil {
		http.Error(w, "Errore durante l'aggiunta dell'utente al gruppo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) //204

}

func (rt *_router) leaveGroupHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())
	conversationID, error := strconv.Atoi(ps.ByName("conversationID"))
	if error != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}

	var request struct {
		UserID int `json:"id"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica del messaggio", http.StatusBadRequest)
		return
	}

	isParticipant, err = rt.db.IsUserParticipantOfConversation(conversationID, request.UserID)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "L'utente non fa parte del gruppo", http.StatusForbidden)
		return
	}

	isGroup, err := rt.db.IsGroup(conversationID)
	if err != nil {
		http.Error(w, "Errore durante la verifica della conversazione", http.StatusInternalServerError)
		return
	}
	if !isGroup {
		http.Error(w, "non puoi rimuovere un utente da una chat privata", http.StatusForbidden)
		return
	}

	err = rt.db.RemoveParticipant(conversationID, request.UserID)
	if err != nil {
		http.Error(w, "Errore durante la rimozione dell'utente dal gruppo", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) //204

}

func (rt *_router) setGroupNameHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())
	conversationID, error := strconv.Atoi(ps.ByName("conversationID"))
	if error != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}

	var request struct {
		GroupName string `json:"name"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica del messaggio", http.StatusBadRequest)
		return
	}

	isGroup, err := rt.db.IsGroup(conversationID)
	if err != nil {
		http.Error(w, "Errore durante la verifica della conversazione", http.StatusInternalServerError)
		return
	}
	if !isGroup {
		http.Error(w, "non puoi cambiare il nome ad una chat privata", http.StatusForbidden)
		return
	}
	pattern := "^.*$"
	re := regexp.MustCompile(pattern)

	if !re.MatchString(request.GroupName) {
		http.Error(w, "Nome utente non valido", http.StatusBadRequest)
		return
	}

	err = rt.db.UpdateGroupName(request.GroupName, conversationID)
	if err != nil {
		http.Error(w, "Errore nella modifica del nome utente", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}

func (rt *_router) setGroupPhotoHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())
	conversationID, error := strconv.Atoi(ps.ByName("conversationID"))
	if error != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}

	var request struct {
		PhotoUrl string `json:"photoUrl"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica del messaggio", http.StatusBadRequest)
		return
	}

	isGroup, err := rt.db.IsGroup(conversationID)
	if err != nil {
		http.Error(w, "Errore durante la verifica della conversazione", http.StatusInternalServerError)
		return
	}
	if !isGroup {
		http.Error(w, "non puoi cambiare la foto di una chat privata", http.StatusForbidden)
		return
	}

	err = rt.db.UpdateGroupPhoto(request.PhotoUrl, conversationID)
	if err != nil {
		http.Error(w, "Errore nella modifica del nla foto profilo", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204
}
