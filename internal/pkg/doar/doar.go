// Package doar - пакет для отправки писем пользователям.
package doar

import (
	"fmt"

	"github.com/training-of-new-employees/qon/internal/config"
	"github.com/training-of-new-employees/qon/internal/logger"
	apisender "github.com/training-of-new-employees/qon/internal/pkg/doar/api-sender"
	smtpsender "github.com/training-of-new-employees/qon/internal/pkg/doar/smtp-sender"
	testsender "github.com/training-of-new-employees/qon/internal/pkg/doar/test-sender"
)

// Интерфейс EmailSender.
type EmailSender interface {
	SendCode(email string, code string) error
	InviteUser(email string, invitationLink string) error
	SendPassword(email string, password string) error
	Mode() string
}

type Mailer interface {
	SendEmail(email string, title string, body string) error
}

// Sender описывает структуру для рассылки писем.
type Sender struct {
	mailer      Mailer
	ClientEmail string
	mode        string
}

// NewSender - конструктор Sender.
func NewSender(mode string, config *config.Config) *Sender {
	var sender Mailer
	if mode == "smtp" {
		logger.Log.Debug(config.SenderEmail + ": " + config.SenderPassword)
		sender = smtpsender.NewSmtpSender(config.SenderEmail, config.SenderPassword)
	}
	if mode == "api" {
		sender = apisender.NewApiSender(config.SenderEmail, config.SenderApiKey)
	}
	if mode == "test" {
		sender = testsender.NewTestSender()
	}

	return &Sender{
		mailer: sender,
		mode:   mode,
	}
}

// SendCode - отправка кода верификации.
func (s *Sender) SendCode(email string, code string) error {
	title := "Подтверждение регистрации (QuickOn)"
	body := fmt.Sprintf("Код верификации: %s", code)

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode - отправка пригласительной ссылки.
func (s *Sender) InviteUser(email string, invitationLink string) error {
	title := "Пригласительная ссылка (QuickOn)"
	body := fmt.Sprintf("Пригласительная cсылка: %s", invitationLink)

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode - отправка нового пароля.
func (s *Sender) SendPassword(email string, password string) error {
	title := "Новый пароль (QuickOn)"
	body := fmt.Sprintf("пароль: %s", password)

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// Mode - чтобы узнать, как отправляются письма.
func (s *Sender) Mode() string {
	return s.mode
}
