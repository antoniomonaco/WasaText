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

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("errore durante il commit: %w", err)
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

	// Inserisco il messaggio nella converszione
	var messageID int64
	result, err := tx.Exec(
		"INSERT INTO messages (conversation_id, sender_id, type, timestamp, status, content) VALUES (?, ?, ?, ?, ?, ?)",
		conversationID, senderID, messageType, timestamp, status, content,
	)
	if err != nil {
		return 0, fmt.Errorf("errore durante l'inserimento del messaggio: %w", err)
	}
	messageID, err = result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("errore durante il recupero dell'ID del messaggio: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return 0, fmt.Errorf("errore durante il commit: %w", err)
	}

	return int(messageID), nil
}

func (db *appdbimpl) DeleteMessage(conversationID int, messageID int) error {
	// Avvia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("errore durante l'avvio della transazione: %w", err)
	}

	_, err = db.c.Exec("DELETE FROM messages WHERE conversation_id = ? AND id = ?",
		conversationID, messageID)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione del messaggio: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("errore durante il commit: %w", err)
	}

	return nil

}

func (db *appdbimpl) GetMessage(conversationID int, messageID int) (*sql.Rows, error) {
	rows, err := db.c.Query(`
	SELECT m.*,u.username
	FROM messages m, users u
	WHERE m.conversation_id = ? AND m.id = ? AND u.id = m.sender_id
`, conversationID, messageID)

	if err != nil {
		return rows, fmt.Errorf("errore durante il recupero del messaggio : %w", err)
	}
	return rows, nil
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

// Serve per verificare che l'utente sia il mittente del messaggiio in modo che possa effettuare
// delle operazioni su di esso
func (db *appdbimpl) IsUserSenderOfMessage(messageID int, senderID int) (bool, error) {
	var count int
	err := db.c.QueryRow(`
        SELECT COUNT(*) 
        FROM messages 
        WHERE id = ? AND sender_id = ?`,
		messageID, senderID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("errore durante la verifica del mittente: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) IsGroup(conversationID int) (bool, error) {
	var count int
	err := db.c.QueryRow(`
        SELECT COUNT(*) 
        FROM conversations 
        WHERE id = ? AND type = "group"`,
		conversationID).Scan(&count)
	if err != nil {
		return false, fmt.Errorf("errore durante la verifica della conversazione: %w", err)
	}
	return count > 0, nil
}

func (db *appdbimpl) AddParticipant(conversationID int, UserID int) error {
	// Avvia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("errore durante l'avvio della transazione: %w", err)
	}

	_, err = db.c.Exec("INSERT INTO participants (user_id, conversation_id) VALUES (?, ?);",
		UserID, conversationID)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiunta del utente al gruppo: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("errore durante il commit: %w", err)
	}

	return nil
}

func (db *appdbimpl) RemoveParticipant(conversationID int, UserID int) error {
	// Avvia una transazione
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("errore durante l'avvio della transazione: %w", err)
	}

	_, err = db.c.Exec("DELETE FROM participants WHERE conversation_id = ? AND user_id = ?",
		conversationID, UserID)
	if err != nil {
		return fmt.Errorf("errore durante l'aggiunta del utente al gruppo: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("errore durante il commit: %w", err)
	}

	return nil
}

func (db *appdbimpl) UpdateGroupName(GroupName string, conversationID int) error {
	_, err := db.c.Exec("UPDATE conversations SET name = ? WHERE id = ?", GroupName, conversationID)

	if err != nil {
		return fmt.Errorf("errore nella modifica del nome : %w", err)
	}

	return nil
}

func (db *appdbimpl) UpdateGroupPhoto(PhotoUrl string, conversationID int) error {
	_, err := db.c.Exec("UPDATE conversations SET photoUrl = ? WHERE id = ?", PhotoUrl, conversationID)

	if err != nil {
		return fmt.Errorf("errore nella modifica della foto profilo : %w", err)
	}

	return nil
}
