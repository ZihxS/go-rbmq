package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/ZihxS/go-rbmq/broker"
	"github.com/ZihxS/go-rbmq/utils"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/redis/go-redis/v9"
)

func run(channel *amqp.Channel) {
	q, err := channel.QueueDeclare("go-rbmq-rpc", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare queue"))
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		panic(errors.Wrap(err, "error when set qos"))
	}

	msgs, err := channel.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when consume"))
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Username: "",
		Password: "",
		DB:       0,
	})

	go func() {
		for msg := range msgs {
			ctx := context.Background()

			log.Println("Processing Message:", string(msg.Body))
			response, err := rdb.Get(ctx, string(msg.Body)).Result()

			if errors.Is(err, redis.Nil) {
				time.Sleep(time.Second) // simulation if there is no data in redis yet
				result := utils.IsPalindrome(string(msg.Body))
				response = strconv.FormatBool(result)
				err = rdb.Set(ctx, string(msg.Body), response, 10*time.Second).Err()
				if err != nil {
					panic(errors.Wrap(err, "error when set data to redis"))
				}

			}

			if err != nil {
				panic(errors.Wrap(err, "error when get cache from redis"))
			}

			err = channel.PublishWithContext(context.Background(), "", msg.ReplyTo, false, false, amqp.Publishing{
				ContentType:   "text/plain",
				CorrelationId: msg.CorrelationId,
				Body:          []byte(response),
			})

			if err != nil {
				panic(errors.Wrap(err, "error when publish message"))
			}

			msg.Ack(false)
		}
	}()
}

func main() {
	exitChan := make(chan os.Signal, 1)
	signal.Notify(exitChan, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	connection, channel := broker.ConnectToRabbitMQ()
	defer broker.CloseConnection(connection, channel)
	go run(channel)

	log.Printf("Waiting For RPC Messages. To Exit Please Press CTRL+C.\n")
	<-exitChan

	fmt.Println("")
	log.Println("Shutdown Signal Received!")
	log.Println("Bye Bye!")
}
