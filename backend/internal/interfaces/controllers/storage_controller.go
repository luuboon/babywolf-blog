package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	
	"babywolf-blog-backend/internal/application/usecases"
)

type StorageController struct {
	storageUseCase usecases.StorageUseCase
}

func NewStorageController(suc usecases.StorageUseCase) *StorageController {
	return &StorageController{storageUseCase: suc}
}

func (c *StorageController) UploadFile(ctx *gin.Context) {
	fileHeader, err := ctx.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}

	folder := ctx.PostForm("folder")
	if folder == "" {
		folder = "uploads"
	}

	file, err := fileHeader.Open()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process file"})
		return
	}
	defer file.Close()

	// Append timestamp to prevent overwriting
	fileName := fmt.Sprintf("%d-%s", time.Now().Unix(), fileHeader.Filename)

	url, err := c.storageUseCase.UploadFile(folder, fileName, file)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed", "details": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     url,
	})
}
