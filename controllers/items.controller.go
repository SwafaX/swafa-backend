package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/calvinnle/todo-app/models"
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

func (ic *ItemController) Finish(c *gin.Context) {
	item_id := c.Param("item_id")

	var item models.Item

	result := ic.DB.First(&item, "id = ?", item_id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "Item not found",
		})
		return
	}

	now := time.Now()

	item.Status = "done"
	item.UpdatedAt = now

	result = ic.DB.Save(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Unable to update item status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   item,
	})
}

func (ic *ItemController) Unfinish(c *gin.Context) {
	item_id := c.Param("item_id")

	var item *models.Item

	result := ic.DB.First(&item, "id = ?", item_id)
	if result.Error != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"message": "Item not found",
		})
		return
	}

	now := time.Now()

	item.Status = "unfinished"
	item.UpdatedAt = now

	result = ic.DB.Save(&item)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Unable to update item status",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   item,
	})
}

func (ic *ItemController) GetAllItems(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var items *[]models.Item

	result := ic.DB.Where("user_id = ?", currentUser.ID).Find(&items)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "fail",
			"message": "Unable to get items.",
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"status": "success",
		"data":   items,
	})
}
