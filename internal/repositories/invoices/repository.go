package invoices

import (
	"context"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	CreateInvoice(ctx context.Context, invoice *Invoice) error
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *Invoice) error
	DeleteInvoice(ctx context.Context, id uuid.UUID) error
	ListInvoices(ctx context.Context, customerID uuid.UUID, limit, offset int) ([]Invoice, error)
	GetInvoicesByStatus(ctx context.Context, status string, limit, offset int) ([]Invoice, error)
	GetTotalInvoiceAmount(ctx context.Context, customerID uuid.UUID) (float64, error)
	ListOverdueInvoices(ctx context.Context, limit, offset int) ([]Invoice, error)
}

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
