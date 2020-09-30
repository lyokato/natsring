package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/lyokato/natsring"
	"github.com/nats-io/nats.go"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

func main() {

	nodes := []string{
		"localhost:4222",
		"localhost:4223",
		"localhost:4224",
		"localhost:4226",
	}

	channels := makeChannels(300)

	pubStopCh := make(chan struct{})
	subStopCh := make(chan struct{})

	go publisher(nodes, channels, pubStopCh)
	go subscriber(nodes, channels, subStopCh)

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit

	close(pubStopCh)
	close(subStopCh)
}

func makeChannels(num int) []string {
	channels := make([]string, num)
	for i := 0; i < num; i++ {
		channels[i] = fmt.Sprintf("%s:%d", randomString(10), i)
	}
	return channels
}

var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func randomString(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func publisher(nodes, channels []string, stopper <-chan struct{}) {

	nr := natsring.New(nodes)
	nr.ConnectAll(nats.MaxReconnects(10))
	// Subscriber側がチャンネルを全部subするのを待つ

	time.Sleep(time.Second * 1)

	idx := 0

	for {
		select {
		case <-stopper:
			nr.CloseAll()
			return
		default:
		}

		if idx >= len(channels) {
			return
		}

		ch := channels[idx]
		nr.Publish(ch, []byte("DUMMY_MESSAGE"))

		idx++
		time.Sleep(time.Millisecond * 50)
	}

}

func subscriber(nodes, channels []string, stopper <-chan struct{}) {

	nr := natsring.New(nodes)
	nr.ConnectAll(nats.MaxReconnects(10))

	// 指定されたワードを最初に全部subする
	for _, ch := range channels {
		nr.Subscribe(ch, func(m *nats.Msg) {
			fmt.Printf("RECEIVED: %s: %s\n", m.Subject, string(m.Data))
		})
	}

	for {
		select {
		case <-stopper:
			nr.CloseAll()
			return
		default:
		}
	}

	nr.CloseAll()
}
