package activities

import (
	"context"
	"errors"
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

func (s SQLRepository) CreateActivity(ctx context.Context, activity *Activity) error {
	if activity.ID == uuid.Nil {
		activity.ID = uuid.New() // Generate a new UUID if not provided
	}
	return s.db.WithContext(ctx).Create(activity).Error
}

func (s SQLRepository) ListRecentActivities(ctx context.Context, limit, offset int) ([]Activity, error) {
	var activities []Activity
	err := s.db.WithContext(ctx).
		Order("created_at DESC"). // Order by most recent
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (s SQLRepository) GetActivitiesByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]Activity, error) {
	var activities []Activity
	err := s.db.WithContext(ctx).
		Where("invoice_id = ?", invoiceID).
		Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (s SQLRepository) SearchActivitiesByType(ctx context.Context, activityType string, limit, offset int) ([]Activity, error) {
	var activities []Activity
	err := s.db.WithContext(ctx).
		Where("description LIKE ?", "%"+activityType+"%").
		Order("created_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&activities).Error
	if err != nil {
		return nil, err
	}
	return activities, nil
}

func (s SQLRepository) DeleteActivity(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&Activity{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no activity found with the given ID")
	}
	return nil
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
