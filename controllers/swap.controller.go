package controllers

import (
	"net/http"
	"time"

	"github.com/SwafaX/swafa-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SwapController struct {
	DB *gorm.DB
}

func NewSwapController(DB *gorm.DB) SwapController {
	return SwapController{DB}
}

// Move to ItemController

// func (sc *SwapController) CreateSwap(c *gin.Context) {
// 	currentUser := c.MustGet("currentUser").(models.User)
//
// 	var payload struct {
// 		RequestedItemID string `json:"request_item_id" binding:"required"`
// 		Message         string `json:"message"`
// 	}
//
// 	if err := c.ShouldBindJSON(&payload); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status":  "fail",
// 			"error":   "Invalid input",
// 			"details": err.Error(),
// 		})
// 		return
// 	}
//
// 	// Check the requested item exists and is valid.
// 	var requestItem models.Item
// 	if err := sc.DB.First(&requestItem, "id = ?", payload.RequestedItemID).Error; err != nil {
// 		c.JSON(http.StatusNotFound, gin.H{
// 			"status": "fail",
// 			"error":  "Requested item not found",
// 		})
// 		return
// 	}
//
// 	// Prevent the user from requesting their own item.
// 	if requestItem.UserID == currentUser.ID {
// 		c.JSON(http.StatusBadRequest, gin.H{
// 			"status": "fail",
// 			"error":  "You cannot swap for your own item",
// 		})
// 		return
// 	}
//
// 	recipientID := requestItem.UserID
//
// 	// Create the new swap.
// 	swap := &models.Swap{
// 		ID:            uuid.New(),
// 		RequesterID:   currentUser.ID,
// 		RecipientID:   recipientID,
// 		RequestItemID: requestItem.ID,
// 		Message:       payload.Message,
// 		Status:        "pending",
// 		CreatedAt:     time.Now(),
// 		UpdatedAt:     time.Now(),
// 	}
//
// 	// Create Swap Request
// 	if err := sc.DB.Create(swap).Error; err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{
// 			"status":  "fail",
// 			"details": "Failed to create swap",
// 		})
// 		return
// 	}
//
// 	c.JSON(http.StatusCreated, gin.H{
// 		"status":  "Created",
// 		"message": "Swap created successfully",
// 		"data":    swap,
// 	})
// }

func (sc *SwapController) AcceptSwap(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	swapID := c.Param("swap_id")

	var swap models.Swap

	if err := sc.DB.First(&swap, "id = ?", swapID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status":  "fail",
			"details": "Swap not found",
		})
		return
	}

	// verify authorization
	if swap.RecipientID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to accept this swap",
		})
		return
	}

	// start a transaction
	tx := sc.DB.Begin()

	// update swap status
	swap.Status = "accepted"
	swap.UpdatedAt = time.Now()

	if err := tx.Save(&swap).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to accept swap",
		})
		return
	}

	// create a new chat between requester and recipient
	newChat := models.Chat{
		ID:           uuid.New(),
		Participant1: swap.RequesterID,
		Participant2: swap.RecipientID,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := tx.Create(&newChat).Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create chat",
		})
		return
	}

	// commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to complete the swap acceptance",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "success",
		"message": "Swap accepted and chat created",
		"swap":    swap,
		"chat":    newChat,
	})
}

func (sc *SwapController) RejectSwap(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	swapID := c.Param("swap_id")

	var payload struct {
		Message string `json:"message"`
	}

	var swap *models.Swap

	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"details": "Invalid input",
		})
		return
	}

	if err := sc.DB.First(&swap, "id = ?", swapID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  "Swap not found",
		})
		return
	}

	// in case the swap displayed on the wrong user (might come from Fat's mistake)
	if swap.RecipientID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{
			"error": "You are not authorized to reject this swap",
		})
		return
	}

	// Update the status to 'rejected' and store rejection message
	swap.Status = "rejected"
	swap.UpdatedAt = time.Now()

	if err := sc.DB.Save(&swap).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reject swap"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "rejected",
		"message": payload.Message,
		"swap":    swap,
	})
}

func (sc *SwapController) GetAllSwaps(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var allSwaps []models.Swap

	// Fetch all swaps of a user (either sent or received)
	query := sc.DB.Where("requester_id = ? OR recipient_id = ?", currentUser.ID, currentUser.ID)
	if err := query.Find(&allSwaps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch swaps"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Swaps fetched successfully",
		"data":    allSwaps,
	})
}

func (sc *SwapController) GetSwapByID(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	swapID := c.Param("swap_id")

	var swap models.Swap

	query := sc.DB.Where("id = ? AND (requester_id = ? OR recipient_id = ?)", swapID, currentUser.ID, currentUser.ID)
	if err := query.First(&swap).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"status": "fail",
			"error":  "Swap not found or you don't have access to it",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Swap fetched successfully",
		"data":    swap,
	})
}

func (sc *SwapController) GetSentSwaps(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var sentSwaps []models.Swap

	if err := sc.DB.Where("requester_id = ?", currentUser.ID).Find(&sentSwaps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch sent swaps"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Sent swaps fetched successfully",
		"data":    sentSwaps,
	})
}

func (sc *SwapController) GetReceivedSwaps(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)

	var receivedSwaps []models.Swap

	if err := sc.DB.Where("recipient_id = ?", currentUser.ID).Find(&receivedSwaps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch received swaps"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Received swaps fetched successfully",
		"data":    receivedSwaps,
	})
}

func (sc *SwapController) DeleteSwap(c *gin.Context) {
	currentUser := c.MustGet("currentUser").(models.User)
	swapID := c.Param("swap_id")

	var swap models.Swap

	if err := sc.DB.First(&swap, "id = ?", swapID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Swap not found"})
		return
	}

	if swap.RequesterID != currentUser.ID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this swap"})
		return
	}

	if swap.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "You can only delete pending swaps"})
		return
	}

	if err := sc.DB.Delete(&swap).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete swap"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Swap deleted successfully",
	})
}
