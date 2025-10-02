// File: jobfair-shared-libs/go/events/publisher.go
package events

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	amqp "github.com/rabbitmq/amqp091-go"
)

type Publisher struct {
	conn    *amqp.Connection
	channel *amqp.Channel
}

// NewPublisher creates a new event publisher with retry mechanism
func NewPublisher(rabbitmqURL string) (*Publisher, error) {
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
				// Declare exchange for events
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
					log.Println("✅ Connected to RabbitMQ (Publisher) successfully")
					return &Publisher{
						conn:    conn,
						channel: channel,
					}, nil
				}
				channel.Close()
			}
			conn.Close()
		}

		waitTime := time.Duration(i+1) * 2 * time.Second
		log.Printf("⏳ Failed to connect to RabbitMQ Publisher (attempt %d/%d). Retrying in %v... Error: %v",
			i+1, maxRetries, waitTime, err)
		time.Sleep(waitTime)
	}

	return nil, fmt.Errorf("failed to connect to RabbitMQ after %d attempts: %w", maxRetries, err)
}

// PublishCompanyRegistered publishes a company registered event
func (p *Publisher) PublishCompanyRegistered(ctx context.Context, data CompanyRegisteredData) error {
	event := CompanyRegisteredEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New().String(),
			EventType: EventTypeCompanyRegistered,
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		Data: data,
	}

	return p.publish(ctx, EventTypeCompanyRegistered, event)
}

// PublishCompanyUpdated publishes a company updated event
func (p *Publisher) PublishCompanyUpdated(ctx context.Context, data CompanyUpdatedData) error {
	event := CompanyUpdatedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New().String(),
			EventType: EventTypeCompanyUpdated,
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		Data: data,
	}

	return p.publish(ctx, EventTypeCompanyUpdated, event)
}

// PublishCompanyDeleted publishes a company deleted event
func (p *Publisher) PublishCompanyDeleted(ctx context.Context, data CompanyDeletedData) error {
	event := CompanyDeletedEvent{
		BaseEvent: BaseEvent{
			EventID:   uuid.New().String(),
			EventType: EventTypeCompanyDeleted,
			Timestamp: time.Now(),
			Version:   "1.0",
		},
		Data: data,
	}

	return p.publish(ctx, EventTypeCompanyDeleted, event)
}

// publish is the internal method to publish events
func (p *Publisher) publish(ctx context.Context, routingKey string, event interface{}) error {
	body, err := json.Marshal(event)
	if err != nil {
		return fmt.Errorf("failed to marshal event: %w", err)
	}

	err = p.channel.PublishWithContext(
		ctx,
		"jobfair.events", // exchange
		routingKey,       // routing key
		false,            // mandatory
		false,            // immediate
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         body,
			DeliveryMode: amqp.Persistent, // persistent messages
			Timestamp:    time.Now(),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to publish event: %w", err)
	}

	log.Printf("✅ Event published: %s", routingKey)
	return nil
}

// Close closes the publisher connection
func (p *Publisher) Close() error {
	if p.channel != nil {
		p.channel.Close()
	}
	if p.conn != nil {
		return p.conn.Close()
	}
	return nil
}
