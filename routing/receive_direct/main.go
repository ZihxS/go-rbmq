package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/ZihxS/go-rbmq/broker"
	"github.com/ZihxS/go-rbmq/utils"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func run(channel *amqp.Channel) {
	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare queue"))
	}

	err = channel.ExchangeDeclare("go-rbmq-exchange-logs-direct", amqp.ExchangeDirect, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare exchange"))
	}

	for _, v := range os.Args[1:] {
		err = channel.QueueBind(q.Name, v, "go-rbmq-exchange-logs-direct", false, nil)
		if err != nil {
			panic(errors.Wrap(err, "error when bind queue"))
		}
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

func osArgsValidation() {
	if len(os.Args) == 1 {
		panic("please use minumum 1 os.Args addition")
	} else {
		route := []string{"info", "warning", "error"}
		for _, v := range os.Args[1:] {
			if !utils.InSlice(v, route) {
				panic(fmt.Sprintf("error when checking os.Args for routing, %v not available in route", v))
			}
		}
	}
}

func main() {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	connection, channel := broker.ConnectToRabbitMQ()
	defer broker.CloseConnection(connection, channel)

	osArgsValidation()
	go run(channel)

	log.Printf("Waiting For Messages. To Exit Please Press CTRL+C.\n")
	<-exitChan

	fmt.Println("")
	log.Println("Shutdown Signal Received!")
	log.Println("Bye Bye!")
}
