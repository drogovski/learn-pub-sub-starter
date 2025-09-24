package pubsub

import (
	"bytes"
	"context"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"time"

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

func PublishGob[T any](ch *amqp.Channel, exchange, key string, val T) error {
	var buffer bytes.Buffer
	encoder := gob.NewEncoder(&buffer)
	err := encoder.Encode(val)
	if err != nil {
		return err
	}

	mandatory := false
	immediate := false
	message := amqp.Publishing{
		ContentType: "application/gob",
		Body:        buffer.Bytes(),
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

func SendChangeGameStatusMessage(pauseGameChannel *amqp.Channel, isPaused bool) error {
	gameState := routing.PlayingState{
		IsPaused: isPaused,
	}
	return PublishJSON(pauseGameChannel, routing.ExchangePerilDirect, routing.PauseKey, gameState)
}

func PublishGameLogMessage(ch *amqp.Channel, message, username string) error {
	gameLog := routing.GameLog{
		CurrentTime: time.Now().UTC(),
		Message:     message,
		Username:    username,
	}

	return PublishGob(
		ch,
		routing.ExchangePerilTopic,
		fmt.Sprintf("%s.%s", routing.GameLogSlug, username),
		gameLog,
	)
}
