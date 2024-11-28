package api

import "time"

// Single user in the system.
type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	PhotoUrl string `json:"photoUrl"`
}

// Conversation between two users or a group of users.
type Conversation struct {
	ID            int     `json:"id"`
	Type          string  `json:"type"` // "chat" o "group"
	Participants  []User  `json:"participants"`
	LatestMessage Message `json:"latestMessage"`
	Name          string  `json:"name"`     // In case of a private chat it will use the name of the other user
	PhotoUrl      string  `json:"photoUrl"` // In case of a private chat it will use the photo of the other user
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
