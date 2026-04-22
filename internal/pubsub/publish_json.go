// Package pubsub - This is where reusable code for interacting with RabbitMQ will go
// so we can use it in both the server and client
package pubsub

import (
	"context"
	"encoding/json"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"
)

func PublishJSON[T any](ch *amqp.Channel, exchange, key string, val T) error {
	// marshal the val to JSON bytes
	msg, err := json.Marshal(val)
	if err != nil {
		return fmt.Errorf("error marshaling value: %w", err)
	}

	// use the channel's .PublishWithContext to publish message to the
	// exchange with the routing key
	err = ch.PublishWithContext(context.Background(), exchange, key, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        msg,
	})
	if err != nil {
		return fmt.Errorf("error publishing message to the AMQP channel: %w", err)
	}
	return nil
}
