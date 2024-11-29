package database

import (
	"database/sql"
	"fmt"
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
	fmt.Printf("Inserting conversation: type=%s, name=%s, photoUrl=%s\n", conversationType, name, photoUrl)

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
