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
	"github.com/rafaeljesus/srv-consumer/pkg/consumer"
	"github.com/rafaeljesus/srv-consumer/pkg/handler"
	"github.com/rafaeljesus/srv-consumer/pkg/listener"
	"github.com/rafaeljesus/srv-consumer/pkg/platform/stats"
	"github.com/rafaeljesus/srv-consumer/pkg/storage/inmem"
	"github.com/streadway/amqp"
)

func main() {
	conn, ch, err := createConnAndChan("amqp://guest:guest@localhost:5672")
	if err != nil {
		log.Fatalf("failed to init rabbit connection: %v", err)
	}
	defer conn.Close()

	s := new(stats.Client)
	lner := listener.New(ch)
	store := inmem.New("memory://localhost")
	events := []struct {
		routingKey string
		exchange   string
		handler    consumer.Handler
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

	var g run.Group
	ctx, cancel := context.WithCancel(context.Background())

	cancelchan := make(chan struct{})
	g.Add(func() error {
		return interrupt(cancelchan)
	}, func(error) {
		close(cancelchan)
	})

	for _, e := range events {
		h, err := consumer.New(e.routingKey, e.exchange, lner, e.handler, s)
		if err != nil {
			// TODO cancel ctx
			log.Fatalf("failed to create consumer: %v", err)
		}

		g.Add(func() error {
			return h.Run(ctx)
		}, func(error) {
			cancel()
		})
	}

	log.Print("running consumers...")
	g.Run()
}

func createConnAndChan(dsn string) (*amqp.Connection, *amqp.Channel, error) {
	conn, err := amqp.Dial(dsn)
	if err != nil {
		return nil, nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, nil, err
	}

	return conn, ch, nil
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
