package yasmtp

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/smtp"
	"strings"

	"github.com/go-playground/validator/v10"
)

const MessageIDHeader = "Message-ID"

var validate = validator.New(validator.WithRequiredStructEnabled())

type From struct {
	ServerHostPort string `validate:"hostname_port,required"`
	Password       string `validate:"required"`
	Email          string `validate:"email,required"`
	Name           string
}

type To struct {
	Email string `validate:"email,required"`
	Name  string
}

type Message struct {
	Id      string
	Subject string
	Body    string
}

type Input struct {
	From              From
	To                []To
	Bcc               []To
	Cc                []To
	Msg               Message
	AdditionalHeaders map[string]string
}

// TODO: validate body
// add text/html support/
// add options support (headers/optional text attach etc.)
func SendHTML(ctx context.Context, i *Input) error {
	if err := validate.Struct(i); err != nil {
		return fmt.Errorf("field validation: %w", err)
	}

	if len(i.To)+len(i.Bcc)+len(i.Cc) == 0 {
		return errors.New("no recipients")
	}

	host, _, _ := net.SplitHostPort(i.From.ServerHostPort)
	auth := smtp.PlainAuth("", i.From.Email, i.From.Password, host)
	if i.Msg.Id != "" {
		i.AdditionalHeaders[MessageIDHeader] = i.Msg.Id
	}

	msgBuilder := &strings.Builder{}
	for k, v := range i.AdditionalHeaders {
		wbf(msgBuilder, "%s: %s\r\n", k, v)
	}

	wbf(msgBuilder, "From: \"%s\" <%s>\r\n", i.From.Name, i.From.Email)
	recipients := []string{}
	for _, to := range i.To {
		wbf(msgBuilder, "To: \"%s\"<%s>\r\n", to.Name, to.Email)
		recipients = append(recipients, to.Email)
	}
	for _, cc := range i.Cc {
		wbf(msgBuilder, "Cc: \"%s\"<%s>\r\n", cc.Name, cc.Email)
		recipients = append(recipients, cc.Email)
	}
	for _, bcc := range i.Bcc {
		wbf(msgBuilder, "Bcc: \"%s\"<%s>\r\n", bcc.Name, bcc.Email)
		recipients = append(recipients, bcc.Email)
	}

	wbf(msgBuilder, "MIME-Version: 1.0\r\n")
	wbf(msgBuilder, "Subject: %s\r\n", i.Msg.Subject)
	wbf(msgBuilder, "Content-Type: text/html; charset=utf-8\r\n\r\n")
	wbf(msgBuilder, "\r\n%s\r\n", i.Msg.Body)

	return send(
		ctx,
		i.From.ServerHostPort,
		auth,
		i.From.Email,
		recipients,
		[]byte(msgBuilder.String()),
	)
}

// send is created based on smtp.SendMail, extended by context
func send(ctx context.Context, addr string, a smtp.Auth, from string, to []string, msg []byte) error {
	if err := validateLine(from); err != nil {
		return err
	}
	for _, recp := range to {
		if err := validateLine(recp); err != nil {
			return err
		}
	}
	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	// TODO: verify need of hello
	// if err = c.hello(); err != nil {
	// 	return err
	// }

	host, _, _ := net.SplitHostPort(addr)
	if err = c.StartTLS(&tls.Config{
		ServerName:         host,
		InsecureSkipVerify: false,
	}); err != nil {
		return err
	}
	if err = c.Auth(a); err != nil {
		return err
	}
	if err = c.Mail(from); err != nil {
		return err
	}
	for _, addr := range to {
		if err = c.Rcpt(addr); err != nil {
			return err
		}
	}
	w, err := c.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(msg)
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}

// validateLine checks to see if a line has CR or LF as per RFC 5321.
func validateLine(line string) error {
	if strings.ContainsAny(line, "\n\r") {
		return errors.New("smtp: A line must not contain CR or LF")
	}
	return nil
}

func wbf(b *strings.Builder, f string, args ...any) {
	b.WriteString(fmt.Sprintf(f, args...))
}
