// Package test-sender - пакет для мок-рассылки писем.
package testsender

import (
	"errors"
	"fmt"

	"github.com/go-resty/resty/v2"
	"github.com/training-of-new-employees/qon/internal/errs"
)

// TestSender - test-sender.
type TestSender struct {
	senderEmail string
	apiKey      string
	client      *resty.Client
}

// NewTestSender - конструктор.
func NewTestSender() *TestSender {
	return &TestSender{}
}

// SendEmail - вывод содержания письма пользователю.
func (s *TestSender) SendEmail(email string, title string, body string) error {
	errs.ErrNotSendEmail = errors.New(fmt.Sprintf("email: %s; subject: %s; body: %s", email, title, body))

	return errs.ErrNotSendEmail
}
