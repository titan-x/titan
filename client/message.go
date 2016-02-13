package client

import "time"

// Message is a chat message.
type Message struct {
	From    string    `json:"from,omitempty"`
	To      string    `json:"to,omitempty"`
	Time    time.Time `json:"time"`
	Message string    `json:"message"`
}
