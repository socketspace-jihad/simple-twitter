package mail_sender

import (
	"encoding/json"
	"simple_twitter/internal/models"
	"simple_twitter/internal/nats_st"
	"simple_twitter/internal/worker"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

type MailSenderWorker struct{}

func (m *MailSenderWorker) Consume() {
	nats_st.NATSServer.Conn.Subscribe("mail_sender", func(msg *nats.Msg) {
		var m models.Mail
		if err := json.Unmarshal(msg.Data, &m); err != nil {
			log.Err(err)
			return
		}
		log.Info().Any("mail", m).Msg("sending the email")
	})
}

func NewMailSenderWorker() worker.Worker {
	return &MailSenderWorker{}
}

func init() {
	worker.RegisterWorker("mail_sender", NewMailSenderWorker())
}
