package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/antoniomonaco/WasaText/service/api/reqcontext"
	"github.com/julienschmidt/httprouter"
)

// getMyConversationsHandler gestisce la richiesta per ottenere tutte le conversazioni dell'utente.
// Restituisce una lista di conversazioni.
func (rt *_router) getMyConversationsHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	// Per ora uso un id di esempio solo per fare test
	userID := 1 // DEBUG ricordati di cambiarlo

	rows, err := rt.db.RetrieveConversations(userID)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle conversazioni", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	// Funzione per creare una lista di conversazioni.
	var conversations []Conversation
	for rows.Next() {
		var conversation Conversation
		var participantIDs string
		var participantUsernames string
		err := rows.Scan(&conversation.ID, &conversation.Type, &participantIDs, &participantUsernames)
		if err != nil {
			http.Error(w, "Errore durante la lettura delle conversazioni", http.StatusInternalServerError)
			return
		}

		// Faccio lo split per ottenere la lista con gli id e i nomi dei partecipanti
		participantIDStrs := strings.Split(participantIDs, ",")
		participantUsernameStrs := strings.Split(participantUsernames, ",")

		//fmt.Fprintf(w, "Il valore di participantIDStrs è %s", participantIDStrs) //DEBUG
		//fmt.Fprintf(w, "Il valore di participantUsernameStrs è %s", participantUsernameStrs) //DEBUG

		//Itero sugli id per creare gli oggetti user
		for i, idStr := range participantIDStrs {
			id, _ := strconv.Atoi(idStr)
			conversation.Participants = append(conversation.Participants, User{
				ID:       id,
				Username: participantUsernameStrs[i],
			})
		}

		conversations = append(conversations, conversation)

		//TODO : Aggiungi il latest message
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getConversationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	conversationID, error := strconv.Atoi(ps.ByName("conversationID"))
	if error != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	rows, err := rt.db.RetrieveConversation(conversationID)
	if err != nil {
		http.Error(w, "Errore durante la lettura della conversazione", http.StatusInternalServerError)
		return
	}
	var messages []Message
	for rows.Next() {
		var messageID int
		var conversationID int
		var senderID int
		var messageType string
		var timestamp string
		var status string
		var content string
		var username string

		err := rows.Scan(&messageID, &conversationID, &senderID, &messageType, &timestamp, &status, &content, &username)
		if err != nil {
			http.Error(w, "Errore durante la lettura dei messaggi", http.StatusInternalServerError)
			return
		}
		t, err := time.Parse(time.RFC3339, timestamp)

		messages = append(messages,
			Message{
				ID:   messageID,
				Type: messageType,
				Sender: User{
					ID:       senderID,
					Username: username,
				},
				Timestamp: t,
				Status:    status,
				Content:   content,
			})

	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}

}
