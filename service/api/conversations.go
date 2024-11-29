package api

import (
	"database/sql"
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

	userID := reqcontext.UserIDFromContext(r.Context())

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
		var name sql.NullString
		var photoUrl sql.NullString

		err := rows.Scan(&conversation.ID, &conversation.Type, &participantIDs, &participantUsernames, &name, &photoUrl)
		if err != nil {
			http.Error(w, "Errore durante la lettura delle conversazioni", http.StatusInternalServerError)
			return
		}

		// Converto i campi da sql.NullString a stringa normale
		conversation.Name = ""
		if name.Valid {
			conversation.Name = name.String
		}

		conversation.PhotoUrl = ""
		if photoUrl.Valid {
			conversation.PhotoUrl = photoUrl.String
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

		// Compongo l'anteprima della conversazione

		rows, err := rt.db.RetrieveLatestMessage(conversation.ID, userID)
		if err != nil {
			http.Error(w, "Errore durante la lettura dell'anteprima", http.StatusInternalServerError)
			return
		}
		var message Message
		for rows.Next() {
			message, err = composeMessage(rows, w)
			if err != nil {
				http.Error(w, "Errore durante la composizione dell'anteprima", http.StatusInternalServerError)
				return
			}
		}

		conversation.LatestMessage = message

		conversations = append(conversations, conversation)

	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(conversations)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getConversationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := reqcontext.UserIDFromContext(r.Context())
	conversationID, error := strconv.Atoi(ps.ByName("conversationID"))
	if error != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	rows, err := rt.db.RetrieveConversation(conversationID, userID)
	if err != nil {
		http.Error(w, "Errore durante la lettura della conversazione", http.StatusInternalServerError)
		return
	}
	var messages []Message

	for rows.Next() {
		message, err := composeMessage(rows, w)
		if err != nil {
			http.Error(w, "Errore durante la lettura dei messaggi", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(messages)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}

}

func (rt *_router) createConversationHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())

	// Decodifica del payload JSON
	var request struct {
		Type         string `json:"type"`         // "chat" o "group"
		Participants []int  `json:"participants"` // Array di ID utenti
		Name         string `json:"name"`         // Nome per il gruppo (opzionale)
		PhotoUrl     string `json:"photoUrl"`     // URL immagine (opzionale)
	}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Errore nella decodifica della richiesta", http.StatusBadRequest)
		return
	}

	if request.Type != "chat" && request.Type != "group" {
		http.Error(w, "Devi scegliere un tipo tra 'chat' o 'group'", http.StatusBadRequest)
		return
	}

	if request.Name == "" && request.Type == "group" {
		http.Error(w, "Devi impostare un nome per il gruppo", http.StatusBadRequest)
		return
	}

	if len(request.Participants) == 0 {
		http.Error(w, "Devi includere almeno un partecipante", http.StatusBadRequest)
		return
	}

	if request.Name != "" && request.Type == "chat" {
		http.Error(w, "Non puoi impostare il nome ad una chat privata", http.StatusBadRequest)
		return
	}

	if request.PhotoUrl != "" && request.Type == "chat" {
		http.Error(w, "Non puoi impostare la foto ad una chat privata", http.StatusBadRequest)
		return
	}

	// Mi assicuro che l'utente corrente sia incluso nei partecipanti
	participantsSet := make(map[int]struct{})
	for _, p := range request.Participants {
		participantsSet[p] = struct{}{}
	}
	participantsSet[IDFromContext] = struct{}{}

	// Converto il set in una slice
	participants := make([]int, 0, len(participantsSet))
	for p := range participantsSet {
		participants = append(participants, p)
	}

	if request.Type == "chat" && len(participants) != 2 {
		http.Error(w, "Una chat privata deve avere esattamente due partecipanti", http.StatusBadRequest)
		return
	}

	// Creazione della conversazione
	conversationID, err := rt.db.CreateConversation(request.Type, request.Name, request.PhotoUrl, participants)
	if err != nil {
		http.Error(w, "Errore durante la creazione della conversazione", http.StatusInternalServerError)
		return
	}

	// Risposta 201 Created
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(struct {
		ConversationID int `json:"conversation_id"`
	}{ConversationID: conversationID})
}

func composeMessage(rows *sql.Rows, w http.ResponseWriter) (Message, error) {
	var message Message

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
		return message, err
	}
	t, _ := time.Parse(time.RFC3339, timestamp)

	message = Message{
		ID:   messageID,
		Type: messageType,
		Sender: User{
			ID:       senderID,
			Username: username,
		},
		Timestamp: t,
		Status:    status,
		Content:   content,
	}

	return message, nil
}
