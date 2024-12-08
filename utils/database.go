package utils

import (
	"fmt"
	"log"
	"template-system/config"
	"template-system/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=Asia/Shanghai",
		config.AppConfig.DB_HOST, config.AppConfig.DB_USER, config.AppConfig.DB_PASSWORD, config.AppConfig.DB_NAME, config.AppConfig.DB_PORT)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	err = DB.AutoMigrate(&models.User{}, &models.Template{}, &models.GeneratedDocument{})
	if err != nil {
		log.Fatal("Failed to migrate database:", err)
	}
}
