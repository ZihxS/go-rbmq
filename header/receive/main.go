package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZihxS/go-rbmq/broker"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func run(header *string, channel *amqp.Channel) {
	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare queue"))
	}

	err = channel.ExchangeDeclare("go-rbmq-exchange-logs-header", amqp.ExchangeHeaders, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare exchange"))
	}

	table := amqp.Table{
		"header": *header,
	}

	err = channel.QueueBind(q.Name, "", "go-rbmq-exchange-logs-header", false, table)
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
	header := flag.String("h", "Header", "Header")
	flag.Parse()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	connection, channel := broker.ConnectToRabbitMQ()
	defer broker.CloseConnection(connection, channel)
	go run(header, channel)

	log.Printf("Waiting For Messages. To Exit Please Press CTRL+C.\n")
	<-exitChan

	fmt.Println("")
	log.Println("Shutdown Signal Received!")
	log.Println("Bye Bye!")
}
