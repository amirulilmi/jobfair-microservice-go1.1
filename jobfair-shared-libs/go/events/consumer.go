// File: jobfair-shared-libs/go/events/consumer.go
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// EventHandler is a function that handles an event
type EventHandler func(ctx context.Context, body []byte) error

// NewConsumer creates a new event consumer with retry mechanism
func NewConsumer(rabbitmqURL string) (*Consumer, error) {
	var conn *amqp.Connection
	var channel *amqp.Channel
	var err error

	// Retry connection up to 15 times with exponential backoff
	maxRetries := 15
	for i := 0; i < maxRetries; i++ {
		conn, err = amqp.Dial(rabbitmqURL)
		if err == nil {
			channel, err = conn.Channel()
			if err == nil {
				// Declare exchange
				err = channel.ExchangeDeclare(
					"jobfair.events", // name
					"topic",          // type
					true,             // durable
					false,            // auto-deleted
					false,            // internal
					false,            // no-wait
					nil,              // arguments
				)
				if err == nil {
					log.Println("✅ Connected to RabbitMQ successfully")
					return &Consumer{
						conn:    conn,
						channel: channel,
					}, nil
				}
				channel.Close()
			}
			conn.Close()
		}

		waitTime := time.Duration(i+1) * 2 * time.Second
		log.Printf("⏳ Failed to connect to RabbitMQ (attempt %d/%d). Retrying in %v... Error: %v",
			i+1, maxRetries, waitTime, err)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
}

// Subscribe subscribes to events and processes them with the provided handler
func (c *Consumer) Subscribe(queueName string, routingKeys []string, handler EventHandler) error {
	// Declare queue
	queue, err := c.channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		return fmt.Errorf("failed to declare queue: %w", err)
	}

	// Bind queue to exchange for each routing key
	for _, routingKey := range routingKeys {
		err = c.channel.QueueBind(
			queue.Name,       // queue name
			routingKey,       // routing key
			"jobfair.events", // exchange
			false,
			nil,
		)
		if err != nil {
			return fmt.Errorf("failed to bind queue: %w", err)
		}
		log.Printf("📥 Subscribed to: %s", routingKey)
	}

	// Set QoS - process one message at a time
	err = c.channel.Qos(
		1,     // prefetch count
		0,     // prefetch size
		false, // global
	)
	if err != nil {
		return fmt.Errorf("failed to set QoS: %w", err)
	}

	// Start consuming
	messages, err := c.channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return fmt.Errorf("failed to register consumer: %w", err)
	}

	// Process messages
	go func() {
		ctx := context.Background()
		for msg := range messages {
			log.Printf("📨 Received event: %s", msg.RoutingKey)

			// Process event
			err := handler(ctx, msg.Body)
			if err != nil {
				log.Printf("❌ Error processing event: %v", err)
				// Reject and requeue the message
				msg.Nack(false, true)
			} else {
				log.Printf("✅ Event processed successfully: %s", msg.RoutingKey)
				// Acknowledge the message
				msg.Ack(false)
			}
		}
	}()

	log.Printf("🚀 Consumer started for queue: %s", queueName)
	return nil
}

// SubscribeCompanyEvents is a helper method to subscribe to company events
func (c *Consumer) SubscribeCompanyEvents(handler EventHandler) error {
	return c.Subscribe(
		"company-service.company-events", // queue name
		[]string{
			EventTypeCompanyRegistered,
			EventTypeCompanyUpdated,
			EventTypeCompanyDeleted,
		},
		handler,
	)
}

// Close closes the consumer connection
func (c *Consumer) Close() error {
	if c.channel != nil {
		c.channel.Close()
	}
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// UnmarshalEvent unmarshals an event from JSON
func UnmarshalEvent(body []byte, v interface{}) error {
	return json.Unmarshal(body, v)
}
