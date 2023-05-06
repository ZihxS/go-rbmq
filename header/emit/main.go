package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/ZihxS/go-rbmq/broker"
	"github.com/ZihxS/go-rbmq/utils"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func emit(wg *sync.WaitGroup, channel *amqp.Channel, eName string, header, message *string, i int) {
	defer wg.Done()

	msg := fmt.Sprintf("emit for %v: header: %v, message(%v): %v.", eName, *header, i, *message)
	err := channel.PublishWithContext(context.Background(), eName, "", false, false, amqp.Publishing{
		Headers: map[string]any{
			"header": *header,
		},
		ContentType: "text/plain",
		Body:        []byte(msg),
	})
	if err != nil {
		panic(errors.Wrap(err, "error when publish message"))
	}

	log.Println(msg)
}

func run(channel *amqp.Channel, header, message *string) {
	exc := "go-rbmq-exchange-logs-header"
	err := channel.ExchangeDeclare(exc, amqp.ExchangeHeaders, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare exchange"))
	}

	var wg sync.WaitGroup

	for i := 1; i <= utils.Random(20, 50); i++ {
		wg.Add(1)
		go emit(&wg, channel, exc, header, message, i)
	}

	wg.Wait()
	log.Println("Process Done. To Exit Please Press CTRL+C.")
}

func main() {
	header := flag.String("h", "Default Header", "Header")
	message := flag.String("m", "Hai Rabbit!", "Message")
	flag.Parse()

	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	connection, channel := broker.ConnectToRabbitMQ()
	defer broker.CloseConnection(connection, channel)
	go run(channel, header, message)

	<-exitChan

	fmt.Println("")
	log.Println("Shutdown Signal Received!")
	log.Println("Bye Bye!")
}
