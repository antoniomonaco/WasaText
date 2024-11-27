package database

import (
	"database/sql"
	"fmt"
)

func (db *appdbimpl) RetrieveConversations(userID int) (*sql.Rows, error) {
	rows, err := db.c.Query(`
WITH tmpConv AS (
	SELECT c.id, c.type, u.id, u.username
	FROM conversations c
	JOIN participants p ON c.id = p.conversation_id
	JOIN users u ON p.user_id = u.id
	WHERE p.user_id = ?
)
SELECT
    c.id,
    c.type,
    GROUP_CONCAT(u.id),
    GROUP_CONCAT(u.username)
FROM tmpConv AS c
JOIN participants AS p ON c.id = p.conversation_id
JOIN users AS u ON p.user_id = u.id
GROUP BY c.id, c.type`, userID)

	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero delle conversazioni: %w", err)
	}
	return rows, nil
}

func (db *appdbimpl) RetrieveConversation(conversationID int) (*sql.Rows, error) {
	rows, err := db.c.Query(`
	SELECT m.* , u.username AS sender
	FROM messages m 
	JOIN users u ON m.sender_id = u.id
	JOIN conversations c ON m.conversation_id = c.id
	WHERE c.id = ?
	`, conversationID)

	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei messaggi: %w", err)
	}
	return rows, nil
}
