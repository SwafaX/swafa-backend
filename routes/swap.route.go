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

	router.PUT(":swap_id/accept", middleware.DeserializeUser(), sc.swapController.AcceptSwap)
	router.PUT(":swap_id/reject", middleware.DeserializeUser(), sc.swapController.RejectSwap)

	router.GET("", middleware.DeserializeUser(), sc.swapController.GetAllSwaps)
	router.GET("sent", middleware.DeserializeUser(), sc.swapController.GetSentSwaps)
	router.GET("received", middleware.DeserializeUser(), sc.swapController.GetReceivedSwaps)
	router.GET(":swap_id", middleware.DeserializeUser(), sc.swapController.GetSwapByID)

	router.DELETE(":swap_id", middleware.DeserializeUser(), sc.swapController.DeleteSwap)
}
