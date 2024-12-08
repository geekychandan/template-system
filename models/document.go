package models

import (
	"time"

	"gorm.io/gorm"
)

type GeneratedDocument struct {
	ID           string         `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID       string         `gorm:"type:uuid;not null"`
	TemplateID   string         `gorm:"type:uuid;not null"`
	DocumentName string         `gorm:"not null"`
	FilePath     string         `gorm:"not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	DeletedAt    gorm.DeletedAt `gorm:"index"`
}
