package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/antoniomonaco/WasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

func (rt *_router) commentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {

	conversationID, _ := strconv.Atoi(ps.ByName("conversationID"))
	messageID, _ := strconv.Atoi(ps.ByName("messageID"))
	userID := reqcontext.UserIDFromContext(r.Context())

	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, userID)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden) // 403
		return
	}

	var request struct {
		Content string `json:"content"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica della richiesta", http.StatusBadRequest) // 400
		return
	}

	timestamp := time.Now()

	commentID, err := rt.db.AddComment(messageID, userID, request.Content, timestamp)
	if err != nil {
		http.Error(w, "Errore durante l'aggiunta del commento", http.StatusInternalServerError)
		return
	}

	user, err := composeUserFromID(userID, rt.db)
	if err != nil {
		http.Error(w, "Errore durante la composizione dell'utente", http.StatusInternalServerError)
		return
	}
	comment := Comment{
		ID:        commentID,
		MessageID: messageID,
		Content:   request.Content,
		Sender:    user,
		Timestamp: timestamp,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(comment)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getComments(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID, _ := strconv.Atoi(ps.ByName("conversationID"))
	messageID, _ := strconv.Atoi(ps.ByName("messageID"))
	userID := reqcontext.UserIDFromContext(r.Context())

	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, userID)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden) // 403
		return
	}

	rows, err := rt.db.GetComments(messageID)
	if err != nil {
		http.Error(w, "Errore durante il recupero dei commenti", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var comments []Comment
	for rows.Next() {
		var comment Comment
		var senderID int
		var timestamp string

		err := rows.Scan(&comment.ID, &senderID, &comment.MessageID, &comment.Content, &timestamp)
		if err != nil {
			http.Error(w, "Errore durante la lettura dei commenti", http.StatusInternalServerError)
			return
		}
		sender, err := composeUserFromID(senderID, rt.db)
		if err != nil {
			http.Error(w, "Errore durante la composizione dell'utente", http.StatusInternalServerError)
			return
		}
		t, _ := time.Parse(time.RFC3339, timestamp)
		comment.Sender = sender
		comment.Timestamp = t
		comments = append(comments, comment)
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(comments)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) uncommentMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID, _ := strconv.Atoi(ps.ByName("conversationID"))
	messageID, _ := strconv.Atoi(ps.ByName("messageID"))
	commentID, _ := strconv.Atoi(ps.ByName("commentID"))
	userID := reqcontext.UserIDFromContext(r.Context())

	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, userID)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden) // 403
		return
	}

	err = rt.db.DeleteComment(messageID, commentID, userID)
	if err != nil {
		http.Error(w, "Commento non torvato", http.StatusNotFound) // 404
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
