package doar

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail() error
}

type Sender struct {
	EmailServiceAddress  string
	EmailServicePassword string
	ClientEmail          string
	CacheKey             string
	EmailTitle           string
	EmailBody            string
}

// TODO
func NewSender(clientEmail string, cacheKey string) *Sender {
	return &Sender{
		EmailServiceAddress:  "",
		EmailServicePassword: "",
		ClientEmail:          clientEmail,
		CacheKey:             cacheKey,
		EmailTitle:           "Подтверждение регистрации",
	}
}

func (s *Sender) SendEmail() error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.EmailServiceAddress)
	m.SetHeader("To", s.ClientEmail)
	m.SetHeader("Subject", s.EmailTitle)
	m.SetBody("text/html", "Код верификации:  "+s.CacheKey)

	d := gomail.NewDialer("smtp.gmail.com", 587, s.EmailServiceAddress, s.EmailServicePassword)
	err := d.DialAndSend(m)
	// TODO error handling
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Email Sent")
	}

	return nil
}
