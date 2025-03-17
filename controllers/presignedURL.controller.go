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
	fileType := c.Query("file_type")

	currentUser := c.MustGet("currentUser").(models.User)

	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "filename is required"})
		return
	}

	if fileType == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_type is required"})
		return
	}

	var path string

	found, err := ic.Minio.BucketExists(context.Background(), bucketName)
	if err != nil {
		fmt.Println(err)
		return
	}
	if found {
		fmt.Println("Bucket found")
	} else {
		err := ic.Minio.MakeBucket(context.Background(), bucketName, minio.MakeBucketOptions{Region: "us-east-1", ObjectLocking: true})
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	fmt.Println("Successfully created mybucket.")
	// Define file type for each case
	if fileType == "avatar" {
		path = fmt.Sprintf("images/%s/avatar/avatar.jpg", currentUser.ID)
	} else if fileType == "item" {
		path = fmt.Sprintf("images/%s/items/item-%d.jpg", currentUser.ID, time.Now().UnixNano())
	} else {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file_type is invalid"})
		return
	}

	// Generate presigned URL
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
