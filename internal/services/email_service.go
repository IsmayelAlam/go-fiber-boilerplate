package services

import (
	"fmt"
	"log/slog"
	"varaden/server/config"

	"gopkg.in/gomail.v2"
)

type EmailService interface {
	SendEmail(to, subject, body string) error
}

type emailService struct {
	Dialer *gomail.Dialer
	From   string
}

func NewEmailService(config *config.SMTPConfig) EmailService {
	return &emailService{
		From: config.From,
		Dialer: gomail.NewDialer(
			config.Host,
			config.Port,
			config.Username,
			config.Password,
		),
	}
}

func (es *emailService) SendEmail(to, subject, body string) error {
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", es.From)
	mailer.SetHeader("To", to)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/plain", body)

	es.background(func() {
		if err := es.Dialer.DialAndSend(mailer); err != nil {
			slog.Error(fmt.Sprintf("Failed to send email to %s: %v", to, err))
		}
	})

	return nil
}

func (es *emailService) background(fn func()) {
	// Increment the WaitGroup counter
	config.SW.Add(1)
	go func() {
		defer func() {
			// Decrement the counter when the goroutine completes
			defer config.SW.Done()
			if err := recover(); err != nil {
				slog.Error(fmt.Sprintf("Panic in background email task: %v", err))
			}
		}()
		fn()
	}()
}
