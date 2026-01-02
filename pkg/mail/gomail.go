package mail

import (
	"bytes"
	"fmt"
	"html/template"
	"sync"

	"github.com/Ablebil/lathi-be/internal/config"
	"github.com/Ablebil/lathi-be/pkg/mail/templates"
	"gopkg.in/gomail.v2"
)

type MailItf interface {
	Send(receiver, subject, tmplName string, data any) error
}

type mail struct {
	dialer   *gomail.Dialer
	template *template.Template
}

var (
	instance MailItf
	once     sync.Once
)

func NewMail(env *config.Env) MailItf {
	once.Do(func() {
		tmpl, err := template.ParseFS(templates.TemplateFS, "*.html")
		if err != nil {
			panic(fmt.Errorf("failed to parse email templates: %w", err))
		}

		instance = &mail{
			dialer: gomail.NewDialer(
				env.SMTPHost,
				env.SMTPPort,
				env.SMTPUsername,
				env.SMTPPassword,
			),
			template: tmpl,
		}
	})

	return instance
}

func (m *mail) Send(receiver, subject, tmplName string, data any) error {
	var tmplOutput bytes.Buffer

	err := m.template.ExecuteTemplate(&tmplOutput, tmplName, data)
	if err != nil {
		return fmt.Errorf("failed to execute email template: %w", err)
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", fmt.Sprintf("%s <%s>", "Lathi", m.dialer.Username))
	msg.SetHeader("To", receiver)
	msg.SetHeader("Subject", subject)
	msg.SetBody("text/html", tmplOutput.String())

	return m.dialer.DialAndSend(msg)
}
