package doar

import (
	"fmt"
	"gopkg.in/gomail.v2"
)

type EmailSender interface {
	SendEmail(email string, key string) error
	InviteUser(email string) error
}

type Sender struct {
	EmailServiceAddress  string
	EmailServicePassword string
	ClientEmail          string
	AdminTitle           string
	AdminBody            string
}

func NewSender(serviceEmail string, password string) *Sender {
	return &Sender{
		EmailServiceAddress:  serviceEmail,
		EmailServicePassword: password,
		ClientEmail:          "",
		AdminTitle:           "Подтверждение регистрации",
		AdminBody:            "Код верификации:  ",
	}
}

func (s *Sender) SendEmail(email string, key string) error {
	m := gomail.NewMessage()
	m.SetHeader("From", s.EmailServiceAddress)
	m.SetHeader("To", email)
	m.SetHeader("Subject", s.AdminTitle)
	m.SetBody("text/html", s.AdminBody+key)

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

func (s *Sender) InviteUser(email string) error {
	return nil
}
