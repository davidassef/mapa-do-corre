package service

import (
	"context"
	"fmt"
	"net/smtp"
	"strings"
)

type EmailMessage struct {
	Destinatario string
	Assunto      string
	Corpo        string
}

type EmailSender interface {
	Enviar(ctx context.Context, message EmailMessage) error
}

type SMTPEmailSender struct {
	address    string
	host       string
	username   string
	password   string
	fromEmail  string
	fromHeader string
}

func NewSMTPEmailSender(host string, port string, username string, password string, fromEmail string, fromName string) *SMTPEmailSender {
	fromHeader := strings.TrimSpace(fromEmail)
	if strings.TrimSpace(fromName) != "" {
		fromHeader = fmt.Sprintf("%s <%s>", strings.TrimSpace(fromName), strings.TrimSpace(fromEmail))
	}

	return &SMTPEmailSender{
		address:    fmt.Sprintf("%s:%s", strings.TrimSpace(host), strings.TrimSpace(port)),
		host:       strings.TrimSpace(host),
		username:   strings.TrimSpace(username),
		password:   password,
		fromEmail:  strings.TrimSpace(fromEmail),
		fromHeader: fromHeader,
	}
}

func (sender *SMTPEmailSender) Enviar(ctx context.Context, message EmailMessage) error {
	if err := ctx.Err(); err != nil {
		return err
	}

	headers := []string{
		fmt.Sprintf("From: %s", sender.fromHeader),
		fmt.Sprintf("To: %s", strings.TrimSpace(message.Destinatario)),
		fmt.Sprintf("Subject: %s", strings.TrimSpace(message.Assunto)),
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
	}

	body := strings.Join(headers, "\r\n") + "\r\n\r\n" + message.Corpo

	var auth smtp.Auth
	if sender.username != "" {
		auth = smtp.PlainAuth("", sender.username, sender.password, sender.host)
	}

	if err := smtp.SendMail(sender.address, auth, sender.fromEmail, []string{strings.TrimSpace(message.Destinatario)}, []byte(body)); err != nil {
		return fmt.Errorf("falha ao enviar email: %w", err)
	}

	return nil
}
