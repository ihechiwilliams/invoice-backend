package invoicesitems

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

const (
	tableName = "invoices_items"
)

type Repository interface {
	CreateInvoiceItem(item *InvoiceItem) error
	GetInvoiceItemsByInvoiceID(invoiceID uuid.UUID) ([]InvoiceItem, error)
	UpdateInvoiceItem(item *InvoiceItem) error
	DeleteInvoiceItem(id uuid.UUID) error
}

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
