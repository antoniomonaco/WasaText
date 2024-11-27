package api

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

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
