package mailer

import (
	"fmt"
	"net/smtp"
)

type SMTPSender struct {
	Host 	 	string
	Port 	 	string
	User 	 	string
	Password 	string
}

func (s *SMTPSender) Send(msg *Message) error {
	addr := s.Host + ":" + s.Port
	auth := smtp.PlainAuth("", s.User, s.Password, s.Host)

	body := []byte(fmt.Sprintf("From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s", 
	s.User, msg.To, msg.Subject, msg.Body,
	))

	return smtp.SendMail(addr, auth, s.User, []string{msg.To}, body)
}