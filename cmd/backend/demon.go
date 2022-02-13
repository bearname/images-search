package main

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
	"github.com/streadway/amqp"
	rabbitmq "photofinish/pkg/common/infrarstructure/amqp"
	"photofinish/pkg/domain/broker"
	"photofinish/pkg/domain/tasks"
	"time"
)

func handleDemon(outboxRepo broker.Repo, amqpChannel *amqp.Channel) {
	for {
		list, err := outboxRepo.FindNotCompletedOutboxList(10)
		if err != nil {
			log.Println(err)
		}
		publish(outboxRepo, list, amqpChannel)
		time.Sleep(1 * time.Minute)
	}
}

func publish(outboxRepo broker.Repo, list *[]broker.Outbox, amqpChannel *amqp.Channel) {
	var task tasks.Task
	var t tasks.AddImageDto
	var data []byte
	var err error
	var status broker.ProcessingStatus
	for _, outbox := range *list {
		value := outbox.BrokerValue
		err = json.Unmarshal([]byte(value), &task)
		if err != nil {
			err = outboxRepo.UpdateStatus(outbox.Id.String(), broker.OutboxNotProcessing)
			continue
		}
		t.BrokerTopic = outbox.BrokerTopic
		t.TaskData = value

		data, err = json.Marshal(t)
		if err != nil {
			err = outboxRepo.UpdateStatus(outbox.Id.String(), broker.OutboxNotProcessing)
			continue
		}
		err = rabbitmq.PublishToQueue(amqpChannel, outbox.BrokerTopic, data)
		if err != nil {
			time.Sleep(2 * time.Second)
			err = rabbitmq.PublishToQueue(amqpChannel, outbox.BrokerTopic, data)
			if err != nil {
				status = broker.OutboxNotProcessing
			}
		}

		status = broker.OutboxProcessing
		err = outboxRepo.UpdateStatus(outbox.Id.String(), status)
		if err != nil {
			continue
		}
	}
}
