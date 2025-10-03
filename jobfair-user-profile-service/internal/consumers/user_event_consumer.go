// File: internal/consumers/user_event_consumer.go
package consumers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	"jobfair-user-profile-service/internal/models"
	"jobfair-user-profile-service/internal/services"

	"github.com/jobfair/shared/events"
)

type UserEventConsumer struct {
	profileService services.ProfileService
	consumer       *events.Consumer
}

func NewUserEventConsumer(
	rabbitmqURL string,
	profileService services.ProfileService,
) (*UserEventConsumer, error) {
	consumer, err := events.NewConsumer(rabbitmqURL)
	if err != nil {
		return nil, fmt.Errorf("failed to create consumer: %w", err)
	}

	return &UserEventConsumer{
		profileService: profileService,
		consumer:       consumer,
	}, nil
}

// Start begins consuming user events
func (c *UserEventConsumer) Start() error {
	log.Println("🚀 Starting user event consumer...")

	return c.consumer.Subscribe(
		"profile-service.user-events", // queue name
		[]string{events.EventTypeUserRegistered},
		c.handleEvent,
	)
}

// handleEvent processes incoming events
func (c *UserEventConsumer) handleEvent(ctx context.Context, body []byte) error {
	// Parse base event to determine type
	var baseEvent events.BaseEvent
	if err := json.Unmarshal(body, &baseEvent); err != nil {
		return fmt.Errorf("failed to unmarshal base event: %w", err)
	}

	log.Printf("📨 Processing event: %s", baseEvent.EventType)

	// Route to appropriate handler based on event type
	switch baseEvent.EventType {
	case events.EventTypeUserRegistered:
		return c.handleUserRegistered(ctx, body)
	default:
		log.Printf("⚠️ Unknown event type: %s", baseEvent.EventType)
		return nil // Don't fail on unknown events
	}
}

// handleUserRegistered creates a new profile record or updates photo
func (c *UserEventConsumer) handleUserRegistered(ctx context.Context, body []byte) error {
	var event events.UserRegisteredEvent
	if err := json.Unmarshal(body, &event); err != nil {
		return fmt.Errorf("failed to unmarshal user registered event: %w", err)
	}

	data := event.Data
	log.Printf("👤 Processing user registered event for user_id: %d, name: %s", data.UserID, data.FullName)

	// Check if profile already exists (idempotency)
	existingProfile, _ := c.profileService.GetProfile(data.UserID)
	if existingProfile != nil {
		log.Printf("ℹ️ Profile already exists for user_id: %d", data.UserID)
		
		// Update profile photo if provided in event
		if data.ProfilePhotoURL != "" && existingProfile.ProfilePictureURL != data.ProfilePhotoURL {
			log.Printf("📸 Updating profile photo for user_id: %d", data.UserID)
			photoURL := data.ProfilePhotoURL
			req := &models.ProfileUpdateRequest{
				ProfilePictureURL: &photoURL,
			}
			if _, err := c.profileService.UpdateProfile(data.UserID, req); err != nil {
				log.Printf("⚠️ Warning: Failed to update profile photo: %v", err)
			} else {
				log.Printf("✅ Profile photo updated for user_id: %d", data.UserID)
			}
		}
		return nil // Not an error, just idempotency
	}

	// Create profile record with photo
	profile, err := c.profileService.CreateProfile(data.UserID, data.FullName, data.PhoneNumber)
	if err != nil {
		return fmt.Errorf("failed to create profile: %w", err)
	}

	// Update profile photo if provided
	if data.ProfilePhotoURL != "" {
		log.Printf("📸 Setting profile photo for newly created profile: user_id=%d", data.UserID)
		photoURL := data.ProfilePhotoURL
		req := &models.ProfileUpdateRequest{
			ProfilePictureURL: &photoURL,
		}
		if _, err := c.profileService.UpdateProfile(data.UserID, req); err != nil {
			log.Printf("⚠️ Warning: Failed to set profile photo: %v", err)
		}
	}

	log.Printf("✅ Profile created successfully: ID=%d, UserID=%d, Name=%s",
		profile.ID, profile.UserID, profile.FullName)

	return nil
}

// Close closes the consumer
func (c *UserEventConsumer) Close() error {
	if c.consumer != nil {
		return c.consumer.Close()
	}
	return nil
}
