package routes

import (
	"github.com/calvinnle/todo-app/controllers"
	"github.com/calvinnle/todo-app/middleware"
	"github.com/gin-gonic/gin"
)

type ImageRouteController struct {
	imageController controllers.ImageController
}

func NewImageRouteController(imageController controllers.ImageController) ImageRouteController {
	return ImageRouteController{imageController}
}

func (ri *ImageRouteController) ImageRoute(rg *gin.RouterGroup) {
	router := rg.Group("presigned-url")

	router.GET("", middleware.DeserializeUSer(), ri.imageController.PresignedURLGenerator)
}
