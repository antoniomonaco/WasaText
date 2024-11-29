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

func (db *appdbimpl) UpdateUserName(username string, IDFromContext int) error {

	// Controllo se il nome utente esiste già
	var existingID int
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&existingID)
	if err == nil && existingID != IDFromContext {
		return fmt.Errorf("il nome utente '%s' è già in uso", username)
	}
	if err != nil && err != sql.ErrNoRows {
		return fmt.Errorf("errore durante il controllo del nome utente: %w", err)
	}

	// Aggiorno il nome utente
	_, err = db.c.Exec("UPDATE users SET username = ? WHERE id = ?", username, IDFromContext)

	if err != nil {
		return fmt.Errorf("errore nella modifica del nome : %w", err)
	}

	return nil
}
