// File: internal/consumers/company_event_consumer.go
package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"jobfair-company-service/internal/models"
	"jobfair-company-service/internal/repository"
	"jobfair-company-service/internal/utils"

	"github.com/jobfair/shared/events"
)

type CompanyEventConsumer struct {
	companyRepo    *repository.CompanyRepository
	consumer       *events.Consumer
	eventPublisher *events.Publisher
}

func NewCompanyEventConsumer(
	rabbitmqURL string,
	companyRepo *repository.CompanyRepository,
) (*CompanyEventConsumer, error) {
	consumer, err := events.NewConsumer(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	// Initialize publisher to re-publish events with company_id
	publisher, err := events.NewPublisher(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create publisher: %w", err)
	}

	return &CompanyEventConsumer{
		companyRepo:    companyRepo,
		consumer:       consumer,
		eventPublisher: publisher,
	}, nil
}

// Start begins consuming company events
func (c *CompanyEventConsumer) Start() error {
	log.Println("🚀 Starting company event consumer...")

	return c.consumer.SubscribeCompanyEvents(c.handleEvent)
}

// handleEvent processes incoming events
func (c *CompanyEventConsumer) handleEvent(ctx context.Context, body []byte) error {
	// Parse base event to determine type
	var baseEvent events.BaseEvent
	if err := json.Unmarshal(body, &baseEvent); err != nil {
		return fmt.Errorf("failed to unmarshal base event: %w", err)
	}

	log.Printf("📨 Processing event: %s", baseEvent.EventType)

	// Route to appropriate handler based on event type
	switch baseEvent.EventType {
	case events.EventTypeCompanyRegistered:
		return c.handleCompanyRegistered(ctx, body)
	case events.EventTypeCompanyUpdated:
		return c.handleCompanyUpdated(ctx, body)
	case events.EventTypeCompanyDeleted:
		return c.handleCompanyDeleted(ctx, body)
	default:
		log.Printf("⚠️ Unknown event type: %s", baseEvent.EventType)
		return nil // Don't fail on unknown events
	}
}

// handleCompanyRegistered creates a new company record
func (c *CompanyEventConsumer) handleCompanyRegistered(ctx context.Context, body []byte) error {
	var event events.CompanyRegisteredEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal company registered event: %w", err)
	}

	data := event.Data
	log.Printf("🏢 Creating company for user_id: %d, name: %s", data.UserID, data.CompanyName)

	// Check if company already exists
	existingCompany, _ := c.companyRepo.GetByUserID(data.UserID)
	if existingCompany != nil {
		log.Printf("ℹ️ Company already exists for user_id: %d, skipping", data.UserID)
		return nil // Not an error, just idempotency
	}

	// Create slug from company name
	slug := utils.GenerateSlug(data.CompanyName)

	// Create company record
	company := &models.Company{
		UserID:      data.UserID,
		Name:        data.CompanyName,
		Email:       data.Email,
		Phone:       data.Phone,
		Website:     data.Website,
		Industry:    data.Industry, // ✅ Now using array directly
		Address:     data.Address,
		LogoURL:     data.LogoURL,
		ContactName: data.ContactName, // ✅ Added ContactName
		CompanySize: "1-10",               // ✅ Default company size
		Slug:        slug,
		IsVerified:  false,
		IsFeatured:  false,
		IsPremium:   false,
	}

	// Save to database
	createdCompany, err := c.companyRepo.Create(company)
	if err != nil {
		return fmt.Errorf("failed to create company: %w", err)
	}

	log.Printf("✅ Company created successfully: ID=%d, UserID=%d, Name=%s",
		createdCompany.ID, createdCompany.UserID, createdCompany.Name)

	// 🚀 Re-publish event WITH company_id for other services (job-service)
	if c.eventPublisher != nil {
		eventDataWithID := events.CompanyRegisteredData{
			UserID:      data.UserID,
			CompanyID:   createdCompany.ID, // ✅ Now includes company_id
			CompanyName: data.CompanyName,
			Email:       data.Email,
			Phone:       data.Phone,
			Website:     data.Website,
			Industry:    data.Industry,
			Address:     data.Address,
			LogoURL:     data.LogoURL,
			CountryCode: data.CountryCode,
			ContactName: data.ContactName,
		}

		if err := c.eventPublisher.PublishCompanyRegistered(ctx, eventDataWithID); err != nil {
			log.Printf("⚠️ Warning: Failed to re-publish company registered event with ID: %v", err)
			// Don't fail the whole operation
		} else {
			log.Printf("📤 Re-published company.registered event with company_id=%d", createdCompany.ID)
		}
	}

	return nil
}

// handleCompanyUpdated updates company record
func (c *CompanyEventConsumer) handleCompanyUpdated(ctx context.Context, body []byte) error {
	var event events.CompanyUpdatedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal company updated event: %w", err)
	}

	data := event.Data
	log.Printf("🔄 Updating company for user_id: %d", data.UserID)

	company, err := c.companyRepo.GetByUserID(data.UserID)
	if err != nil {
		return fmt.Errorf("company not found for user_id %d: %w", data.UserID, err)
	}

	// Update fields from event
	// TODO: Implement field updates based on data.UpdatedFields

	if err := c.companyRepo.Update(company); err != nil {
		return fmt.Errorf("failed to update company: %w", err)
	}

	log.Printf("✅ Company updated successfully: ID=%d, UserID=%d", company.ID, company.UserID)
	return nil
}

// handleCompanyDeleted soft deletes company record
func (c *CompanyEventConsumer) handleCompanyDeleted(ctx context.Context, body []byte) error {
	var event events.CompanyDeletedEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal company deleted event: %w", err)
	}

	data := event.Data
	log.Printf("🗑️ Deleting company for user_id: %d", data.UserID)

	company, err := c.companyRepo.GetByUserID(data.UserID)
	if err != nil {
		return fmt.Errorf("company not found for user_id %d: %w", data.UserID, err)
	}

	if err := c.companyRepo.Delete(company.ID); err != nil {
		return fmt.Errorf("failed to delete company: %w", err)
	}

	log.Printf("✅ Company deleted successfully: ID=%d, UserID=%d", company.ID, company.UserID)
	return nil
}

// Close closes the consumer and publisher
func (c *CompanyEventConsumer) Close() error {
	var err error
	if c.consumer != nil {
		err = c.consumer.Close()
	}
	if c.eventPublisher != nil {
		if pubErr := c.eventPublisher.Close(); pubErr != nil && err == nil {
			err = pubErr
		}
	}
	return err
}
