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
	"github.com/ZihxS/go-rbmq/utils"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func emit(wg *sync.WaitGroup, channel *amqp.Channel, eName string, parent int, delay time.Duration, amount int, topic string) {
	defer wg.Done()
	for i := 1; i <= amount; i++ {
		go func() {
			switch topic {
			case "lazy.#":
				randCond := utils.Random(1, 2)
				if randCond == 1 {
					topic = "lazy" + "." + "orange"
				} else {
					topic = "lazy" + "." + uuid.NewString()
				}
			case "*.*.rabbit":
				randCond := utils.Random(1, 4)
				switch randCond {
				case 1:
					topic = "lazy" + "." + uuid.NewString() + "." + "rabbit"
				case 2:
					topic = uuid.NewString() + "." + "orange" + "." + "rabbit"
				case 3:
					topic = "lazy" + "." + "orange" + "." + "rabbit"
				case 4:
					topic = uuid.NewString() + "." + uuid.NewString() + "." + "rabbit"
				}
			case "*.orange.*":
				randCond := utils.Random(1, 4)
				switch randCond {
				case 1:
					topic = "lazy" + "." + "orange" + "." + uuid.NewString()
				case 2:
					topic = uuid.NewString() + "." + "orange" + "." + "rabbit"
				case 3:
					topic = "lazy" + "." + "orange" + "." + "rabbit"
				case 4:
					topic = uuid.NewString() + "." + "orange" + "." + uuid.NewString()
				}
			}
			msg := fmt.Sprintf("emit for %v (%v): this is message from parent: %v, child: %v.", eName, topic, parent, i)
			err := channel.PublishWithContext(context.Background(), eName, topic, false, false, amqp.Publishing{
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
	exc := "go-rbmq-exchange-logs-topic"
	err := channel.ExchangeDeclare(exc, amqp.ExchangeTopic, true, false, false, false, nil)
	if err != nil {
		panic(errors.Wrap(err, "error when declare exchange"))
	}

	var wg sync.WaitGroup

	// for emit to topic: lazy.#
	min, max := 1, utils.Random(2, 5)

	for i := min; i <= max; i++ {
		wg.Add(1)
		go emit(&wg, channel, exc, i, ((time.Second / 4) * time.Duration(i)), ((max+1)*50)-(15*i), "lazy.#")
	}

	// for emit to topic: *.*.rabbit
	max = utils.Random(2, 4)

	for i := min; i <= max; i++ {
		wg.Add(1)
		go emit(&wg, channel, exc, i, ((time.Second / 4) * time.Duration(i)), ((max+1)*50)-(15*i), "*.*.rabbit")
	}

	// for emit to topic: *.orange.*
	max = utils.Random(2, 3)

	for i := min; i <= max; i++ {
		wg.Add(1)
		go emit(&wg, channel, exc, i, ((time.Second / 4) * time.Duration(i)), ((max+1)*50)-(15*i), "*.orange.*")
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
