package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"github.com/ZihxS/go-rbmq/broker"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func produce(wg *sync.WaitGroup, channel *amqp.Channel, qName string, parent int, delay time.Duration, amount int) {
	defer wg.Done()
	for i := 1; i <= amount; i++ {
		go func() {
			msg := fmt.Sprintf("from producer for %v: this is message from parent: %v, child: %v.", qName, parent, i)
			err := channel.PublishWithContext(context.Background(), "", qName, false, false, amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(msg),
			})
			if err != nil {
				panic(errors.Wrap(err, "error when publish message"))
			}
			log.Println(msg)
		}()
		time.Sleep(delay)
	}
}

func run(channel *amqp.Channel) {
	q, err := channel.QueueDeclare("go-rbmq-queue-basic", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare queue"))
	}

	var wg sync.WaitGroup

	min, max := 1, 5

	for i := min; i <= max; i++ {
		wg.Add(1)
		go produce(&wg, channel, q.Name, i, ((time.Second / 4) * time.Duration(i)), ((max+1)*50)-(15*i))
	}

	wg.Wait()
	log.Println("Process Done. To Exit Please Press CTRL+C.")
}

func main() {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	connection, channel := broker.ConnectToRabbitMQ()
	defer broker.CloseConnection(connection, channel)
	go run(channel)

	<-exitChan

	fmt.Println("")
	log.Println("Shutdown Signal Received!")
	log.Println("Bye Bye!")
}
