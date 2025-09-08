package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	const CONNECTION_STRING = "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")

	connection, err := amqp.Dial(CONNECTION_STRING)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer connection.Close()
	fmt.Println("Client connected successfully to RabbitMQ!.")

	publishCh, err := connection.Channel()
	if err != nil {
		log.Fatalf("could not create channel: %v", err)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not get username: %v", err)
	}
	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)
	gameState := gamelogic.NewGameState(username)
	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		queueName,
		routing.PauseKey,
		pubsub.Transient,
		handlerPause(gameState),
	)
	if err != nil {
		log.Fatalf("could not subscribe to GameState queue: %v", err)
	}

	armyMovesKey := fmt.Sprintf("%s.*", routing.ArmyMovesPrefix)
	armyMovesQueueName := fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, gameState.GetUsername())

	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		armyMovesQueueName,
		armyMovesKey,
		pubsub.Transient,
		handlerMove(publishCh, gameState),
	)
	if err != nil {
		log.Fatalf("could not subscribe to army moves queue: %v", err)
	}

	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilTopic,
		routing.WarRecognitionsPrefix,
		routing.WarRecognitionsPrefix+".*",
		pubsub.Durable,
		handlerWar(gameState),
	)
	if err != nil {
		log.Fatalf("could not subscribe to war declarations queue: %v", err)
	}

	gamelogic.PrintClientHelp()
	for {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "move":
			move, err := gameState.CommandMove(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
			err = pubsub.PublishJSON(
				publishCh,
				routing.ExchangePerilTopic,
				fmt.Sprintf("%s.%s", routing.ArmyMovesPrefix, username),
				move,
			)
			if err != nil {
				fmt.Printf("error: %v", err)
			}
			fmt.Printf("Successfully published the move: %v", move)
		case "spawn":
			err := gameState.CommandSpawn(words)
			if err != nil {
				fmt.Println(err)
				continue
			}
		case "status":
			gameState.CommandStatus()
		case "help":
			gamelogic.PrintClientHelp()
		case "spam":
			fmt.Println("Spamming not allowed yet!")
		case "quit":
			gamelogic.PrintQuit()
			return
		default:
			fmt.Println("unknown command")
		}
	}
}
