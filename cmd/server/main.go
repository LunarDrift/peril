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
	fmt.Println("Starting Peril server...")

	const connectionStr = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatal("Could not make connection to RabbitMQ:", err)
	}
	defer connection.Close()
	fmt.Println("Peril game server successfully connected to RabbitMQ!")

	publishChannel, err := connection.Channel()
	if err != nil {
		log.Fatal("Could not create AMQP Channel:", err)
	}

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "pause":
			fmt.Println("Publishing paused game state")
			err = pubsub.PublishJSON(publishChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}

		case "resume":
			fmt.Println("Publishing resumes game state")
			err = pubsub.PublishJSON(publishChannel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Printf("could not publish time: %v", err)
			}

		case "quit":
			log.Println("goodbye")
			return

		default:
			fmt.Println("unknown command")
		}
	}
}
