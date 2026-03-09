package mail_sender

import (
	"encoding/json"
	"errors"
	"os"
	"simple_twitter/internal/models"
	"simple_twitter/internal/nats_st"
	"simple_twitter/internal/worker"
	"strconv"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
	"gopkg.in/gomail.v2"
)

type MailSenderWorker struct {
	email    string
	passwd   string
	smtpHost string
	smtpPort int
	*gomail.Dialer
	gomail.SendCloser
}

type MailSenderWorkerConfig func(*MailSenderWorker)

func WithEmailEnv(env string) MailSenderWorkerConfig {
	return func(msw *MailSenderWorker) {
		if val := os.Getenv(env); val != "" {
			msw.email = val
		}
	}
}

func WithPasswdEnv(env string) MailSenderWorkerConfig {
	return func(msw *MailSenderWorker) {
		if val := os.Getenv(env); val != "" {
			msw.passwd = val
		}
	}
}

func WithSMTPPHostEnv(env string) MailSenderWorkerConfig {
	return func(msw *MailSenderWorker) {
		if val := os.Getenv(env); val != "" {
			msw.smtpHost = val
		}
	}
}

func WithSMTPPortEnv(env string) MailSenderWorkerConfig {
	return func(msw *MailSenderWorker) {
		if val := os.Getenv(env); val != "" {
			val, err := strconv.Atoi(val)
			if err != nil {
				panic(err)
			}
			msw.smtpPort = val
		}
	}
}

func (m *MailSenderWorker) Consume() {
	nats_st.NATSServer.Conn.Subscribe("mail_sender", func(msg *nats.Msg) {
		var mail models.Mail
		if err := json.Unmarshal(msg.Data, &mail); err != nil {
			log.Err(err)
			return
		}
		message := gomail.NewMessage()
		message.SetHeader("From", m.email)
		message.SetHeader("To", mail.To)
		message.SetHeader("Subject", mail.Subj)
		message.SetBody("text/html", "Hello!")
		if m.Dialer == nil {
			log.Err(errors.New("mail server isn't setup yet")).Msg("error sending email")
			return
		}
		if err := m.Dialer.DialAndSend(message); err != nil {
			log.Err(err).Msg("error sending email")
		}
		log.Info().Any("mail", m).Msg("sending the email")

	})
}

func (m *MailSenderWorker) Connect() error {
	m.Dialer = gomail.NewDialer(
		m.smtpHost,
		m.smtpPort,
		m.email,
		m.passwd,
	)
	sender, err := m.Dialer.Dial()
	if err != nil {
		return err
	}
	m.SendCloser = sender
	return nil
}

func (m *MailSenderWorker) Disconnect() error {
	return m.SendCloser.Close()
}

func NewMailSenderWorker(configs ...MailSenderWorkerConfig) worker.Worker {
	wrk := &MailSenderWorker{
		email:    "test@example.com",
		passwd:   "test",
		smtpHost: "localhost",
		smtpPort: 587,
	}
	for _, config := range configs {
		config(wrk)
	}
	if err := wrk.Connect(); err != nil {
		panic(err)
	}
	return wrk
}

func init() {
	wrk := NewMailSenderWorker(
		WithEmailEnv("SYSTEM_MAIL_AUTH_EMAIL"),
		WithPasswdEnv("SYSTEM_MAIL_AUTH_PASSWORD"),
		WithSMTPPHostEnv("SYSTEM_MAIL_SMTP_HOST"),
		WithSMTPPortEnv("SYSTEM_MAIL_SMTP_PORT"),
	)
	worker.RegisterWorker("mail_sender", wrk)
}
