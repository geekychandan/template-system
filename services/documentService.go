package services

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"template-system/config"
	"template-system/models"
	"template-system/utils"
	"time"

	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"

	// "github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/google/uuid"
)

func GenerateDocument(templateID string, placeholders map[string]string, userID string) (string, error) {
	var tmpl models.Template
	if err := utils.DB.Where("id = ?", templateID).First(&tmpl).Error; err != nil {
		log.Printf("Error finding template: %v", err)
		return "", err
	}

	// Open the DOCX file as a zip archive
	reader, err := zip.OpenReader(tmpl.FilePath)
	if err != nil {
		log.Printf("Error opening DOCX file: %v", err)
		return "", err
	}
	defer reader.Close()

	// Create a buffer to store the modified zip content
	buf := new(bytes.Buffer)
	writer := zip.NewWriter(buf)

	// Loop through each file in the zip archive
	for _, file := range reader.File {
		f, err := file.Open()
		if err != nil {
			log.Printf("Error opening file inside DOCX: %v", err)
			return "", err
		}
		defer f.Close()

		// Read the file content
		content, err := ioutil.ReadAll(f)
		if err != nil {
			log.Printf("Error reading file content: %v", err)
			return "", err
		}

		// If it's the document.xml file, replace placeholders
		if strings.HasSuffix(file.Name, "document.xml") {
			docContent := string(content)
			for key, value := range placeholders {
				docContent = strings.ReplaceAll(docContent, "{{"+key+"}}", value)
			}
			content = []byte(docContent)
		}

		// Create a new file in the zip archive and write the content
		w, err := writer.Create(file.Name)
		if err != nil {
			log.Printf("Error creating file in zip archive: %v", err)
			return "", err
		}
		_, err = w.Write(content)
		if err != nil {
			log.Printf("Error writing content to file in zip archive: %v", err)
			return "", err
		}
	}

	// Close the zip writer to finalize the archive
	writer.Close()

	// Ensure the documents directory exists
	documentDir := "documents"
	if _, err := os.Stat(documentDir); os.IsNotExist(err) {
		err := os.Mkdir(documentDir, 0755)
		if err != nil {
			log.Printf("Error creating documents directory: %v", err)
			return "", err
		}
	}

	// Write the buffer to a new DOCX file
	documentID := uuid.New().String()
	documentPath := filepath.Join(documentDir, documentID+".docx")
	err = ioutil.WriteFile(documentPath, buf.Bytes(), 0644)
	if err != nil {
		log.Printf("Error writing DOCX file: %v", err)
		return "", err
	}

	// Upload to S3
	s3Path, err := utils.UploadFileToS3(documentPath)
	if err != nil {
		log.Printf("Error uploading file to S3: %v", err)
		return "", err
	}
	log.Printf("Uploaded document to S3 path: %s", s3Path)

	// Save document metadata
	generatedDoc := models.GeneratedDocument{
		ID:           documentID,
		UserID:       userID,
		TemplateID:   templateID,
		DocumentName: documentID + ".docx",
		FilePath:     s3Path,
	}
	if err := utils.DB.Create(&generatedDoc).Error; err != nil {
		log.Printf("Error saving document metadata: %v", err)
		return "", err
	}

	return s3Path, nil
}

func GetGeneratedDocumentsByUserID(userID string) ([]models.GeneratedDocument, error) {
	var documents []models.GeneratedDocument
	err := utils.DB.Where("user_id = ?", userID).Find(&documents).Error
	if err != nil {
		log.Printf("Error getting generated documents: %v", err)
	}
	return documents, err
}

func GetDocumentPathByID(documentID string) (string, error) {
	var document models.GeneratedDocument

	// Fetch the document path from the database
	err := utils.DB.Where("id = ?", documentID).First(&document).Error
	if err != nil {
		log.Printf("Error getting document by ID: %v", err)
		return "", err
	}

	// Sanitize the file path to avoid Windows-specific issues
	document.FilePath = strings.ReplaceAll(document.FilePath, "\\", "/")
	return document.FilePath, nil
}

// func GeneratePresignedURL(fullURL string) (string, error) {

// 	// Extract the S3 key from the full URL
// 	s3Key, err := utils.ExtractS3Key(fullURL)
// 	if err != nil {
// 		return "", fmt.Errorf("failed to extract S3 key: %v", err)
// 	}

// 	// Initialize the AWS session
// 	sess, err := session.NewSession(&aws.Config{
// 		Region: aws.String(config.AppConfig.S3_REGION),
// 		Credentials: credentials.NewStaticCredentials(
// 			config.AppConfig.S3_ACCESS_KEY,
// 			config.AppConfig.S3_SECRET_KEY,
// 			"",
// 		),
// 	})
// 	if err != nil {
// 		log.Printf("Error initializing AWS session: %v", err)
// 		return "", err
// 	}

// 	// Create S3 client
// 	svc := s3.New(sess)

// 	// Log details for debugging
// 	log.Printf("Generating presigned URL for bucket: %s, key: %s", config.AppConfig.S3_BUCKET, s3Key)

// 	// Generate the presigned URL
// 	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
// 		Bucket: aws.String(config.AppConfig.S3_BUCKET),
// 		Key:    aws.String(s3Key),
// 	})

// 	urlStr, err := req.Presign(15 * time.Minute)
// 	if err != nil {
// 		log.Printf("Error generating presigned URL: %v", err)
// 		return "", err
// 	}

// 	log.Printf("Generated presigned URL: %s", urlStr)
// 	return urlStr, nil
// }

func GeneratePresignedURL(fullURL string) (string, error) {
	// Extract the S3 key from the full URL
	s3Key, err := utils.ExtractS3Key(fullURL)
	if err != nil {
		return "", fmt.Errorf("failed to extract S3 key: %v", err)
	}

	// Check if the presigned URL is already in the cache
	if cachedURL, found := utils.Cache.Get(s3Key); found {
		log.Printf("Cache hit: Presigned URL found for key: %s", s3Key)
		return cachedURL.(string), nil
	}

	// Initialize the AWS session
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AppConfig.S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			config.AppConfig.S3_ACCESS_KEY,
			config.AppConfig.S3_SECRET_KEY,
			"",
		),
	})
	if err != nil {
		log.Printf("Error initializing AWS session: %v", err)
		return "", err
	}

	// Create S3 client
	svc := s3.New(sess)

	// Log details for debugging
	log.Printf("Generating presigned URL for bucket: %s, key: %s", config.AppConfig.S3_BUCKET, s3Key)

	// Generate the presigned URL
	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.AppConfig.S3_BUCKET),
		Key:    aws.String(s3Key),
	})

	urlStr, err := req.Presign(15 * time.Minute)
	if err != nil {
		log.Printf("Error generating presigned URL: %v", err)
		return "", err
	}

	// Cache the generated presigned URL
	utils.Cache.Set(s3Key, urlStr, 15*time.Minute)
	log.Printf("Generated and cached presigned URL: %s", urlStr)
	return urlStr, nil
}
