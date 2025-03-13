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

	router.GET("", middleware.DeserializeUser(), rc.ItemController.GetAllItems)
	router.POST("", middleware.DeserializeUser(), rc.ItemController.CreateItems)
	router.POST(":item_id/swap_request", middleware.DeserializeUser(), rc.ItemController.CreateSwapRequest)

	router.GET("me", middleware.DeserializeUser(), rc.ItemController.GetMyItems)
	router.GET(":item_id", middleware.DeserializeUser(), rc.ItemController.GetItemByID)
	router.PUT(":item_id", middleware.DeserializeUser(), rc.ItemController.UpdateItem)
	router.DELETE(":item_id", middleware.DeserializeUser(), rc.ItemController.DeleteItem)
}
