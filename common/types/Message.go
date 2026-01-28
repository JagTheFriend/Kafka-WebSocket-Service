// Package types Stores Common structs
package types

type Message struct {
	Content    string `json:"content"`
	ChatID     string `json:"chatId"`
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
}
