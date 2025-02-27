package routes

import (
	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/middleware"
	"github.com/gin-gonic/gin"
)

type WebSocketRouteController struct {
	websocketController controllers.WebSocketController
}

func NewWebSocketRouteController(websocketController controllers.WebSocketController) WebSocketRouteController {
	return WebSocketRouteController{
		websocketController: websocketController,
	}
}

func (wc *WebSocketRouteController) WebSocketRoute(rg *gin.RouterGroup) {
	router := rg.Group("/ws")
	router.GET("/chat", middleware.DeserializeUSer(), wc.websocketController.HandleWebSocket)
} 