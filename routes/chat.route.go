package routes

import (
	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/middleware"
	"github.com/gin-gonic/gin"
)

type ChatRouteController struct {
	chatController controllers.ChatController
}

func NewChatRouteController(chatController controllers.ChatController) ChatRouteController {
	return ChatRouteController{chatController}
}

func (crc *ChatRouteController) ChatRoute(rg *gin.RouterGroup) {
	router := rg.Group("chats")
	router.Use(middleware.DeserializeUser())

	router.GET("me", crc.chatController.GetMyChats)
	router.GET(":chatId/messages", crc.chatController.GetChatMessages)
	router.POST("", crc.chatController.CreateChat)
	router.GET("ws", crc.chatController.HandleWebSocket)
} 