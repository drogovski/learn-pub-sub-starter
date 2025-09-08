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
	routingKey := fmt.Sprintf("%s.*", routing.GameLogSlug)
	_, queue, err := pubsub.DeclareAndBind(
		connection,
		routing.ExchangePerilTopic,
		routing.GameLogSlug,
		routingKey,
		pubsub.Durable,
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}
	fmt.Printf("Queue %v declared and bound!\n", queue.Name)

	for {
		words := gamelogic.GetInput()

		if len(words) > 0 {
			if words[0] == "pause" {
				log.Println("Sending a pause message.")
				err = pubsub.SendChangeGameStatusMessage(pauseGameChannel, true)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
				continue
			}
			if words[0] == "resume" {
				log.Println("Sending a resume message.")
				err = pubsub.SendChangeGameStatusMessage(pauseGameChannel, false)
				if err != nil {
					log.Printf("could not publish time: %v", err)
				}
				continue
			}
			if words[0] == "quit" {
				log.Println("Exiting the program.")
				break
			}
			log.Println("Wrong command. Please try again.")
		}
	}
}
