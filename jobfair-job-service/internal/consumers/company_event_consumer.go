// File: internal/consumers/company_event_consumer.go
package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"jobfair-job-service/internal/repository"

	"github.com/jobfair/shared/events"
)

type CompanyEventConsumer struct {
	companyRepo *repository.CompanyRepository
	consumer    *events.Consumer
}

func NewCompanyEventConsumer(
	rabbitmqURL string,
	companyRepo *repository.CompanyRepository,
) (*CompanyEventConsumer, error) {
	consumer, err := events.NewConsumer(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &CompanyEventConsumer{
		companyRepo: companyRepo,
		consumer:    consumer,
	}, nil
}

// Start begins consuming company events
func (c *CompanyEventConsumer) Start() error {
	log.Println("🚀 [JOB-SERVICE] Starting company event consumer...")

	return c.consumer.SubscribeCompanyEvents(c.handleEvent)
}

// handleEvent processes incoming events
func (c *CompanyEventConsumer) handleEvent(ctx context.Context, body []byte) error {
	// Parse base event to determine type
	var baseEvent events.BaseEvent
	if err := json.Unmarshal(body, &baseEvent); err != nil {
		return fmt.Errorf("failed to unmarshal base event: %w", err)
	}

	log.Printf("📨 [JOB-SERVICE] Processing event: %s", baseEvent.EventType)

	// Route to appropriate handler based on event type
	switch baseEvent.EventType {
	case events.EventTypeCompanyRegistered:
		return c.handleCompanyRegistered(ctx, body)
	case events.EventTypeCompanyUpdated:
		return c.handleCompanyUpdated(ctx, body)
	case events.EventTypeCompanyDeleted:
		return c.handleCompanyDeleted(ctx, body)
	default:
		log.Printf("⚠️ [JOB-SERVICE] Unknown event type: %s", baseEvent.EventType)
		return nil // Don't fail on unknown events
	}
}

// handleCompanyRegistered creates company mapping
func (c *CompanyEventConsumer) handleCompanyRegistered(ctx context.Context, body []byte) error {
	var event events.CompanyRegisteredEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal company registered event: %w", err)
	}

	data := event.Data
	log.Printf("🏢 [JOB-SERVICE] Creating company mapping for user_id: %d, company: %s", 
		data.UserID, data.CompanyName)

	// Check if company_id is provided (re-published event from company-service)
	if data.CompanyID == 0 {
		log.Printf("⚠️ [JOB-SERVICE] Company event without company_id, skipping mapping (will be created by re-published event)")
		return nil
	}

	// Create or update company mapping
	if err := c.companyRepo.UpsertCompanyMapping(data.UserID, data.CompanyID, data.CompanyName); err != nil {
		return fmt.Errorf("failed to upsert company mapping: %w", err)
	}

	log.Printf("✅ [JOB-SERVICE] Company mapping created: UserID=%d, CompanyID=%d, Name=%s",
		data.UserID, data.CompanyID, data.CompanyName)

	return nil
}

// handleCompanyUpdated updates company mapping
func (c *CompanyEventConsumer) handleCompanyUpdated(ctx context.Context, body []byte) error {
	var event events.CompanyUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal company updated event: %w", err)
	}

	data := event.Data
	log.Printf("🔄 [JOB-SERVICE] Updating company mapping for user_id: %d", data.UserID)

	// Extract updated fields
	var companyID uint
	var companyName string
	
	if id, ok := data.UpdatedFields["company_id"].(float64); ok {
		companyID = uint(id)
	}
	
	if name, ok := data.UpdatedFields["company_name"].(string); ok {
		companyName = name
	}

	// If we have company_id, create or update mapping
	if companyID > 0 {
		if err := c.companyRepo.UpsertCompanyMapping(data.UserID, companyID, companyName); err != nil {
			return fmt.Errorf("failed to upsert company mapping: %w", err)
		}
		
		log.Printf("✅ [JOB-SERVICE] Company mapping updated: UserID=%d, CompanyID=%d", 
			data.UserID, companyID)
	}

	return nil
}

// handleCompanyDeleted removes company mapping
func (c *CompanyEventConsumer) handleCompanyDeleted(ctx context.Context, body []byte) error {
	var event events.CompanyDeletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal company deleted event: %w", err)
	}

	data := event.Data
	log.Printf("🗑️ [JOB-SERVICE] Deleting company mapping for user_id: %d", data.UserID)

	if err := c.companyRepo.DeleteMappingByUserID(data.UserID); err != nil {
		log.Printf("⚠️ [JOB-SERVICE] Failed to delete mapping (may not exist): %v", err)
		// Don't return error as mapping might not exist
	}

	log.Printf("✅ [JOB-SERVICE] Company mapping deleted for UserID=%d", data.UserID)
	return nil
}

// Close closes the consumer
func (c *CompanyEventConsumer) Close() error {
	if c.consumer != nil {
		return c.consumer.Close()
	}
	return nil
}
