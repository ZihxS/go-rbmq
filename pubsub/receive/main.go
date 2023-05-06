package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZihxS/go-rbmq/broker"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func run(channel *amqp.Channel) {
	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare queue"))
	}

	err = channel.ExchangeDeclare("go-rbmq-exchange-logs", amqp.ExchangeFanout, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare exchange"))
	}

	err = channel.QueueBind(q.Name, "", "go-rbmq-exchange-logs", false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when bind queue"))
	}

	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when consume queue"))
	}

	go func() {
		for msg := range msgs {
			log.Printf("Received a message: %s", msg.Body)
		}
	}()
}

func main() {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	connection, channel := broker.ConnectToRabbitMQ()
	defer broker.CloseConnection(connection, channel)
	go run(channel)

	log.Printf("Waiting For Messages. To Exit Please Press CTRL+C.\n")
	<-exitChan

	fmt.Println("")
	log.Println("Shutdown Signal Received!")
	log.Println("Bye Bye!")
}
