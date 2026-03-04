package nats_st

import (
	"os"

	"github.com/nats-io/nats.go"
	"github.com/rs/zerolog/log"
)

var NATSServer *NATS

type NATS struct {
	Addr string
	*nats.Conn
}

func (n *NATS) Connect() error {
	conn, err := nats.Connect(n.Addr)
	if err != nil {
		return err
	}
	n.Conn = conn
	return nil
}

func (n *NATS) Disconnect() {
	n.Conn.Close()
}

func (n *NATS) Publish(topic string, data []byte) error {
	return n.Conn.Publish(topic, data)
}

type natsConfig func(*NATS)

func NewNATS(configs ...natsConfig) *NATS {
	n := &NATS{
		Addr: "nats://localhost:4222",
	}
	for _, config := range configs {
		config(n)
	}
	return n
}

func WithAddrEnv(env string) natsConfig {
	return func(n *NATS) {
		if val := os.Getenv(env); val != "" {
			n.Addr = val
		}
	}
}

func init() {
	NATSServer = NewNATS(
		WithAddrEnv("NATS_URL"),
	)
	if err := NATSServer.Connect(); err != nil {
		panic(err)
	}
	if !NATSServer.IsConnected() {
		panic("NATS server isn't ready")
	}
	log.Info().Msg("NATS URL SETUP")
}
