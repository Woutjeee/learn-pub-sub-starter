package pubsub

import (
	"context"
	"encoding/json"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	b, err := json.Marshal(val)
	if err != nil {
		panic(err)
	}

	x := amqp.Publishing{
		ContentType: "application/json",
		Body:        b,
	}

	ch.PublishWithContext(context.Background(), exchange, key, false, false, x)
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueuName,
	key string,
	simpleQueueType int,
) (*amqp.Channel, amqp.Queue, error) {
	c, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, errors.New("Something went wrong making a channel.")
	}

	var isDurable, autoDelete, exclusive bool
	if simpleQueueType == 0 {
		isDurable = true
		autoDelete = false
		exclusive = false
	} else {
		isDurable = false
		autoDelete = true
		exclusive = true
	}

	q, err := c.QueueDeclare(
		queueuName,
		isDurable,
		autoDelete,
		exclusive,
		false,
		nil,
	)

	c.QueueBind(q.Name, key, exchange, false, nil)

	return c, q, nil

}
