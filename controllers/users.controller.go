package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SwafaX/swafa-backend/initializers"
	"github.com/SwafaX/swafa-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
	"gorm.io/gorm"
)

type UserController struct {
	DB    *gorm.DB
	MinIO *minio.Client
}

func NewUserController(DB *gorm.DB, MinIO *minio.Client) UserController {
	return UserController{
		DB:    DB,
		MinIO: MinIO,
	}
}

func (uc *UserController) GetMyProfile(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var dbUser models.User
	result := initializers.DB.First(&dbUser, "id = ?", currentUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "User no longer exists",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":    dbUser.ID,
			"name":  dbUser.Name,
			"email": dbUser.Email,
			"age":   dbUser.Age,
		},
	})
}

func (uc *UserController) UpdateMyProfile(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var dbUser models.User
	result := uc.DB.First(&dbUser, "id = ?", currentUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusForbidden, gin.H{
			"status":  "fail",
			"message": "User not found",
		})
		return
	}

	var updateInfo models.UpdateUserInput
	if err := c.ShouldBindJSON(&updateInfo); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"message": fmt.Sprintf("Invalid input: %v", err),
		})
		return
	}

	// Update
	if updateInfo.Name != "" {
		dbUser.Name = updateInfo.Name
	}
	if updateInfo.Age > 0 {
		dbUser.Age = int64(updateInfo.Age)
	}

	now := time.Now()

	if err := uc.DB.Save(&dbUser).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to update profile",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":         dbUser.ID,
			"name":       dbUser.Name,
			"email":      dbUser.Email,
			"age":        dbUser.Age,
			"created_at": dbUser.CreatedAt,
			"updated_at": now,
		},
	})
}

func (uc *UserController) GetUserProfile(c *gin.Context) {
	userID := c.Param("user_id")

	var user models.User

	// Fetch the user profile by ID
	if err := uc.DB.First(&user, "id = ?", userID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  "User not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data": gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
			"age":   user.Age,
		},
	})
}

// Get all items belongs to a user
func (uc *UserController) ShowItems(c *gin.Context) {
	userID := c.Param("user_id")

	var items *[]models.Item
	if err := uc.DB.Where("user_id = ?", userID).Find(&items).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  "No items found for this user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Items fetched successfully",
		"data":    items,
	})
}
