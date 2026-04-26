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
	fmt.Println("Starting Peril client...")
	const connectionStr = "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionStr)
	if err != nil {
		log.Fatal("Could not make connection to RabbitMQ:", err)
	}
	defer connection.Close()
	fmt.Println("Peril game client successfully connected to RabbitMQ!")

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("Could not get username: %v", err)
	}

	gs := gamelogic.NewGameState(username)

	err = pubsub.SubscribeJSON(
		connection,
		routing.ExchangePerilDirect,
		routing.PauseKey+"."+gs.GetUsername(),
		routing.PauseKey,
		pubsub.SimpleQueueTransient,
		handlerPause(gs),
	)
	if err != nil {
		log.Fatalf("could not subscribe to pause: %v", err)
	}

	for {
		words := gamelogic.GetInput()
		if len(words) == 0 {
			continue
		}
		switch words[0] {
		case "spawn":
			err := gs.CommandSpawn(words)
			if err != nil {
				log.Printf("could not spawn unit: %v", err)
			}

		case "move":
			move, err := gs.CommandMove(words)
			if err != nil {
				log.Printf("could not move unit: %v", err)
			} else {
				fmt.Println("Move successful:", move)
			}

		case "status":
			gs.CommandStatus()

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
