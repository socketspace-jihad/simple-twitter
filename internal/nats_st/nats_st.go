package nats_st

import (
	"github.com/nats-io/nats.go"
)

var NATSServer *NATS

type NATS struct {
	Addr string
	*nats.Conn
}

func (n *NATS) Connect() error {
	conn, err := nats.Connect("nats://localhost:4222")
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

func init() {
	NATSServer = NewNATS()
	if err := NATSServer.Connect(); err != nil {
		panic(err)
	}
	if !NATSServer.IsConnected() {
		panic("NATS server isn't ready")
	}
}
