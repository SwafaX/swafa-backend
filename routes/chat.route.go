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

	router.GET("/messages/:chatId", crc.chatController.GetChatMessages)
}
