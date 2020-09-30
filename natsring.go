package natsring

import (
	nats "github.com/nats-io/nats.go"
	"github.com/serialx/hashring"
)

type NatsRing struct {
	nodes []string
	ring  *hashring.HashRing
	conns map[string]*nats.Conn
}

func New(nodes []string) *NatsRing {
	ring := hashring.New(nodes)
	return &NatsRing{
		nodes: nodes,
		ring:  ring,
		conns: make(map[string]*nats.Conn),
	}
}

func (nr *NatsRing) ConnectAll(options ...nats.Option) error {
	for _, node := range nr.nodes {
		conn, err := nats.Connect(node, options...)
		if err != nil {
			nr.CloseAll()
			return err
		}
		nr.conns[node] = conn
	}
	return nil
}

func (nr *NatsRing) CloseAll() {
	for _, conn := range nr.conns {
		conn.Close()
	}
}

func (nr *NatsRing) Publish(subj string, data []byte) error {
	node, _ := nr.ring.GetNode(subj)
	conn, _ := nr.conns[node]
	return conn.Publish(subj, data)
}

func (nr *NatsRing) Subscribe(subj string, cb nats.MsgHandler) (*nats.Subscription, error) {
	node, _ := nr.ring.GetNode(subj)
	conn, _ := nr.conns[node]
	return conn.Subscribe(subj, cb)
}

func (nr *NatsRing) QueueSubscribe(subj, queue string, cb nats.MsgHandler) (*nats.Subscription, error) {
	node, _ := nr.ring.GetNode(subj)
	conn, _ := nr.conns[node]
	return conn.QueueSubscribe(subj, queue, cb)
}
