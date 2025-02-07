/*
Package database is the middleware between the app database and the code. All data (de)serialization (save/load) from a
persistent database are handled here. Database specific logic should never escape this package.

To use this package you need to apply migrations to the database if needed/wanted, connect to it (using the database
data source name from config), and then initialize an instance of AppDatabase from the DB connection.

For example, this code adds a parameter in `webapi` executable for the database data source name (add it to the
main.WebAPIConfiguration structure):

	DB struct {
		Filename string `conf:""`
	}

This is an example on how to migrate the DB and connect to it:

	// Start Database
	logger.Println("initializing database support")
	db, err := sql.Open("sqlite3", "./foo.db")
	if err != nil {
		logger.WithError(err).Error("error opening SQLite DB")
		return fmt.Errorf("opening SQLite: %w", err)
	}
	defer func() {
		logger.Debug("database stopping")
		_ = db.Close()
	}()

Then you can initialize the AppDatabase and pass it to the api package.
*/
package database

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

// AppDatabase is the high level interface for the DB
type AppDatabase interface {
	GetName() (string, error)
	SetName(name string) error

	Ping() error

	Login(username string) (int, error)

	RetrieveUsers(username string, IDFromContext int) (*sql.Rows, error)
	RetrieveUserFromID(userID int) (*sql.Row, error)
	UpdateUserName(username string, IDFromContext int) error
	SetUserPhoto(photoUrl string, IDFromContext int) error

	CreateConversation(conversationType, name, photoUrl string, participants []int) (int, error)
	RetrieveConversations(userID int) (*sql.Rows, error)
	RetrieveMessages(conversationID int, userID int) (*sql.Rows, error)
	RetrieveConversationInfo(conversationID int) (*sql.Row, error)
	RetrieveLatestMessage(conversationID int, userID int) (*sql.Rows, error)
	IsUserParticipantOfConversation(conversationID int, userID int) (bool, error)
	IsUserSenderOfMessage(conversationID int, userID int) (bool, error)
	SendMessage(conversationID int, IDFromContext int, messageType string, timestamp time.Time, status string, content string, replyTo int) (int, error)
	DeleteMessage(conversationID int, messageID int) error
	GetMessage(conversationID int, messageID int) (*sql.Rows, error)
	IsGroup(conversationID int) (bool, error)
	AddParticipant(conversationID int, UserID int) error
	RemoveParticipant(conversationID int, UserID int) error
	UpdateGroupName(GroupName string, conversationID int) error
	UpdateGroupPhoto(PhotoUrl string, conversationID int) error
	MarkMessagesAsRead(conversationID, userID int) error

	AddComment(messageID int, userID int, content string, timestamp time.Time) (int, error)
	GetComments(messageID int) (*sql.Rows, error)
	DeleteComment(messageID int, commentID int, userID int) error
}

type appdbimpl struct {
	c *sql.DB
}

// New returns a new instance of AppDatabase based on the SQLite connection `db`.
// `db` is required - an error will be returned if `db` is `nil`.
func New(db *sql.DB) (AppDatabase, error) {
	if db == nil {
		return nil, errors.New("database is required when building a AppDatabase")

	}

	// Check if table exists. If not, the database is empty, and we need to create the structure
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS users (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT UNIQUE NOT NULL,
			photoUrl TEXT
		);

		CREATE TABLE IF NOT EXISTS conversations (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			type TEXT NOT NULL CHECK (type IN ('chat', 'group')),
			name TEXT,
			photoUrl TEXT 
		);

		CREATE TABLE IF NOT EXISTS participants (
			user_id INTEGER REFERENCES users(id),
			conversation_id INTEGER REFERENCES conversations(id),
			PRIMARY KEY (user_id, conversation_id)
		);

		CREATE TABLE IF NOT EXISTS messages (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			conversation_id INTEGER REFERENCES conversations(id),
			sender_id INTEGER REFERENCES users(id),
			type TEXT NOT NULL CHECK (type IN ('text', 'media')),
			timestamp TEXT NOT NULL,
			status TEXT NOT NULL CHECK (status IN ('received', 'read')),
			content TEXT,
			reply_to INTEGER REFERENCES messages(id)
		);

		CREATE TABLE IF NOT EXISTS comments (
			id INTEGER PRIMARY KEY AUTOINCREMENT,
			message_id INTEGER REFERENCES messages(id),
			sender_id INTEGER REFERENCES users(id),
			content TEXT NOT NULL,
			timestamp TEXT NOT NULL
		);
	`)
	if err != nil {
		return nil, fmt.Errorf("error creating database structure: %w", err)
	}

	return &appdbimpl{
		c: db,
	}, nil
}

func (db *appdbimpl) Ping() error {
	return db.c.Ping()
}
