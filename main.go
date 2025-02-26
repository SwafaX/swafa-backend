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

	SwapController      controllers.SwapController
	SwapRouteController routes.SwapRouteController

	WebSocketController      controllers.WebSocketController
	WebSocketRouteController routes.WebSocketRouteController
)

func init() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	// connect to dependencies
	initializers.ConnectRedis(&config)
	initializers.ConnectDB(&config)
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

	// swap
	SwapController = controllers.NewSwapController(initializers.DB)
	SwapRouteController = routes.NewSwapRouteController(SwapController)

	// websocket
	WebSocketController = controllers.NewWebSocketController(initializers.DB)
	WebSocketRouteController = routes.NewWebSocketRouteController(WebSocketController)

	// server
	server = gin.Default()
}

func main() {
	config, err := initializers.LoadConfig(".")

	if err != nil {
		log.Fatal("Could not load environment variables", err)
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

	// core
	AuthRouteController.AuthRoute(router)
	ItemRouteController.ItemRoute(router)
	UserRouteController.UserRoute(router)
	SwapRouteController.SwapRoute(router)

	// photo uploading
	PresignedURLRouteController.PresignedURLRoute(router)

	// websocket
	WebSocketRouteController.WebSocketRoute(router)

	log.Fatal(server.Run("0.0.0.0:" + config.ServerPort))
}
