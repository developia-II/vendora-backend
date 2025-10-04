package utils

import (
	"fmt"
	"os"

	"gopkg.in/gomail.v2"
)

func SendEmail(to, subject, body string) error {
	smtpHost := os.Getenv("SMTP_HOST")
	smtpPort := 587
	smtpUser := os.Getenv("SMTP_USER") // Keep as 886548001@smtp-brevo.com
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	// Use your verified email as sender
	senderEmail := "opiafavourjr@gmail.com"
	senderName := "Vendora"

	if smtpHost == "" || smtpUser == "" || smtpPassword == "" {
		return fmt.Errorf("SMTP configuration missing")
	}

	m := gomail.NewMessage()
	m.SetAddressHeader("From", senderEmail, senderName) // "Vendora <yourname@gmail.com>"
	m.SetHeader("To", to)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(smtpHost, smtpPort, smtpUser, smtpPassword)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}

	return nil
}
