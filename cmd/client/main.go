package main

import (
	"fmt"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	CONNECTION_STRING := "amqp://guest:guest@localhost:5672/"
	fmt.Println("Starting Peril client...")
	connection, _ := amqp.Dial(CONNECTION_STRING)
	defer connection.Close()
	fmt.Println("Client connected successfully.")
	username, err := gamelogic.ClientWelcome()
	if err != nil {
		fmt.Printf("%v", err)
	}
	queueName := fmt.Sprintf("%s.%s", routing.PauseKey, username)

	pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan
	println("\nSignal received. Terminating the program.")
}
