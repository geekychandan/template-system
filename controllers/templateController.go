package controllers

import (
	"net/http"
	"path/filepath"
	"template-system/models"
	"template-system/services"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func UploadTemplate(c *gin.Context) {
	userID := c.MustGet("userId").(string)
	file, _ := c.FormFile("template")
	fileName := filepath.Base(file.Filename)
	fileID := uuid.New().String()
	filePath := "templates/" + fileID + "_" + fileName

	if err := c.SaveUploadedFile(file, filePath); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to upload file"})
		return
	}

	template := models.Template{
		UserID:       userID,
		TemplateName: fileName,
		FilePath:     filePath,
	}

	if err := services.CreateTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Template uploaded successfully", "template": template})
}

func GetTemplates(c *gin.Context) {
	userID := c.MustGet("userId").(string)
	templates, err := services.GetTemplatesByUserID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"templates": templates})
}

func GetPlaceholders(c *gin.Context) {
	templateID := c.Param("id")
	placeholders, err := services.ExtractPlaceholders(templateID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"placeholders": placeholders})
}
