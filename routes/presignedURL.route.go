package routes

import (
	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/middleware"
	"github.com/gin-gonic/gin"
)

type PresignedURLRouteController struct {
	presignedURLController controllers.PresignedURLController
}

func NewPresignedURLRouteController(presignedURLController controllers.PresignedURLController) PresignedURLRouteController {
	return PresignedURLRouteController{presignedURLController}
}

func (ri *PresignedURLRouteController) PresignedURLRoute(rg *gin.RouterGroup) {
	router := rg.Group("presigned-url")

	router.GET("", middleware.DeserializeUSer(), ri.presignedURLController.PresignedURLGenerator)
}
