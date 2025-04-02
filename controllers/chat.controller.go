package controllers

import (
	"log"
	"net/http"
	"time"

	"github.com/SwafaX/swafa-backend/models"
	socketio "github.com/doquangtan/socket.io/v4"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ChatController struct {
	DB *gorm.DB
	IO *socketio.Io
}

func NewChatController(DB *gorm.DB, io *socketio.Io) ChatController {
	return ChatController{DB: DB, IO: io}
}

func (cc *ChatController) SetupSocketHandlers() {
	cc.IO.OnConnection(func(socket *socketio.Socket) {
		log.Printf("User connected: %s", socket.Id)

		// Join conversation handler
		socket.On("joinConversation", func(event *socketio.EventPayload) {
			if len(event.Data) > 0 {
				chatID := event.Data[0].(string)
				socket.Join(chatID)
				log.Printf("User %s joined conversation: %s", socket.Id, chatID)
			}
		})

		// Send message handler
		socket.On("sendMessage", func(event *socketio.EventPayload) {
			if len(event.Data) < 1 {
				return
			}

			// Add type checking and conversion
			var messageData map[string]interface{}
			switch data := event.Data[0].(type) {
			case string:
				// Handle string data (likely JSON)
				log.Printf("Received string data: %v", data)
				return
			case map[string]interface{}:
				messageData = data
			default:
				log.Printf("Unexpected data type: %T", data)
				return
			}

			chatID, ok1 := messageData["chat_id"].(string)
			senderID, ok2 := messageData["sender_id"].(string)
			content, ok3 := messageData["content"].(string)

			if !ok1 || !ok2 || !ok3 {
				log.Printf("Invalid message format")
				return
			}

			// Save message to database
			message, err := cc.saveMessage(chatID, senderID, content)
			if err != nil {
				log.Printf("Failed to save message: %v", err)
				return
			}

			// Broadcast to conversation room
			cc.IO.To(chatID).Emit("newMessage", message)

			// Notify all clients about conversation update
			cc.IO.Emit("conversationUpdated", map[string]interface{}{
				"chatId":          chatID,
				"lastMessage":     message.Content,
				"lastMessageTime": message.CreatedAt,
			})
		})

		// Disconnect handler
		socket.On("disconnect", func(event *socketio.EventPayload) {
			log.Printf("User disconnected: %s", socket.Id)
		})
	})
}

// GetChatMessages retrieves messages for a specific conversation (keep existing HTTP endpoint)
func (cc *ChatController) GetChatMessages(ctx *gin.Context) {
	chatID := ctx.Param("chatId")

	var messages []models.Message
	result := cc.DB.Where("chat_id = ?", chatID).Order("created_at asc").Find(&messages)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

// Database operations
func (cc *ChatController) saveMessage(chatID, senderID, content string) (*models.Message, error) {
	chatUUID, err := uuid.Parse(chatID)
	if err != nil {
		return nil, err
	}

	senderUUID, err := uuid.Parse(senderID)
	if err != nil {
		return nil, err
	}

	message := models.Message{
		ID:        uuid.New(),
		ChatID:    chatUUID,
		SenderID:  senderUUID,
		Content:   content,
		CreatedAt: time.Now(),
	}

	result := cc.DB.Create(&message)
	if result.Error != nil {
		return nil, result.Error
	}

	return &message, nil
}
