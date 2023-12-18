// Package doar - пакет для отправки писем пользователям (smpt).
package smtpsender

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// SmtpSender - smtp-sender.
type SmtpSender struct {
	EmailServiceAddress  string
	EmailServicePassword string
}

func NewSmtpSender(serviceEmail string, password string) *SmtpSender {
	return &SmtpSender{
		EmailServiceAddress:  serviceEmail,
		EmailServicePassword: password,
	}
}

// SendEmail - отправить письмо.
func (s *SmtpSender) SendEmail(email string, title string, body string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.EmailServiceAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", title)
	m.SetBody("text/html", body)

	d := gomail.NewDialer("smtp.gmail.com", 587, s.EmailServiceAddress, s.EmailServicePassword)
	err := d.DialAndSend(m)
	if err != nil {
		return err
	}

	fmt.Println("Email Sent")

	return nil
}
