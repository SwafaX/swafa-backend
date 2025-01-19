package routes

import (
	"github.com/calvinnle/todo-app/controllers"
	"github.com/calvinnle/todo-app/middleware"
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

	router.GET("profile", middleware.DeserializeUSer(), uc.userController.GetProfile)
	router.PUT("profile", middleware.DeserializeUSer(), uc.userController.UpdateProfile)
}
