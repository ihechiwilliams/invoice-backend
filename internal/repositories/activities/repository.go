package activities

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// CreateActivity Create a new activity log
	CreateActivity(activity *Activity) error

	// ListRecentActivities Retrieve all activities, optionally filtered by type
	ListRecentActivities(limit, offset int) ([]Activity, error)

	// GetActivitiesByInvoiceID Retrieve activities related to a specific invoice
	GetActivitiesByInvoiceID(invoiceID uuid.UUID) ([]Activity, error)

	// SearchActivitiesByType Search activities by type (e.g., Invoice Created, Payment Confirmed)
	SearchActivitiesByType(activityType string, limit, offset int) ([]Activity, error)

	DeleteActivity(id uuid.UUID) error
}

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
