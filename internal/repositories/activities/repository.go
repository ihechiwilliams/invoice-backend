package activities

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// CreateActivity Create a new activity log
	CreateActivity(ctx context.Context, activity *Activity) error

	// ListRecentActivities Retrieve all activities, optionally filtered by type
	ListRecentActivities(ctx context.Context, limit, offset int) ([]Activity, error)

	// GetActivitiesByInvoiceID Retrieve activities related to a specific invoice
	GetActivitiesByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]Activity, error)

	// SearchActivitiesByType Search activities by type (e.g., Invoice Created, Payment Confirmed)
	SearchActivitiesByType(ctx context.Context, activityType string, limit, offset int) ([]Activity, error)

	DeleteActivity(ctx context.Context, id uuid.UUID) error
}

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
