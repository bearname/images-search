package amqp

import (
	"github.com/streadway/amqp"
)

const TargetQueue = "QueueService"

func Dial(amqpServerURL string, queueName string) (*amqp.Channel, error) {
	connectRabbitMQ, err := amqp.Dial(amqpServerURL)
	if err != nil {
		panic(err)
	}

	channel, err := connectRabbitMQ.Channel()
	if err != nil {
		panic(err)
	}

	err = initialize(channel, queueName)
	if err != nil {
		return channel, err
	}
	return channel, nil
}

func Consume(channel *amqp.Channel, queueName string) (<-chan amqp.Delivery, error) {
	return channel.Consume(
		queueName, // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
}

func initialize(channel *amqp.Channel, queueName string) error {
	_, err := channel.QueueDeclare(
		queueName, // queue name
		true,      // durable
		false,     // auto delete
		false,     // exclusive
		false,     // no wait
		nil,       // arguments
	)
	if err != nil {
		return err
	}
	return nil
}

func Publish(channel *amqp.Channel, queueName string, message amqp.Publishing) error {
	return channel.Publish(
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,     // immediate
		message,   // message to publish
	)
}
