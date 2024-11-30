package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *appdbimpl) CreateConversation(conversationType, name, photoUrl string, participants []int) (int, error) {
	// Avvia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return 0, fmt.Errorf("errore durante l'avvio della transazione: %w", err)
	}

	// Gestisco eventuali errori in modo da evitare di "sporcare" il database

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Inserisco la conversazione
	var conversationID int

	err = tx.QueryRow(
		"INSERT INTO conversations (type, name, photoUrl) VALUES (?, ?, ?) RETURNING id",
		conversationType, name, photoUrl,
	).Scan(&conversationID)
	if err != nil {
		return 0, fmt.Errorf("errore durante l'inserimento della conversazione: %w", err)
	}

	// Inserisco i partecipanti
	for _, participant := range participants {
		_, err := tx.Exec(
			"INSERT INTO participants (user_id, conversation_id) VALUES (?, ?)",
			participant, conversationID,
		)
		if err != nil {
			return 0, fmt.Errorf("errore durante l'inserimento dei partecipanti: %w", err)
		}
	}

	return conversationID, nil
}

func (db *appdbimpl) RetrieveConversations(userID int) (*sql.Rows, error) {
	rows, err := db.c.Query(`
WITH tmpConv AS (
	SELECT c.id, c.type, u.id, u.username, c.name, c.photoUrl
	FROM conversations c
	JOIN participants p ON c.id = p.conversation_id
	JOIN users u ON p.user_id = u.id
	WHERE p.user_id = ?
)
SELECT
    c.id,
    c.type,
    GROUP_CONCAT(u.id),
    GROUP_CONCAT(u.username),
    COALESCE(c.name, ''),    -- COALESCE CONVERTE I VALORI NULL IN STRINGHE VUOTE
    COALESCE(c.photoUrl, '')
FROM tmpConv AS c
JOIN participants AS p ON c.id = p.conversation_id
JOIN users AS u ON p.user_id = u.id
GROUP BY c.id, c.type`, userID)

	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero delle conversazioni: %w", err)
	}
	return rows, nil
}

func (db *appdbimpl) RetrieveConversation(conversationID int, userID int) (*sql.Rows, error) {
	rows, err := db.c.Query(`
		SELECT 
			m.*, 
			u.username AS sender 
		FROM messages m 
		JOIN users u ON m.sender_id = u.id 
		JOIN conversations c ON m.conversation_id = c.id 
		JOIN participants p ON c.id = p.conversation_id 
		WHERE c.id = ? AND p.user_id = ?  
	`, conversationID, userID)

	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei messaggi: %w", err)
	}
	return rows, nil
}

func (db *appdbimpl) RetrieveLatestMessage(conversationID int, userID int) (*sql.Rows, error) {
	rows, err := db.c.Query(`
	SELECT 
		m.*, 
		u.username AS sender 
	FROM messages m 
	JOIN users u ON m.sender_id = u.id 
	JOIN conversations c ON m.conversation_id = c.id 
	JOIN participants p ON c.id = p.conversation_id 
	WHERE c.id = ? AND p.user_id = ?
	  AND m.timestamp = (
      SELECT MAX(m2.timestamp)
      FROM messages m2
      WHERE m2.conversation_id = c.id
  )  
`, conversationID, userID)

	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dell'anteprima: %w", err)
	}
	return rows, nil
}

func (db *appdbimpl) SendMessage(conversationID int, senderID int, messageType string, timestamp time.Time, status string, content string) (int, error) {
	// Avvia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return 0, fmt.Errorf("errore durante l'avvio della transazione: %w", err)
	}

	// Gestisco eventuali errori in modo da evitare di "sporcare" il database

	defer func() {
		if p := recover(); p != nil {
			tx.Rollback()
			panic(p)
		} else if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	// Inserisco il messaggio nella converszione
	var messageID int
	err = tx.QueryRow(
		"INSERT INTO messages (conversation_id, sender_id, type, timestamp, status, content) VALUES (?, ?, ?, ?, ?, ?) RETURNING id",
		conversationID, senderID, messageType, timestamp, status, content).Scan(&messageID)
	if err != nil {
		return 0, fmt.Errorf("errore durante l'inserimento del messaggio: %w", err)
	}

	return messageID, nil
}

// Serve per verificare che l'utente faccia effettivamente parte della conversazione in modo che possa effettuare
// delle operazioni su di essa
func (db *appdbimpl) IsUserParticipantOfConversation(conversationID int, userID int) (bool, error) {
	var count int
	err := db.c.QueryRow(`
        SELECT COUNT(*) 
        FROM participants 
        WHERE conversation_id = ? AND user_id = ?`,
		conversationID, userID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("errore durante la verifica della partecipazione: %w", err)
	}
	return count > 0, nil
}
