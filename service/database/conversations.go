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

/*
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
GROUP BY c.id, c.type , userID
*/

/*
SELECT
    c.id,
    c.type,
    GROUP_CONCAT(u.id),
    GROUP_CONCAT(u.username)
FROM conversations AS c
JOIN participants AS p ON c.id = p.conversation_id
JOIN users AS u ON p.user_id = u.id
GROUP BY c.id, c.type
*/
