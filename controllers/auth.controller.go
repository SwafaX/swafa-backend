package controllers

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/calvinnle/todo-app/initializers"
	"github.com/calvinnle/todo-app/models"
	"github.com/calvinnle/todo-app/utils"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type AuthController struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewAuthController(DB *gorm.DB, Redis *redis.Client) AuthController {
	return AuthController{
		DB:    DB,
		Redis: Redis,
	}
}

// Register
//
//	@Summary		Register a new user
//	@Description	Create a new user account with name, email, and password.
//	@Tags			Auth
//	@Accept			json
//	@Produce		json
//	@Param			registerInfo	body		models.RegisterInput	true	"User registration data"
//	@Router			/auth/register [post]
func (ac *AuthController) Register(c *gin.Context) {
	var registerInfo *models.RegisterInput

	if err := c.ShouldBindJSON(&registerInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": fmt.Sprintf("Validation failed: %v", err),
		})
		return
	}

	if err := utils.ValidateEmail(registerInfo.Email); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid email",
		})
		return
	}

	if registerInfo.Password != registerInfo.PasswordConfirm {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Passwords do not match",
		})
		return
	}

	hashedPassword, err := utils.HashPassword(registerInfo.Password)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": err.Error(),
		})
		return
	}

	now := time.Now()

	newUser := &models.User{
		ID:        uuid.New(),
		Name:      registerInfo.Name,
		Email:     registerInfo.Email,
		Age:       registerInfo.Age,
		Password:  hashedPassword,
		CreatedAt: now,
		UpdatedAt: now,
	}

	result := ac.DB.Create(&newUser)

	if result.Error != nil {
		if strings.Contains(result.Error.Error(), "duplicate key") {
			c.JSON(http.StatusConflict, gin.H{
				"status":  "fail",
				"message": "User already exists",
			})
			return
		}
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "error",
			"message": result.Error.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status": "success",
		"data":   newUser,
	})
}

func (ac *AuthController) LogIn(c *gin.Context) {
	var SignInInput *models.SignInInput

	if err := c.ShouldBindJSON(&SignInInput); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Check db if user existed
	var user models.User

	result := ac.DB.First(&user, "email = ?", strings.ToLower(SignInInput.Email))
	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid email or password",
		})
		return
	}

	if err := utils.VerifyPassword(user.Password, SignInInput.Password); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Invalid email or password",
		})
		return
	}
	// Get token config
	config, _ := initializers.LoadConfig(".")

	// Generate tokens and return to users
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.AccessTokenPrivate)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":        "logged in",
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}

func (ac *AuthController) Refresh(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Refresh token is required",
		})
		return
	}

	refreshToken := requestBody.RefreshToken

	// Load configuration
	config, _ := initializers.LoadConfig(".")

	// Validate the refresh token
	sub, err := utils.ValidateToken(refreshToken, config.RefreshTokenPrivate)
	if err != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": err.Error(),
		})
		return
	}

	// Check if token is blacklisted
	isBlacklisted, err := utils.IsTokenBlacklisted(refreshToken, ac.Redis)
	if isBlacklisted {
		c.JSON(http.StatusUnauthorized, gin.H{"status": "fail", "message": "Invalid token, please log in again"})
		return
	}

	// Lookup the user in the database
	var user models.User
	if result := ac.DB.First(&user, "id = ?", sub); result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "User no longer exists",
		})
		return
	}

	// Blacklist the old refresh token
	utils.BlacklistToken(refreshToken, ac.Redis, config.RefreshTokenExpiresIn)

	// Generate new tokens
	access_token, err := utils.CreateToken(config.AccessTokenExpiresIn, user.ID, config.AccessTokenPrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Could not create access token"})
		return
	}

	refresh_token, err := utils.CreateToken(config.RefreshTokenExpiresIn, user.ID, config.RefreshTokenPrivate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"status": "fail", "message": "Could not create refresh token"})
		return
	}

	// Send new tokens in response
	c.JSON(http.StatusOK, gin.H{
		"status":        "success",
		"access_token":  access_token,
		"refresh_token": refresh_token,
	})
}

func (ac *AuthController) LogOut(c *gin.Context) {
	var requestBody struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": "Refresh token is required for logout",
		})
		return
	}

	user, exists := c.Get("currentUser")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{
			"status":  "fail",
			"message": "User not authenticated",
		})
		return
	}

	// Check if user existed
	currentUser := user.(models.User)
	var dbUser models.User
	result := initializers.DB.First(&dbUser, "id = ?", currentUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "User no longer exists",
		})
		return
	}

	// Blacklist
	err := utils.BlacklistToken(requestBody.RefreshToken, initializers.RedisClient, time.Hour*24)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to blacklist refresh token",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Logged out successfully",
	})
}
