package controllers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/SwafaX/swafa-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/minio/minio-go/v7"
)

type PresignedURLController struct {
	Minio *minio.Client
}

func NewPresignedURLController(Minio *minio.Client) PresignedURLController {
	return PresignedURLController{Minio}
}

const (
	bucketName = "images"
)

// Gen presigned url
func (ic *PresignedURLController) PresignedURLGenerator(c *gin.Context) {
	filename := c.Query("filename")

	currentUser := c.MustGet("currentUser").(models.User)

	var requestBody struct {
		File_type string `json:"file_type" binding:"required"`
	}

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"status":  "fail",
			"details": "Cannot get file-type",
		})
		return
	}

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	var path string

	// define file type for each case
	if requestBody.File_type == "avatar" {
		path = fmt.Sprintf("users/%s/avatar/avatar.jpg", currentUser.ID)
	} else if requestBody.File_type == "item" {
		path = fmt.Sprintf("users/%s/items/item-%d.jpg", currentUser.ID, time.Now().UnixNano())
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_type is invalid"})
		return
	}

	presignedURL, err := ic.Minio.PresignedPutObject(context.Background(), bucketName, path, time.Minute*15)
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	log.Printf("Generated presigned URL for upload: %s", presignedURL.String())
	c.JSON(http.StatusOK, gin.H{
		"url": presignedURL.String(),
	})
}
