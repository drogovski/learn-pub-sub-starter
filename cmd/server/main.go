package main

import (
	"fmt"
	"log"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	CONNECTION_STRING := "amqp://guest:guest@localhost:5672/"

	fmt.Println("Starting Peril server...")
	connection, err := amqp.Dial(CONNECTION_STRING)
	if err != nil {
		log.Fatal(err)
	}
	defer connection.Close()
	fmt.Println("Connection was successful.")
	gamelogic.PrintServerHelp()
	pauseGameChannel, err := connection.Channel()
	if err != nil {
		log.Fatal(err)
	}

	for {
		words := gamelogic.GetInput()

		if len(words) > 0 {
			if words[0] == "pause" {
				log.Println("Sending a pause message.")
				pubsub.SendChangeGameStatusMessage(pauseGameChannel, true)
				continue
			}
			if words[0] == "resume" {
				log.Println("Sending a resume message.")
				pubsub.SendChangeGameStatusMessage(pauseGameChannel, false)
				continue
			}
			if words[0] == "quit" {
				log.Println("Exiting the program.")
				break
			}
			log.Println("Wrong command. Please try again.")
		}
	}

	// wait for ctrl+c
	// signalChan := make(chan os.Signal, 1)
	// signal.Notify(signalChan, os.Interrupt)
	// <-signalChan
	// println("\nSignal received. Terminating the program.")
}
