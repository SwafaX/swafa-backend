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
		ImageUrl:    payload.ImageUrl,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	result := ic.DB.Create(&newItem)
	if result.Error != nil {
		c.JSON(http.StatusBadGateway, gin.H{
			"status":  "fail",
			"details": "Couldn't create item.",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "sucess",
		"data":    newItem,
	})
}

// Get all items for newsfeed
func (ic *ItemController) GetAllItems(c *gin.Context) {
	var allItems *[]models.Item

	result := ic.DB.Find(&allItems)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Unable to fetch items.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   allItems,
	})
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

func (ic *ItemController) GetItemByID(c *gin.Context) {
	item_id := c.Param("item_id")

	var item *models.Item

	result := ic.DB.First(&item, "id = ?", item_id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"details": "Item not found",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"item":   item,
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

func (ic *ItemController) UpdateItem(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	itemID := c.Param("item_id")

	var item models.Item
	result := ic.DB.First(&item, "id = ? AND user_id = ?", itemID, currentUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"details": "Couldn't find item",
		})
		return
	}

	var payload struct {
		Title       string `json:"title"`
		Description string `json:"description"`
	}
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"details": "Invalid input",
		})
		return
	}

	item.Title = payload.Title
	item.Description = payload.Description
	item.UpdatedAt = time.Now()

	// Save changes
	if err := ic.DB.Save(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to update item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item updated successfully",
		"data":    item,
	})
}

func (ic *ItemController) DeleteItem(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	itemID := c.Param("item_id")

	// Fetch the existing item
	var item models.Item
	result := ic.DB.First(&item, "id = ? AND user_id = ?", itemID, currentUser.ID)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"details": "Couldn't find item",
		})
	}

	if err := ic.DB.Delete(&item).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Failed to delete item",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Item deleted successfully",
	})
}
