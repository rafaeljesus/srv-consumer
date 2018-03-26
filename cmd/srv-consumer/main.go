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
	"github.com/rafaeljesus/srv-consumer/handler"
	"github.com/rafaeljesus/srv-consumer/platform/message"
	"github.com/rafaeljesus/srv-consumer/platform/message/amqp"
	"github.com/rafaeljesus/srv-consumer/platform/stats"
	"github.com/rafaeljesus/srv-consumer/register"
	"github.com/rafaeljesus/srv-consumer/storage/inmem"
)

func main() {
	amqpDSN := os.Getenv("AMQP_DSN")
	if amqpDSN == "" {
		amqpDSN = "amqp://guest:guest@localhost:5672"
	}
	conn, ch, err := amqp.NewConnection(amqpDSN)
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

	sts := new(stats.Client)
	consumer := amqp.NewConsumer(ch)
	for _, e := range events {
		reg, err := register.New(e.routingKey, e.exchange, consumer, e.handler, sts)
		if err != nil {
			log.Fatalf("failed to create consumer: %v", err)
		}

		ctx, cancel := context.WithCancel(context.Background())
		g.Add(func() error {
			return reg.Run(ctx)
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
