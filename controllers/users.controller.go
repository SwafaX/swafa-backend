package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/calvinnle/todo-app/initializers"
	"github.com/calvinnle/todo-app/models"
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

func (uc *UserController) GetProfile(c *gin.Context) {
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

func (uc *UserController) UpdateProfile(c *gin.Context) {
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
