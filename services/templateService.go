package services

import (
	"archive/zip"
	"bytes"
	"fmt"

	// "fmt"
	"regexp"
	"strings"
	"template-system/models"
	"template-system/utils"
	// "github.com/unidoc/unioffice/document"
)

func CreateTemplate(template *models.Template) error {
	return utils.DB.Create(template).Error
}

func GetTemplatesByUserID(userID string) ([]models.Template, error) {
	var templates []models.Template
	err := utils.DB.Where("user_id = ?", userID).Find(&templates).Error
	return templates, err
}

func ExtractPlaceholders(templateID string) ([]string, error) {
	var tmpl models.Template
	if err := utils.DB.Where("id = ?", templateID).First(&tmpl).Error; err != nil {
		return nil, err
	}

	// Open the DOCX file
	r, err := zip.OpenReader(tmpl.FilePath)
	if err != nil {
		return nil, err
	}
	defer r.Close()

	var content string
	// Loop through files in the DOCX archive
	for _, file := range r.File {
		if strings.HasSuffix(file.Name, "document.xml") {
			// Open document.xml file
			rc, err := file.Open()
			if err != nil {
				return nil, err
			}

			buf := new(bytes.Buffer)
			buf.ReadFrom(rc)
			content = buf.String()
			rc.Close()
			break
		}
	}

	if content == "" {
		return nil, fmt.Errorf("document.xml not found in DOCX file")
	}

	// fmt.Printf("content is: %s\n", content)

	// Ensure that we're reading the correct format of placeholders, e.g., {{placeholder}}
	re := regexp.MustCompile(`{{\s*([a-zA-Z0-9_]+)\s*}}`)
	matches := re.FindAllStringSubmatch(content, -1)

	placeholders := []string{}
	for _, match := range matches {
		placeholders = append(placeholders, match[1])
	}

	return placeholders, nil
}
