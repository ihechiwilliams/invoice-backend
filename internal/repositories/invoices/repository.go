package invoices

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateInvoice(invoice *Invoice) error
	GetInvoiceByID(id uuid.UUID) (*Invoice, error)
	UpdateInvoice(invoice *Invoice) error
	DeleteInvoice(id uuid.UUID) error
	ListInvoices(customerID uuid.UUID, limit, offset int) ([]Invoice, error)
	GetInvoicesByStatus(status string, limit, offset int) ([]Invoice, error)
	GetTotalInvoiceAmount(customerID uuid.UUID) (float64, error)
	ListOverdueInvoices(limit, offset int) ([]Invoice, error)
}

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
