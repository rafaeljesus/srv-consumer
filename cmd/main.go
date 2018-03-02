package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/oklog/run"
	"github.com/rafaeljesus/srv-consumer/pkg/handler"
	"github.com/rafaeljesus/srv-consumer/pkg/listener"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/message"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/message/amqp"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/stats"
	"github.com/rafaeljesus/srv-consumer/pkg/storage/inmem"
)

func main() {
	conn, ch, err := amqp.NewConnection("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("failed to init rabbit connection: %v", err)
	}
	defer conn.Close()

	store := inmem.New("memory://localhost")
	events := []struct {
		routingKey string
		exchange   string
		handler    message.Handler
	}{
		{
			"user.created",
			"users",
			handler.NewUserCreated(store),
		},
		{
			"user.status.changed",
			"users",
			handler.NewUserStatusChanged(store),
		},
		{
			"user.email.changed",
			"users",
			handler.NewUserEmailChanged(store),
		},
	}

	cancelchan := make(chan struct{})
	var g run.Group
	g.Add(func() error {
		return interrupt(cancelchan)
	}, func(error) {
		close(cancelchan)
	})

	s := new(stats.Client)
	consumer := amqp.NewConsumer(ch)
	for _, e := range events {
		h, err := listener.New(e.routingKey, e.exchange, consumer, e.handler, s)
		if err != nil {
			log.Fatalf("failed to create consumer: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			return h.Run(ctx)
		}, func(error) {
			cancel()
		})
	}

	log.Print("running consumers...")
	if err := g.Run(); err != nil {
		log.Fatalf("failed to run actors group: %v", err)
	}
}

func interrupt(cancel <-chan struct{}) error {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
	select {
	case sig := <-c:
		return fmt.Errorf("received signal %s", sig)
	case <-cancel:
		return errors.New("canceled")
	}
}
