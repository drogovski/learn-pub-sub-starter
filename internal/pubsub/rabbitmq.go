package pubsub

import (
	"context"
	"encoding/json"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
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

	return ch.PublishWithContext(
		context.Background(),
		exchange,
		key,
		mandatory,
		immediate,
		message,
	)
}

func SendChangeGameStatusMessage(pauseGameChannel *amqp.Channel, isPaused bool) {
	gameState := routing.PlayingState{
		IsPaused: isPaused,
	}
	PublishJSON(pauseGameChannel, routing.ExchangePerilDirect, routing.PauseKey, gameState)
}
