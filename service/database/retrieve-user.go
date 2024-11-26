package database

import (
	"database/sql"
	"fmt"
)

// RetrieveUser cerca un utente nel database in base al nome utente.
// Se l'utente non esiste, lo crea e restituisce il suo ID.
func (db *appdbimpl) RetrieveUser(username string) (int, error) {
	var userID int
	// Prova a trovare l'utente nel database
	err := db.c.QueryRow("SELECT id FROM users WHERE username = ?", username).Scan(&userID)
	if err != nil {
		if err == sql.ErrNoRows {
			// Se l'utente non esiste, lo creiamo
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
