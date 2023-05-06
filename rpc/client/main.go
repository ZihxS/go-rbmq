package main

import (
	"context"
	"log"
	"os"

	"github.com/ZihxS/go-rbmq/broker"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func run(channel *amqp.Channel) {
	q, err := channel.QueueDeclare("", false, false, true, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare queue"))
	}

	msgs, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when consume queue"))
	}

	corrID := uuid.NewString()

	err = channel.PublishWithContext(context.Background(), "", "go-rbmq-rpc", false, false, amqp.Publishing{
		ContentType:   "text/plain",
		CorrelationId: corrID,
		ReplyTo:       q.Name,
		Body:          []byte(os.Args[1]),
	})
	if err != nil {
		panic(errors.Wrap(err, "error when publish message"))
	}

	log.Println("Send Message to RPC:", os.Args[1])
	log.Println("Waiting For Response From RPC Server...")

	for msg := range msgs {
		if msg.CorrelationId == corrID {
			log.Println("Result From RPC Server:", string(msg.Body))
			break
		}
	}
}

func main() {
	connection, channel := broker.ConnectToRabbitMQ()
	defer broker.CloseConnection(connection, channel)
	run(channel)
}
