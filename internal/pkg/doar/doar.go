// Package doar - пакет для отправки писем пользователям.
package doar

import (
	"fmt"

	"github.com/training-of-new-employees/qon/internal/config"
	apisender "github.com/training-of-new-employees/qon/internal/pkg/doar/api-sender"
	smtpsender "github.com/training-of-new-employees/qon/internal/pkg/doar/smtp-sender"
)

// Интерфейс EmailSender.
type EmailSender interface {
	SendCode(email string, code string) error
	InviteUser(email string, invitationLink string) error
	SendPassword(email string, password string) error
}

type Mailer interface {
	SendEmail(email string, title string, body string) error
}

// Sender описывает структуру для рассылки писем.
type Sender struct {
	mailer      Mailer
	ClientEmail string
}

// NewSender - конструктор Sender.
func NewSender(mode string, config *config.Config) *Sender {
	var sender Mailer
	if mode == "smtp" {
		sender = smtpsender.NewSmtpSender(config.SenderEmail, config.SenderPassword)
	}
	if mode == "api" {
		sender = apisender.NewApiSender(config.SenderEmail, config.SenderApiKey)
	}

	return &Sender{
		mailer: sender,
	}
}

// SendCode отправляет код верификации.
func (s *Sender) SendCode(email string, code string) error {
	title := "Подтверждение регистрации (QuickOn)"
	body := fmt.Sprintf("Код верификации: %s", code)

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode отправляет пригласительную ссылку.
func (s *Sender) InviteUser(email string, invitationLink string) error {
	title := "Пригласительная ссылка (QuickOn)"
	body := fmt.Sprintf("Пригласительная cсылка: %s", invitationLink)

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode отправляет новый пароль.
func (s *Sender) SendPassword(email string, password string) error {
	title := "Новый пароль (QuickOn)"
	body := fmt.Sprintf("пароль: %s", password)

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}
