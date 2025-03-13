package routes

import (
	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/middleware"
	"github.com/gin-gonic/gin"
)

type AuthRouteController struct {
	authController controllers.AuthController
}

func NewAuthRouteController(authController controllers.AuthController) AuthRouteController {
	return AuthRouteController{authController}
}

func (rc *AuthRouteController) AuthRoute(rg *gin.RouterGroup) {
	router := rg.Group("auth")

	router.POST("register", rc.authController.Register)
	router.POST("login", rc.authController.LogIn)
	router.GET("refresh", rc.authController.Refresh)
	router.POST("logout", middleware.DeserializeUser(), rc.authController.LogOut)
}
