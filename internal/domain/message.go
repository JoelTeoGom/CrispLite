package domain

import "time"

type Message struct {
	Content    string    `json:"content"`
	SenderId   string    `json:"senderId"`
	ReceiverId string    `json:"receiverId"`
	Timestamp  time.Time `json:"timestamp"`
}
