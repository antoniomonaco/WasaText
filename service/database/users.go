package database

import (
	"database/sql"
	"fmt"
)

// Login cerca un utente nel database in base al nome utente.
// Se l'utente non esiste, lo crea e restituisce il suo ID.
func (db *appdbimpl) Login(username string) (int, error) {
	var userID int
	// cerco l'utente nel database
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se l'utente non esiste, lo creo
			err = db.c.QueryRow("INSERT INTO users (username) VALUES (?) RETURNING id", username).Scan(&userID)
			if err != nil {
				return 0, fmt.Errorf("errore durante la creazione dell'utente: %w", err)
			}
		} else {
			// Gestisce altri errori di database
			return 0, fmt.Errorf("errore durante il controllo dell'utente: %w", err)
		}
	}
	// Restituisce l'userID che può essere utilizzato per l'accesso
	return userID, nil
}

// RetrieveUsers cerca un utente nel database in base al nome utente.
func (db *appdbimpl) RetrieveUsers(username string, IDFromContext int) (*sql.Rows, error) {

	var rows *sql.Rows
	var err error

	if username != "" {
		rows, err = db.c.Query("SELECT * FROM users WHERE username = ? AND id != ?", username, IDFromContext)
		if err != nil {
			return rows, fmt.Errorf("nessun utente trovato con questo nome : %w", err)
		}
	} else {
		rows, err = db.c.Query("SELECT * FROM users WHERE id != ? ", IDFromContext)
		if err != nil {
			return rows, fmt.Errorf("nessun utente trovato : %w", err)
		}
	}

	return rows, nil
}
