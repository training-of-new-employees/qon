// Package test-sender - пакет для мок-рассылки писем.
package testsender

import (
	"fmt"
)

// TestSender - test-sender.
type TestSender struct {
}

// NewTestSender - конструктор.
func NewTestSender() *TestSender {
	return &TestSender{}
}

// SendEmail - вывод содержания письма пользователю.
func (s *TestSender) SendEmail(email string, title string, body string) error {
	return fmt.Errorf("<mock-sender>: email: %s; subject: %s; body: %s", email, title, body)
}
