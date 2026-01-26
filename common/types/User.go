// Package types Stored user types
package types

type User struct {
	ChatID     string `json:"chatId"`
	SenderID   string `json:"senderId"`
	ReceiverID string `json:"receiverId"`
}
