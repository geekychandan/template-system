package utils

import (
	"bytes"
	"os"
	"template-system/config"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func UploadFileToS3(filePath string) (string, error) {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(config.AppConfig.S3_REGION),
		Credentials: credentials.NewStaticCredentials(
			config.AppConfig.S3_ACCESS_KEY, config.AppConfig.S3_SECRET_KEY, ""),
	})
	if err != nil {
		return "", err
	}

	uploader := s3manager.NewUploader(sess)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	buffer := new(bytes.Buffer)
	buffer.ReadFrom(file)

	result, err := uploader.Upload(&s3manager.UploadInput{
		Bucket: aws.String(config.AppConfig.S3_BUCKET),
		Key:    aws.String(filePath),
		Body:   buffer,
	})
	if err != nil {
		return "", err
	}

	return result.Location, nil
}
