// Package api-sender - пакет для рассылки писем через api.
package apisender

import (
	"context"

	"github.com/go-resty/resty/v2"
	"github.com/training-of-new-employees/qon/internal/logger"
)

// ApiSender - api-sender.
type ApiSender struct {
	senderEmail string
	apiKey      string
	client      *resty.Client
}

// NewApiSender - конструктор.
func NewApiSender(senderEmail string, key string) *ApiSender {
	return &ApiSender{
		senderEmail: senderEmail,
		apiKey:      key,
		client:      resty.New(),
	}
}

// SendEmail - отправить письмо.
func (s *ApiSender) SendEmail(email string, title string, body string) error {
	res, err := s.client.R().
		SetQueryParams(map[string]string{
			"api_key":        s.apiKey,
			"sender_email":   s.senderEmail,
			"email":          email,
			"subject":        title,
			"body":           body,
			"sender_name":    "QuickOn",
			"format":         "json",
			"error_checking": "1",
			"list_id":        "1",
		}).
		SetContext(context.TODO()).Get("https://api.unisender.com/ru/api/sendEmail")

	if err != nil {
		return err
	}

	logger.Log.Warn(string(res.Body()))

	return nil
}
