package invoicesitems

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	tableName = "invoices_items"
)

type Repository interface {
	CreateInvoiceItem(context context.Context, item *InvoiceItem) error
	GetInvoiceItemsByInvoiceID(context context.Context, invoiceID uuid.UUID) ([]InvoiceItem, error)
	UpdateInvoiceItem(context context.Context, item *InvoiceItem) error
	DeleteInvoiceItem(context context.Context, id uuid.UUID) error
}

type SQLRepository struct {
	db *gorm.DB
}

func (s *SQLRepository) CreateInvoiceItem(ctx context.Context, item *InvoiceItem) error {
	return s.db.WithContext(ctx).Create(item).Error
}

func (s *SQLRepository) GetInvoiceItemsByInvoiceID(ctx context.Context, invoiceID uuid.UUID) ([]InvoiceItem, error) {
	var items []InvoiceItem
	err := s.db.WithContext(ctx).
		Where("invoice_id = ?", invoiceID).
		Find(&items).Error
	return items, err
}

func (s *SQLRepository) UpdateInvoiceItem(ctx context.Context, item *InvoiceItem) error {
	return s.db.WithContext(ctx).Save(item).Error
}

func (s *SQLRepository) DeleteInvoiceItem(ctx context.Context, id uuid.UUID) error {
	return s.db.WithContext(ctx).
		Where("id = ?", id).
		Delete(&InvoiceItem{}).Error
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
