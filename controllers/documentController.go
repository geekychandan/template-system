package controllers

import (
	"log"
	"net/http"
	"os"
	"path/filepath"
	"template-system/services"

	"github.com/gin-gonic/gin"
)

func GenerateDocument(c *gin.Context) {
	templateID := c.Param("id")
	var input map[string]string
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	documentPath, err := services.GenerateDocument(templateID, input, c.MustGet("userId").(string))
	if err != nil {
		log.Printf("Error generating document: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Document generated successfully", "document_path": documentPath})
}

func GetGeneratedDocuments(c *gin.Context) {
	userID := c.MustGet("userId").(string)
	documents, err := services.GetGeneratedDocumentsByUserID(userID)
	if err != nil {
		log.Printf("Error getting generated documents: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"documents": documents})
}

func DownloadDocument(c *gin.Context) {
	documentID := c.Param("id")
	documentPath, err := services.GetDocumentPathByID(documentID)
	if err != nil {
		log.Printf("Error getting document path: %v", err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}
	log.Printf("Document path: %s", documentPath)

	// Ensure the tmp directory exists
	tmpDir := "tmp"
	if _, err := os.Stat(tmpDir); os.IsNotExist(err) {
		err := os.Mkdir(tmpDir, 0755)
		if err != nil {
			log.Printf("Error creating tmp directory: %v", err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	// Download from S3 to a temporary local file
	localPath := filepath.Join(tmpDir, filepath.Base(documentPath))
	err = services.DownloadFileFromS3(documentPath, localPath)
	if err != nil {
		log.Printf("Error downloading file from S3: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to download document"})
		return
	}

	// Serve the file
	c.File(localPath)
}
