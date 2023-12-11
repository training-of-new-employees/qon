// Package doar - пакет для отправки писем пользователям.
package doar

import (
	"fmt"

	"gopkg.in/gomail.v2"
)

// Интерфейс EmailSender.
type EmailSender interface {
	SendCode(email string, code string) error
	InviteUser(email string, invitationLink string) error
	SendPassword(email string, password string) error
}

// Sender описывает структуру для рассылки писем.
type Sender struct {
	EmailServiceAddress  string
	EmailServicePassword string
	ClientEmail          string
}

// NewSender - конструктор Sender.
func NewSender(serviceEmail string, password string) *Sender {
	return &Sender{
		EmailServiceAddress:  serviceEmail,
		EmailServicePassword: password,
		ClientEmail:          "",
	}
}

// SendCode отправляет код верификации.
func (s *Sender) SendCode(email string, code string) error {
	title := "Подтверждение регистрации (QuickOn)"
	body := fmt.Sprintf("Код верификации: %s", code)

	if err := s.sendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode отправляет пригласительную ссылку.
func (s *Sender) InviteUser(email string, invitationLink string) error {
	title := "Пригласительная ссылка (QuickOn)"
	body := fmt.Sprintf("пригласительная cсылка: %s", invitationLink)

	if err := s.sendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode отправляет новый пароль.
func (s *Sender) SendPassword(email string, password string) error {
	title := "Новый пароль (QuickOn)"
	body := fmt.Sprintf("пароль: %s", password)

	if err := s.sendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// sendEmail - для отправки писем.
func (s *Sender) sendEmail(email string, title string, body string) error {
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
