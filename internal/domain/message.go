package domain

type Message struct {
	ConversationId string
	SenderId       string
	ReceiverId     string
	Timestamp      string
	Content        string
}
