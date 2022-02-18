package amqp

import (
	"github.com/streadway/amqp"
)

const TargetQueue = "QueueService"

type BrokerService interface {
	Consume(queueName string) (<-chan amqp.Delivery, error)
	PublishToQueue(queueName string, data []byte) error
}

type Service struct {
	channel *amqp.Channel
}

func NewAmqpService(channel *amqp.Channel) *Service {
	s := new(Service)
	s.channel = channel
	return s
}
func (s *Service) Consume(queueName string) (<-chan amqp.Delivery, error) {
	return s.channel.Consume(
		queueName, // queue name
		"",        // consumer
		true,      // auto-ack
		false,     // exclusive
		false,     // no local
		false,     // no wait
		nil,       // arguments
	)
}

func (s *Service) PublishToQueue(queueName string, data []byte) error {
	message := amqp.Publishing{
		ContentType: "text/plain",
		Body:        data,
	}

	err := s.publish(s.channel, queueName, message)
	if err != nil {
		return err
	}
	return err
}

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

func (s *Service) publish(channel *amqp.Channel, queueName string, message amqp.Publishing) error {
	return channel.Publish(
		"",        // exchange
		queueName, // queue name
		false,     // mandatory
		false,     // immediate
		message,   // message to publish
	)
}
