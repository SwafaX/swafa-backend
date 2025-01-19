package routes

import (
	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/middleware"
	"github.com/gin-gonic/gin"
)

type SwapRouteController struct {
	swapController controllers.SwapController
}

func NewSwapRouteController(swapController controllers.SwapController) SwapRouteController {
	return SwapRouteController{swapController}
}

func (sc *SwapRouteController) SwapRoute(rg *gin.RouterGroup) {
	router := rg.Group("swap")

	router.PUT(":swap_id/accept", middleware.DeserializeUSer(), sc.swapController.AcceptSwap)
	router.PUT(":swap_id/reject", middleware.DeserializeUSer(), sc.swapController.RejectSwap)

	router.GET("", middleware.DeserializeUSer(), sc.swapController.GetAllSwaps)
	router.GET("sent", middleware.DeserializeUSer(), sc.swapController.GetSentSwaps)
	router.GET("received", middleware.DeserializeUSer(), sc.swapController.GetReceivedSwaps)

	router.DELETE(":swap_id", middleware.DeserializeUSer(), sc.swapController.DeleteSwap)
}
