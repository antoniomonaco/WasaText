package database

import (
	"database/sql"
	"fmt"
	"time"
)

func (db *appdbimpl) AddComment(messageID int, userID int, content string, timestamp time.Time) (int, error) {
	var commentID int
	query := `
        INSERT INTO comments (message_id, sender_id, content, timestamp)
        VALUES (?, ?, ?, ?)
        RETURNING id
    `
	err := db.c.QueryRow(query, messageID, userID, content, timestamp).Scan(&commentID)
	if err != nil {
		return 0, fmt.Errorf("errore durante l'aggiunta del commento: %w", err)
	}
	return commentID, nil
}

func (db *appdbimpl) GetComments(messageID int) (*sql.Rows, error) {
	query := `
        SELECT c.id, c.sender_id,c.message_id, c.content, c.timestamp
        FROM comments c
        WHERE c.message_id = ?
        ORDER BY c.timestamp ASC
    `
	rows, err := db.c.Query(query, messageID)
	if err != nil {
		return nil, fmt.Errorf("errore durante il recupero dei commenti: %w", err)
	}
	return rows, nil
}

func (db *appdbimpl) DeleteComment(messageID int, commentID int, userID int) error {
	// Inizio della transazione
	tx, err := db.c.Begin()
	if err != nil {
		return fmt.Errorf("errore durante l'avvio della transazione: %w", err)
	}

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

	// Controllo che il commento appartenga all'utente
	var count int
	err = tx.QueryRow(
		"SELECT COUNT(*) FROM comments WHERE id = ? AND message_id = ? AND sender_id = ?",
		commentID, messageID, userID,
	).Scan(&count)
	if err != nil {
		return fmt.Errorf("errore durante la verifica del commento: %w", err)
	}

	if count == 0 {
		return fmt.Errorf("commento non trovato o non autorizzato")
	}

	// Eliminazione del commento
	_, err = tx.Exec(
		"DELETE FROM comments WHERE id = ? AND message_id = ? AND sender_id = ?",
		commentID, messageID, userID,
	)
	if err != nil {
		return fmt.Errorf("errore durante l'eliminazione del commento: %w", err)
	}

	return nil
}
