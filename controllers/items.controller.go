package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/SwafaX/swafa-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ItemController struct {
	DB *gorm.DB
}

func NewItemController(DB *gorm.DB) ItemController {
	return ItemController{
		DB: DB,
	}
}

func (ic *ItemController) CreateItems(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var payload *models.ItemCreation

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid input",
			"details": fmt.Sprintf("Validation failed: %v", err),
		})
		return
	}

	now := time.Now()

	newItem := &models.Item{
		ID:          uuid.New(),
		Title:       payload.Title,
		Description: payload.Description,
		UserID:      currentUser.ID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := ic.DB.Create(&newItem)

	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{"status": "error", "message": result.Error.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "sucess",
		"data":    newItem,
	})
}

// Get all items for newsfeed
func (ic *ItemController) GetAllItems(c *gin.Context) {
	// To be implemented
}

func (ic *ItemController) GetMyItems(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var items *[]models.Item

	result := ic.DB.Where("user_id = ?", currentUser.ID).Find(&items)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Unable to get my items.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   items,
	})
}

func (ic *ItemController) CreateSwapRequest(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	itemID := c.Param("item_id")

	var payload struct {
		Message string `json:"message"`
	}

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"error":   "Invalid input",
			"details": err.Error(),
		})
		return
	}

	// Check if item exists
	var item models.Item
	if err := ic.DB.First(&item, "id = ?", itemID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  "Item not found",
		})
		return
	}

	// Prevent the user from requesting their own item
	if item.UserID == currentUser.ID {
		c.JSON(http.StatusBadRequest, gin.H{
			"status": "fail",
			"error":  "You cannot swap for your own item",
		})
		return
	}

	swap := &models.Swap{
		RequesterID:   currentUser.ID,
		RecipientID:   item.UserID,
		RequestItemID: item.ID,
		Message:       payload.Message,
		Status:        "pending",
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if err := ic.DB.Create(&swap).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"details": "Failed to create swap",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"status":  "Created",
		"message": "Swap created successfully",
		"data":    swap,
	})
}
