package utils

import (
	"fmt"
	"strconv"
	"template-system/config"

	"gopkg.in/gomail.v2"
)

func SendHTMLEmail(to, subject, body string) error {
	from := config.AppConfig.EMAIL_ADDRESS
	password := config.AppConfig.EMAIL_PASSWORD
	smtpHost := config.AppConfig.SMTP_HOST
	smtpPort := config.AppConfig.SMTP_PORT

	// Validate SMTP port
	port, err := strconv.Atoi(smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port: %v", err)
	}

	// Create a new email message
	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	// Create a dialer
	d := gomail.NewDialer(smtpHost, port, from, password)

	// Attempt to send the email
	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email to %s via %s:%d. Error: %v", to, smtpHost, port, err)
	}

	fmt.Printf("Email sent successfully to: %s\n", to)
	return nil
}
