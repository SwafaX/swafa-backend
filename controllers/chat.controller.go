package controllers

import (
	"net/http"
	"time"

	"github.com/SwafaX/swafa-backend/models"
	"github.com/SwafaX/swafa-backend/utils"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"gorm.io/gorm"
)

type ChatController struct {
	DB *gorm.DB
}

func NewChatController(DB *gorm.DB) ChatController {
	return ChatController{DB}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // Allow all origins for now, adjust for production
	},
}

// Map to store active connections
var clients = make(map[uuid.UUID]*websocket.Conn)

// GetChatsByUser retrieves all chats for a specific user
func (cc *ChatController) GetMyChats(ctx *gin.Context) {
	currentUser := ctx.MustGet("currentUser").(models.User)

	var chats []models.Chat
	result := cc.DB.Where("participant1 = ? OR participant2 = ?", currentUser.ID, currentUser.ID).Find(&chats)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch chats"})
		return
	}

	ctx.JSON(http.StatusOK, chats)
}

// GetChatMessages retrieves messages for a specific chat
func (cc *ChatController) GetChatMessages(ctx *gin.Context) {
	chatIDStr := ctx.Param("chatId")
	chatID, err := uuid.Parse(chatIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid chat ID"})
		return
	}

	// Fetch messages for the chat
	var messages []models.Message
	result := cc.DB.Where("chat_id = ?", chatID).Order("created_at asc").Find(&messages)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch messages"})
		return
	}

	ctx.JSON(http.StatusOK, messages)
}

// CreateChat creates a new chat between two users
func (cc *ChatController) CreateChat(ctx *gin.Context) {
	var req struct {
		Participant1 string `json:"participant1_id" binding:"required"`
		Participant2 string `json:"participant2_id" binding:"required"`
	}

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	participant1_id, err := uuid.Parse(req.Participant1)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant1 ID"})
		return
	}

	participant2_id, err := uuid.Parse(req.Participant2)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid participant2 ID"})
		return
	}

	// Check if chat already exists between these users
	var existingChat models.Chat
	result := cc.DB.Where(
		"(participant1 = ? AND participant2 = ?) OR (participant1 = ? AND participant2 = ?)",
		participant1_id, participant2_id, participant2_id, participant1_id,
	).First(&existingChat)

	if result.Error == nil {
		// Chat already exists
		ctx.JSON(http.StatusOK, existingChat)
		return
	}

	// Create new chat
	newChat := models.Chat{
		ID:           uuid.New(),	
		Participant1: participant1_id,
		Participant2: participant2_id,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	result = cc.DB.Create(&newChat)
	if result.Error != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create chat"})
		return
	}

	ctx.JSON(http.StatusCreated, newChat)
}

// HandleWebSocket handles the WebSocket connection for real-time chat
func (cc *ChatController) HandleWebSocket(ctx *gin.Context) {
	userIDStr := ctx.Query("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	conn, err := upgrader.Upgrade(ctx.Writer, ctx.Request, nil)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upgrade connection"})
		return
	}

	// Register client
	clients[userID] = conn

	// Handle incoming messages
	go func() {
		defer func() {
			conn.Close()
			delete(clients, userID)
		}()

		for {
			// Read message
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				return
			}

			// Parse message
			var msg struct {
				ChatID  string `json:"chat_id"`
				Content string `json:"content"`
				To      string `json:"to"`
			}

			if err := utils.ParseJSON(string(p), &msg); err != nil {
				continue
			}

			// Parse the recipient UUID
			toID, err := uuid.Parse(msg.To)
			if err != nil {
				continue
			}

			// Parse the ChatID as UUID
			chatID, err := uuid.Parse(msg.ChatID)
			if err != nil {
				continue
			}

			// Store message in database
			newMessage := models.Message{
				ID:        uuid.New(),
				ChatID:    chatID,
				SenderID:  userID,
				Content:   msg.Content,
				CreatedAt: time.Now(),
			}

			if result := cc.DB.Create(&newMessage); result.Error != nil {
				continue
			}

			// Send message to recipient if online
			if recipient, ok := clients[toID]; ok {
				messageJSON, _ := utils.ToJSON(newMessage)
				recipient.WriteMessage(messageType, []byte(messageJSON))
			}
		}
	}()
}
