package domain

type Message struct {
	ConversationId string `json:"conversationId"`
	SenderId       string `json:"senderId"`
	ReceiverId     string `json:"receiverId"`
	Timestamp      string `json:"timestamp"`
	Content        string `json:"content"`
}
