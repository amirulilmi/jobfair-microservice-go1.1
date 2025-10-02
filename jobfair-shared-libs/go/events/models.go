// File: jobfair-shared-libs/go/events/models.go
package events

import "time"

// Event Types
const (
	EventTypeCompanyRegistered = "company.registered"
	EventTypeCompanyUpdated    = "company.updated"
	EventTypeCompanyDeleted    = "company.deleted"
)

// BaseEvent contains common fields for all events
type BaseEvent struct {
	EventID   string    `json:"event_id"`
	EventType string    `json:"event_type"`
	Timestamp time.Time `json:"timestamp"`
	Version   string    `json:"version"`
}

// CompanyRegisteredEvent is published when a company completes registration
type CompanyRegisteredEvent struct {
	BaseEvent
	Data CompanyRegisteredData `json:"data"`
}

type CompanyRegisteredData struct {
	UserID      uint     `json:"user_id"`
	CompanyName string   `json:"company_name"`
	Email       string   `json:"email"`
	Phone       string   `json:"phone"`
	Website     string   `json:"website"`
	Industry    []string `json:"industry"` // Changed to array for multiple industries
	Address     string   `json:"address"`
	LogoURL     string   `json:"logo_url"`
	CountryCode string   `json:"country_code"`
	ContactName string   `json:"contact_name"` // Added for UI support
}

// CompanyUpdatedEvent is published when company profile is updated
type CompanyUpdatedEvent struct {
	BaseEvent
	Data CompanyUpdatedData `json:"data"`
}

type CompanyUpdatedData struct {
	UserID      uint              `json:"user_id"`
	UpdatedFields map[string]interface{} `json:"updated_fields"`
}

// CompanyDeletedEvent is published when company is deleted
type CompanyDeletedEvent struct {
	BaseEvent
	Data CompanyDeletedData `json:"data"`
}

type CompanyDeletedData struct {
	UserID uint `json:"user_id"`
}
