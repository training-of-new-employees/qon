// Package doar - пакет для отправки писем пользователям.
package doar

import (
	"fmt"
	"strings"

	"github.com/training-of-new-employees/qon/internal/logger"
	apisender "github.com/training-of-new-employees/qon/internal/pkg/doar/api-sender"
	smtpsender "github.com/training-of-new-employees/qon/internal/pkg/doar/smtp-sender"
	testsender "github.com/training-of-new-employees/qon/internal/pkg/doar/test-sender"
	emailtemplate "github.com/training-of-new-employees/qon/internal/pkg/email-template"
)

type SenderMode string

const (
	TestMode SenderMode = "test"
	ApiMode  SenderMode = "api"
	SmtpMode SenderMode = "smtp"
)

// Интерфейс EmailSender.
type EmailSender interface {
	SendCode(email string, code string) error
	InviteUser(email string, name string, invitationLink string) error
	SendPassword(email string, name string, password string, linkLogin string) error
	Mode() SenderMode
}

type Mailer interface {
	SendEmail(email string, title string, body string) error
}

// Sender описывает структуру для рассылки писем.
type Sender struct {
	mailer      Mailer
	ClientEmail string
	mode        SenderMode
}

// NewSender - конструктор Sender.
func NewSender(config *SenderConfig) *Sender {
	var sender Mailer

	switch config.Mode {
	case SmtpMode: // рассылка с помощью SMTP
		logger.Log.Debug(config.SenderEmail + ": " + config.SenderPassword)
		sender = smtpsender.NewSmtpSender(config.SenderEmail, config.SenderPassword)
	case ApiMode: // рассылка с помощью API-сервиса
		sender = apisender.NewApiSender(config.SenderEmail, config.SenderApiKey)
	case TestMode: // мок-рассылка
		sender = testsender.NewTestSender()
	}

	return &Sender{
		mailer: sender,
		mode:   config.Mode,
	}
}

// SendCode - отправка кода верификации.
func (s *Sender) SendCode(email string, code string) error {
	title := "Подтверждение регистрации (QuickOn)"
	body := fmt.Sprintf("Код верификации: %s", code)

	if s.mode != TestMode {
		var err error
		body, err = emailtemplate.HandleMailTemplate(
			emailtemplate.Verification,
			emailtemplate.Content{VerificationCode: strings.Split(code, "")},
		)
		if err != nil {
			return err
		}
	}

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode - отправка пригласительной ссылки.
func (s *Sender) InviteUser(email string, name string, invitationLink string) error {
	title := "Пригласительная ссылка (QuickOn)"
	body := fmt.Sprintf("Добрый день, %s! Для Вас пригласительная cсылка: %s", name, invitationLink)

	if s.mode != TestMode {
		var err error
		body, err = emailtemplate.HandleMailTemplate(
			emailtemplate.InvitationLink,
			emailtemplate.Content{Name: name, InvitationLink: invitationLink},
		)
		if err != nil {
			return err
		}
	}

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// SendCode - отправка нового пароля.
func (s *Sender) SendPassword(email string, name string, password string, linkLogin string) error {
	title := "Новый пароль (QuickOn)"
	body := fmt.Sprintf("пароль: %s", password)

	if s.mode != TestMode {
		var err error
		body, err = emailtemplate.HandleMailTemplate(
			emailtemplate.PasswordRecovery,
			emailtemplate.Content{Name: name, Password: password, LinkLogin: linkLogin},
		)
		if err != nil {
			return err
		}
	}

	if err := s.mailer.SendEmail(email, title, body); err != nil {
		return err
	}

	return nil
}

// Mode - чтобы узнать, как отправляются письма.
func (s *Sender) Mode() SenderMode {
	return s.mode
}
