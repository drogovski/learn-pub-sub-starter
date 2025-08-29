package pubsub

import (
	"context"
	"encoding/json"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

type simpleQueueType int

const (
	Transient simpleQueueType = iota
	Durable
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	data, err := json.Marshal(val)
	if err != nil {
		return err
	}

	mandatory := false
	immediate := false
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        data,
	}
	err = ch.PublishWithContext(context.Background(), exchange, key, mandatory, immediate, message)
	if err != nil {
		return err
	}
	return nil
}

func DeclareAndBind(
	conn *amqp.Connection,
	exchange,
	queueName,
	key string,
	queueType simpleQueueType, // an enum to represent "durable" or "transient"
) (*amqp.Channel, amqp.Queue, error) {
	channel, err := conn.Channel()
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	var durable, autoDelete, exclusive, noWait bool

	switch queueType {
	case Durable:
		durable = true
		autoDelete = false
		exclusive = false
		noWait = false
	case Transient:
		durable = false
		autoDelete = true
		exclusive = true
		noWait = false
	}

	queue, err := channel.QueueDeclare(queueName, durable, autoDelete, exclusive, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	err = channel.QueueBind(queue.Name, key, exchange, noWait, nil)
	if err != nil {
		return nil, amqp.Queue{}, err
	}

	return channel, queue, nil
}

func SendChangeGameStatusMessage(pauseGameChannel *amqp.Channel, isPaused bool) {
	gameState := routing.PlayingState{
		IsPaused: isPaused,
	}
	PublishJSON(pauseGameChannel, routing.ExchangePerilDirect, routing.PauseKey, gameState)
}
