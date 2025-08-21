package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	CONNECTION_STRING := "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")
	connection, _ := amqp.Dial(CONNECTION_STRING)
	defer connection.Close()
	fmt.Println("Connection was successful.")

	pauseGameChannel, _ := connection.Channel()
	gameState := routing.PlayingState{
		IsPaused: true,
	}
	pubsub.PublishJSON(pauseGameChannel, routing.ExchangePerilDirect, routing.PauseKey, gameState)

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	println("\nSignal received. Terminating the program.")
}
