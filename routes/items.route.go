package routes

import (
	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/middleware"
	"github.com/gin-gonic/gin"
)

type ItemRouteController struct {
	ItemController controllers.ItemController
}

func NewItemRouteController(ItemController controllers.ItemController) ItemRouteController {
	return ItemRouteController{ItemController}
}

func (rc *ItemRouteController) ItemRoute(rg *gin.RouterGroup) {
	router := rg.Group("items")

	router.POST("", middleware.DeserializeUSer(), rc.ItemController.CreateItems)
	router.POST(":item_id/finish", middleware.DeserializeUSer(), rc.ItemController.Finish)
	router.POST(":item_id/unfinish", middleware.DeserializeUSer(), rc.ItemController.Unfinish)

	router.GET("", middleware.DeserializeUSer(), rc.ItemController.GetAllItems)
}
