package main

import (
	"log"
	"net/http"

	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/initializers"
	"github.com/SwafaX/swafa-backend/routes"
	"github.com/gin-gonic/gin"

	docs "github.com/SwafaX/swafa-backend/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	server *gin.Engine

	AuthController      controllers.AuthController
	AuthRouteController routes.AuthRouteController

	ItemController      controllers.ItemController
	ItemRouteController routes.ItemRouteController

	UserController      controllers.UserController
	UserRouteController routes.UserRouteController

	PresignedURLController      controllers.PresignedURLController
	PresignedURLRouteController routes.PresignedURLRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// connect to dependencies
	initializers.ConnectDB(&config)
	initializers.ConnectRedis(&config)
	initializers.ConnectMinio(&config)

	// auth
	AuthController = controllers.NewAuthController(initializers.DB, initializers.RedisClient)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	// item
	ItemController = controllers.NewItemController(initializers.DB)
	ItemRouteController = routes.NewItemRouteController(ItemController)

	// user
	UserController = controllers.NewUserController(initializers.DB, initializers.MinioClient)
	UserRouteController = routes.NewUserRouteController(UserController)

	// presigned URL
	PresignedURLController = controllers.NewPresignedURLController(initializers.MinioClient)
	PresignedURLRouteController = routes.NewPresignedURLRouteController(PresignedURLController)

	// server
	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("🚀 Could not load environment variables", err)
	}
	docs.SwaggerInfo.BasePath = "/api/v1"
	router := server.Group("/api/v1")
	router.GET("/healthcheck", func(c *gin.Context) {
		message := "Welcome to my todo app"
		c.JSON(http.StatusOK, gin.H{
			"message": message,
		})
	})
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	AuthRouteController.AuthRoute(router)
	ItemRouteController.ItemRoute(router)
	UserRouteController.UserRoute(router)
	PresignedURLRouteController.PresignedURLRoute(router)

	log.Fatal(server.Run("localhost:" + config.ServerPort))
}
