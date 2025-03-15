package routes

import (
	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/middleware"
	"github.com/gin-gonic/gin"
)

type UserRouteController struct {
	userController controllers.UserController
}

func NewUserRouteController(userController controllers.UserController) UserRouteController {
	return UserRouteController{userController}
}

func (uc *UserRouteController) UserRoute(rg *gin.RouterGroup) {
	router := rg.Group("user")

	router.GET("me/profile", middleware.DeserializeUser(), uc.userController.GetMyProfile)
	router.PUT("me/profile", middleware.DeserializeUser(), uc.userController.UpdateMyProfile)

	router.GET(":user_id/profile", middleware.DeserializeUser(), uc.userController.GetUserProfile)
	router.GET(":user_id/items", middleware.DeserializeUser(), uc.userController.ShowItems)
}
