package models

import (
	"gofiber-chat-api/src/configs"
	"time"

	"gorm.io/gorm"
)

type Chat struct {
	Base
	ChatUsers   []ChatUser `json:"users" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	LastMessage Message    `json:"last_message"`
	Status      string     `json:"status"` // NONE, REPLY
}

type ChatUser struct {
	Base
	ChatID string `gorm:"type:uuid;" json:"chat_id"`
	UserID string `gorm:"type:uuid;" json:"user_id"`
	User   User   `gorm:"foreignKey:UserID" json:"user"`
}

type Message struct {
	Base
	ChatID string    `gorm:"type:uuid;" json:"chat_id"`
	UserID string    `gorm:"type:uuid;" json:"user_id"`
	User   User      `gorm:"foreignKey:UserID" json:"user"`
	Body   string    `json:"body"`
	Status string    `gorm:"default:NONE" json:"status"` // NONE, SENT, UNSEEN, SEEN
	SeenAt time.Time `gorm:"default:NULL" json:"seen_at"`
}

func SelectChatsbyUserID(user_id *string) []*Chat {
	var chats []*Chat

	configs.DB.Preload("LastMessage", func(db *gorm.DB) *gorm.DB {
		return db.Preload("User").Order("created_at DESC").Limit(1)
	}).Preload("ChatUsers.User").Joins("INNER JOIN messages ON messages.chat_id = chats.id AND messages.user_id IN (?)", user_id).
		Group("chats.id").Find(&chats)

	return chats
}

func SelectChatbyID(id *string) *Chat {
	var chat *Chat
	configs.DB.First(&chat, "id = ?", &id)
	return chat
}

func SelectChatUserbyChatUserID(chat_id, user_id *string) *ChatUser {
	var chatMember *ChatUser
	configs.DB.First(&chatMember, "chat_id = ? AND user_id = ?", &chat_id, &user_id)
	return chatMember
}

func CreateChat(chat *Chat) error {
	result := configs.DB.Create(&chat)
	return result.Error
}

func CreateMessage(message *Message) error {
	result := configs.DB.Create(&message)
	return result.Error
}
