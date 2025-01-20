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

	router.GET("", middleware.DeserializeUSer(), rc.ItemController.GetAllItems)
	router.POST("", middleware.DeserializeUSer(), rc.ItemController.CreateItems)
	router.POST(":item_id/swap_request", middleware.DeserializeUSer(), rc.ItemController.CreateSwapRequest)

	router.GET("me", middleware.DeserializeUSer(), rc.ItemController.GetMyItems)
	router.GET(":item_id", middleware.DeserializeUSer(), rc.ItemController.GetItemByID)
	router.PUT(":item_id", middleware.DeserializeUSer(), rc.ItemController.UpdateItem)
	router.DELETE(":item_id", middleware.DeserializeUSer(), rc.ItemController.DeleteItem)
}
