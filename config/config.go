package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DB_HOST       string
	DB_PORT       string
	DB_USER       string
	DB_PASSWORD   string
	DB_NAME       string
	S3_BUCKET     string
	S3_REGION     string
	S3_ACCESS_KEY string
	S3_SECRET_KEY string
	JWT_SECRET    string
}

var AppConfig Config

func InitConfig() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	AppConfig = Config{
		DB_HOST:       os.Getenv("DB_HOST"),
		DB_PORT:       os.Getenv("DB_PORT"),
		DB_USER:       os.Getenv("DB_USER"),
		DB_PASSWORD:   os.Getenv("DB_PASSWORD"),
		DB_NAME:       os.Getenv("DB_NAME"),
		S3_BUCKET:     os.Getenv("S3_BUCKET"),
		S3_REGION:     os.Getenv("S3_REGION"),
		S3_ACCESS_KEY: os.Getenv("S3_ACCESS_KEY"),
		S3_SECRET_KEY: os.Getenv("S3_SECRET_KEY"),
		JWT_SECRET:    os.Getenv("JWT_SECRET"),
	}
}
