package rabbitmq

import (
	"github.com/streadway/amqp"
	"encoding/json"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
)

type RabbitMQClient struct {
	conn  *amqp.Connection
	ch    *amqp.Channel
	queue amqp.Queue
}

func NewRabbitMQClient(amqpURL, queueName string) (*RabbitMQClient, error) {
	conn, err := amqp.Dial(amqpURL)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		conn.Close()
		return nil, err
	}

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		ch.Close()
		conn.Close()
		return nil, err
	}
	
	return &RabbitMQClient{
		conn:  conn,
		ch:    ch,
		queue: q,
	}, nil
}

func (c *RabbitMQClient) SendLogRequest(item audit.LogItem) error {
	body, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return c.ch.Publish(
		"",
		c.queue.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)
}

func (c *RabbitMQClient) Close() {
	c.ch.Close()
	c.conn.Close()
}