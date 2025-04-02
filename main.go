package main

import (
	"log"
	"net/http"

	"github.com/SwafaX/swafa-backend/controllers"
	"github.com/SwafaX/swafa-backend/initializers"
	"github.com/SwafaX/swafa-backend/routes"
	socketio "github.com/doquangtan/socket.io/v4"
	"github.com/gin-contrib/static"
	"github.com/gin-gonic/gin"

	docs "github.com/SwafaX/swafa-backend/docs"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	server *gin.Engine
	io     *socketio.Io

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

	ChatController      controllers.ChatController
	ChatRouteController routes.ChatRouteController
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

	// Initialize Gin server first
	server = gin.Default()

	// Initialize Socket.IO
	io = socketio.New()

	// Initialize all controllers
	AuthController = controllers.NewAuthController(initializers.DB, initializers.RedisClient)
	AuthRouteController = routes.NewAuthRouteController(AuthController)

	ItemController = controllers.NewItemController(initializers.DB)
	ItemRouteController = routes.NewItemRouteController(ItemController)

	UserController = controllers.NewUserController(initializers.DB, initializers.MinioClient)
	UserRouteController = routes.NewUserRouteController(UserController)

	PresignedURLController = controllers.NewPresignedURLController(initializers.MinioClient)
	PresignedURLRouteController = routes.NewPresignedURLRouteController(PresignedURLController)

	SwapController = controllers.NewSwapController(initializers.DB)
	SwapRouteController = routes.NewSwapRouteController(SwapController)

	ChatController = controllers.NewChatController(initializers.DB, io)
	ChatRouteController = routes.NewChatRouteController(ChatController)

	// Setup Socket.IO handlers
	ChatController.SetupSocketHandlers()

	// Setup static files and Socket.IO handler
	server.Use(static.Serve("/", static.LocalFile("./public", false)))
	server.GET("/socket.io/*any", gin.WrapH(io.HttpHandler()))
}

func main() {
	config, err := initializers.LoadConfig(".")
	if err != nil {
		log.Fatal("Could not load environment variables", err)
	}

	docs.SwaggerInfo.BasePath = "/api/v1"
	router := server.Group("/api/v1")
	router.GET("/healthcheck", func(c *gin.Context) {
		message := "Welcome to my app"
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
	ChatRouteController.ChatRoute(router)

	log.Fatal(server.Run("0.0.0.0:" + config.ServerPort))
}
