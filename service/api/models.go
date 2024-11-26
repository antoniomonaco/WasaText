package api

import "time"

// Single user in the system.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
}

// Conversation between two users or a group of users.
type Conversation struct {
	ID            int     `json:"id"`
	Type          string  `json:"type"` // "chat" o "group"
	Participants  []User  `json:"participants"`
	LatestMessage Message `json:"latestMessage"`
}

// Message in a conversation.
type Message struct {
	ID        int       `json:"id"`
	Type      string    `json:"type"` // "text" o "media"
	Sender    User      `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
	Status    string    `json:"status"` // "received" o "read"
	Content   string    `json:"content"`
}

// A comment on a message in a conversation.
type Comment struct {
	ID        int       `json:"id"`
	MessageID int       `json:"messageId"`
	Content   string    `json:"content"`
	Sender    User      `json:"sender"`
	Timestamp time.Time `json:"timestamp"`
}
