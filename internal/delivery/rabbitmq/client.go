package rabbitmq

import (
	"context"
	"encoding/json"
	"valerii/crudbananas/internal/domain"

	"github.com/streadway/amqp"
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

func (c *RabbitMQClient) SendLogRequest(ctx context.Context, item audit.LogItem) error {
	ctxMap := make(map[string]interface{})

    if metadata, ok := ctx.Value("metadata").(map[string]interface{}); ok {
        for key, value := range metadata {
            ctxMap[key] = value
        }
    }

	message := domain.LogMessage{
		Context: ctxMap,
		LogItem: item,
	}

	body, err := json.Marshal(message)
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