package controllers

import (
	"log"
	"net/http"

	// "os"
	// "path/filepath"
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
	// Get document ID from the request
	documentID := c.Param("id")

	// Fetch the S3 path for the document
	documentPath, err := services.GetDocumentPathByID(documentID)
	if err != nil {
		log.Printf("Error fetching document path for ID %s: %v", documentID, err)
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	// Generate a presigned URL for the document
	urlStr, err := services.GeneratePresignedURL(documentPath)
	if err != nil {
		log.Printf("Error generating pre-signed URL for document ID %s: %v", documentID, err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate download link"})
		return
	}

	// Return the presigned URL in the response
	// c.Redirect(http.StatusFound, urlStr)
	c.JSON(http.StatusOK, gin.H{"download_url": urlStr})

}
