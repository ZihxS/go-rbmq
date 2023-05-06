package broker

import (
	"github.com/pkg/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

func ConnectToRabbitMQ() (*amqp.Connection, *amqp.Channel) {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672")
	if err != nil {
		panic(errors.Wrap(err, "error when connect to RabbitMQ"))
	}

	channel, err := connection.Channel()
	if err != nil {
		panic(errors.Wrap(err, "error when get channel"))
	}

	return connection, channel
}

func CloseConnection(connection *amqp.Connection, channel *amqp.Channel) {
	channel.Close()
	connection.Close()
}
