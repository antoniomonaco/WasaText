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

const chat = "chat"
const group = "group"

// getMyConversations gestisce la richiesta per ottenere tutte le conversazioni dell'utente.
// Restituisce una lista di conversazioni.
func (rt *_router) getMyConversations(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := reqcontext.UserIDFromContext(r.Context())

	rows, err := rt.db.RetrieveConversations(userID)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle conversazioni", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var conversations []Conversation
	for rows.Next() {
		var conversation Conversation
		var participantIDs, participantUsernames, participantPhotos string
		var name, photoUrl sql.NullString

		err := rows.Scan(
			&conversation.ID,
			&conversation.Type,
			&participantIDs,
			&participantUsernames,
			&participantPhotos,
			&name,
			&photoUrl,
		)
		if err != nil {
			http.Error(w, "Errore durante la lettura delle conversazioni", http.StatusInternalServerError)
			return
		}

		conversation.Name = ""
		conversation.PhotoUrl = ""
		conversation.Participants = []User{}

		if name.Valid {
			conversation.Name = name.String
		}
		if photoUrl.Valid {
			conversation.PhotoUrl = photoUrl.String
		}

		if participantIDs != "" && participantUsernames != "" {
			idsList := strings.Split(participantIDs, ",")
			usernamesList := strings.Split(participantUsernames, ",")
			var photosList []string
			if participantPhotos != "" {
				photosList = strings.Split(participantPhotos, "|||")
			}

			for i := range idsList {
				id, _ := strconv.Atoi(idsList[i])
				userPhotoUrl := ""
				if i < len(photosList) {
					userPhotoUrl = photosList[i]
				}

				user := User{
					ID:       id,
					Username: usernamesList[i],
					PhotoUrl: userPhotoUrl,
				}
				conversation.Participants = append(conversation.Participants, user)
			}
		}

		// Gestione chat private
		if conversation.Type == chat && len(conversation.Participants) > 0 {
			for _, participant := range conversation.Participants {
				if participant.ID != userID {
					conversation.Name = participant.Username
					if conversation.PhotoUrl == "" {
						conversation.PhotoUrl = participant.PhotoUrl
					}
					break
				}
			}
		}

		// Inizializza un messaggio vuoto se non ne esistono
		conversation.LatestMessage = Message{}

		// Recupera l'ultimo messaggio solo se esiste
		latestMessageRows, err := rt.db.RetrieveLatestMessage(conversation.ID, userID)
		if err == nil && latestMessageRows != nil {
			defer latestMessageRows.Close()
			if latestMessageRows.Next() {
				if message, err := composeMessage(latestMessageRows, w); err == nil {
					conversation.LatestMessage = message
				}
			}
			if err := latestMessageRows.Err(); err != nil {
				http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
				return
			}
		}
		conversations = append(conversations, conversation)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(conversations); err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) getConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := reqcontext.UserIDFromContext(r.Context())
	conversationID, error := strconv.Atoi(ps.ByName("conversationID"))
	if error != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	messageRows, err := rt.db.RetrieveMessages(conversationID, userID)
	if err != nil {
		http.Error(w, "Errore durante la lettura della conversazione", http.StatusInternalServerError)
		return
	}
	defer messageRows.Close()
	var messages []Message
	for messageRows.Next() {
		message, err := composeMessage(messageRows, w)
		if err != nil {
			http.Error(w, "Errore durante la lettura dei messaggi", http.StatusInternalServerError)
			return
		}
		messages = append(messages, message)
	}
	if err := messageRows.Err(); err != nil {
		http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
		return
	}
	/*
		if messages == nil {
			http.Error(w, "Conversazione non trovata", http.StatusNotFound)
			return
		}
	*/
	infoRows, err := rt.db.RetrieveConversationInfo(conversationID)
	if err != nil {
		http.Error(w, "Errore durante il recupero delle informazioni della conversazione", http.StatusInternalServerError)
		return
	}
	var conversation Conversation
	var name, photoUrl, participantIDs, participantUsernames, participantPhotoUrls sql.NullString
	err = infoRows.Scan(&conversation.ID, &conversation.Type, &name, &photoUrl, &participantIDs, &participantUsernames, &participantPhotoUrls)
	if err != nil {
		http.Error(w, "Errore durante la lettura delle info della conversazione", http.StatusInternalServerError)
		return
	}
	if photoUrl.Valid {
		conversation.PhotoUrl = photoUrl.String
	}
	if name.Valid {
		conversation.Name = name.String
	}
	participantIDStrs := []string{}
	if participantIDs.Valid {
		participantIDStrs = strings.Split(participantIDs.String, ",")
	}
	participantUsernameStrs := []string{}
	if participantUsernames.Valid {
		participantUsernameStrs = strings.Split(participantUsernames.String, ",")
	}
	participantPhotoUrlStrs := []string{}
	if participantPhotoUrls.Valid {
		participantPhotoUrlStrs = strings.Split(participantPhotoUrls.String, "|||")
	}
	for i, idStr := range participantIDStrs {
		id, _ := strconv.Atoi(idStr)
		userphoto := ""
		if len(participantPhotoUrlStrs) > i {
			userphoto = participantPhotoUrlStrs[i]
		}
		conversation.Participants = append(conversation.Participants, composeUser(id, participantUsernameStrs[i], userphoto))
	}
	response := struct {
		Conversation Conversation `json:"conversation"`
		Messages     []Message    `json:"messages"`
	}{
		Conversation: conversation,
		Messages:     messages,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) createConversation(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())

	// Decodifica del payload JSON
	var request struct {
		Type         string `json:"type"`         // chat o "group"
		Participants []int  `json:"participants"` // Array di ID utenti
		Name         string `json:"name"`         // Nome per il gruppo (opzionale)
		PhotoUrl     string `json:"photoUrl"`     // URL immagine (opzionale)
	}
	err := json.NewDecoder(r.Body).Decode(&request)

	if err != nil {
		http.Error(w, "Errore nella decodifica della richiesta", http.StatusBadRequest)
		return
	}

	if request.Type != chat && request.Type != group {
		http.Error(w, "Devi scegliere un tipo tra 'chat' o 'group'", http.StatusBadRequest)
		return
	}

	if request.Name == "" && request.Type == group {
		http.Error(w, "Devi impostare un nome per il gruppo", http.StatusBadRequest)
		return
	}

	if len(request.Participants) == 0 {
		http.Error(w, "Devi includere almeno un partecipante", http.StatusBadRequest)
		return
	}

	if request.Name != "" && request.Type == chat {
		http.Error(w, "Non puoi impostare il nome ad una chat privata", http.StatusBadRequest)
		return
	}

	if request.PhotoUrl != "" && request.Type == chat {
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

	if request.Type == chat && len(participants) != 2 {
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
	response := struct {
		ConversationID int `json:"conversation_id"`
	}{ConversationID: conversationID}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) sendMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())

	conversationID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, IDFromContext)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}

	var request struct {
		Type    string `json:"type"`
		Content string `json:"content"`
		ReplyTo int    `json:"replyTo"`
	}

	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Errore nella decodifica del messaggio", http.StatusBadRequest)
		return
	}

	if request.Type != "media" && request.Type != "text" {
		http.Error(w, "Il messaggio deve essere di tipo 'text' o 'media'", http.StatusBadRequest)
		return
	}

	/*
		var replyedMessage Message

		if request.ReplyTo > 0 {
			rows, err := rt.db.GetMessage(conversationID, request.ReplyTo)
			if err != nil {
				http.Error(w, "Errore nel recupero del messaggio di risposta", http.StatusInternalServerError)
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			defer rows.Close()

			for rows.Next() {
				replyedMessage, err = composeMessage(rows, w)
				if err != nil {
					http.Error(w, "Errore nella composizione del messaggio", http.StatusInternalServerError)
					http.Error(w, err.Error(), http.StatusInternalServerError)
					return
				}
			}
			if err := rows.Err(); err != nil {
				http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
				return
			}


		}

	*/

	timestamp := time.Now()
	messageStatus := "received"

	messageID, err := rt.db.SendMessage(conversationID, IDFromContext, request.Type, timestamp, messageStatus, request.Content, request.ReplyTo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rows, err := rt.db.GetMessage(conversationID, messageID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		http.Error(w, "Messaggio non trovato dopo la creazione", http.StatusInternalServerError)
		return
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
		return
	}

	message, err := composeMessage(rows, w)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	/*

		if request.ReplyTo > 0 {
			response := struct {
				ID        int       `json:"id"`
				Type      string    `json:"type"`
				Sender    User      `json:"sender"`
				Timestamp time.Time `json:"timestamp"`
				Status    string    `json:"status"`
				Content   string    `json:"content"`
				ReplyTo   Message   `json:"replyTo"`
			}{
				ID:        conversationID,
				Type:      message.Type,
				Sender:    message.Sender,
				Timestamp: message.Timestamp,
				Status:    message.Status,
				Content:   message.Content,
				ReplyTo:   replyedMessage,
			}

			// In questa response includo direttamente il messaggio a cui si risponde, se c'è, così facilito il frontend

			err = json.NewEncoder(w).Encode(response)
			if err != nil {
				http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
				return
			}
		}
		err = json.NewEncoder(w).Encode(message)
		if err != nil {
			http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
			return
		}
	*/
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func (rt *_router) deleteMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())

	conversationID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	messageID, err := strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}

	// Verifico che l'utente sia un partecipante della conversazione.
	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}

	isSender, err := rt.db.IsUserSenderOfMessage(messageID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica del mittente", http.StatusInternalServerError)
		return
	}
	if !isSender {
		http.Error(w, "Operazione non consentita su questo messaggio", http.StatusForbidden) // 403
		return
	}

	err = rt.db.DeleteMessage(conversationID, messageID)
	if err != nil {
		http.Error(w, "Messaggio inesistente", http.StatusNotFound) // 404
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNoContent) // 204
}

func (rt *_router) forwardMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	IDFromContext := reqcontext.UserIDFromContext(r.Context())
	conversationID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	// Verifico che l'utente sia un partecipante della conversazione.
	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}
	messageID, err := strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "ID messaggio non valido", http.StatusBadRequest)
		return
	}
	var request struct {
		TargetConversationID int `json:"id"`
	}
	err = json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}
	// Verifico che l'utente sia un partecipante della conversazione.
	isParticipant, err = rt.db.IsUserParticipantOfConversation(request.TargetConversationID, IDFromContext)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden) // 403
		return
	}
	rows, err := rt.db.GetMessage(conversationID, messageID)
	if err != nil {
		http.Error(w, "Messaggio non trovato", http.StatusNotFound) // 404
		return
	}
	var message Message
	for rows.Next() {
		message, err = composeMessage(rows, w)
		if err != nil {
			http.Error(w, "Errore durante la composizione del messaggio", http.StatusInternalServerError)
			return
		}
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
		return
	}

	if err != nil {
		http.Error(w, "Errore nella composizione del messaggio", http.StatusInternalServerError)
		return
	}
	message.Timestamp = time.Now()
	message.Status = "received"
	newMessageID, err := rt.db.SendMessage(request.TargetConversationID, IDFromContext, message.Type, message.Timestamp, message.Status, message.Content, 0)
	if err != nil {
		http.Error(w, "Errore nell'invio del messaggio", http.StatusInternalServerError)
		return
	}

	// Imposto i nuovi valori nel messaggio per poterlo restituire

	message.ID = newMessageID
	newSender, err := composeUserFromID(IDFromContext, rt.db)
	if err != nil {
		http.Error(w, "Errore nella creazione dell'utente", http.StatusInternalServerError)
		return
	}
	message.Sender = newSender

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated) // 201
	err = json.NewEncoder(w).Encode(message)
	if err != nil {
		http.Error(w, "Errore nella codifica della risposta", http.StatusInternalServerError)
		return
	}
}

func composeMessage(rows *sql.Rows, w http.ResponseWriter) (Message, error) {
	var message Message
	var messageID, conversationID, senderID int
	var messageType, timestamp, status, content, username string
	var replyToID int

	err := rows.Scan(
		&messageID,
		&conversationID,
		&senderID,
		&messageType,
		&timestamp,
		&status,
		&content,
		&replyToID,
		&username,
	)
	if err != nil {
		http.Error(w, "errore durante lo scan del messaggio", http.StatusInternalServerError)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return message, err
	}

	parsedTime, err := time.Parse("2006-01-02 15:04:05", timestamp)
	if err != nil {
		http.Error(w, "errore durante il parsing del timestamp", http.StatusInternalServerError)
		return message, err
	}
	parsedTime = parsedTime.Local()

	message = Message{
		ID:   messageID,
		Type: messageType,
		Sender: User{
			ID:       senderID,
			Username: username,
		},
		Timestamp: parsedTime,
		Status:    status,
		Content:   content,
		ReplyTo:   replyToID,
	}

	return message, nil
}
func (rt *_router) markMessagesAsRead(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := reqcontext.UserIDFromContext(r.Context())
	conversationID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "ID conversazione non valido", http.StatusBadRequest)
		return
	}

	// Verifica che l'utente sia un partecipante della conversazione
	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, userID)
	if err != nil {
		http.Error(w, "Errore durante la verifica dei partecipanti", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Accesso non autorizzato alla conversazione", http.StatusForbidden)
		return
	}

	// Segna i messaggi come letti
	err = rt.db.MarkMessagesAsRead(conversationID, userID)
	if err != nil {
		http.Error(w, "Errore durante l'aggiornamento dello stato dei messaggi", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent) // 204 No Content
}
func (rt *_router) getMessage(w http.ResponseWriter, r *http.Request, ps httprouter.Params, ctx reqcontext.RequestContext) {
	userID := reqcontext.UserIDFromContext(r.Context())
	conversationID, err := strconv.Atoi(ps.ByName("conversationID"))
	if err != nil {
		http.Error(w, "Invalid conversation ID", http.StatusBadRequest)
		return
	}

	messageID, err := strconv.Atoi(ps.ByName("messageID"))
	if err != nil {
		http.Error(w, "Invalid message ID", http.StatusBadRequest)
		return
	}

	isParticipant, err := rt.db.IsUserParticipantOfConversation(conversationID, userID)
	if err != nil {
		http.Error(w, "Error verifying participation", http.StatusInternalServerError)
		return
	}
	if !isParticipant {
		http.Error(w, "Unauthorized access", http.StatusForbidden)
		return
	}

	rows, err := rt.db.GetMessage(conversationID, messageID)
	if err != nil {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	defer rows.Close()

	if !rows.Next() {
		http.Error(w, "Message not found", http.StatusNotFound)
		return
	}
	if err := rows.Err(); err != nil {
		http.Error(w, "Errore durante l'iterazione delle righe", http.StatusInternalServerError)
		return
	}

	message, err := composeMessage(rows, w)
	if err != nil {
		http.Error(w, "Error composing message", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(message); err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}
