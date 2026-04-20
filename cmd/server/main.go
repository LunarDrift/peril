package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

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

	// wait for ctrl+c
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)
	<-signalChan

	fmt.Println("\nRabbitMQ connection closed...")
}
