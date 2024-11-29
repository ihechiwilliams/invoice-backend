package invoices

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"
	"invoice-backend/internal/shared"
	"time"
)

const (
	tableName = "invoices"
)

type Repository interface {
	CreateInvoice(ctx context.Context, invoice *DBInvoice) (*Invoice, error)
	GetInvoiceByID(ctx context.Context, id uuid.UUID) (*Invoice, error)
	UpdateInvoice(ctx context.Context, invoice *Invoice) error
	DeleteInvoice(ctx context.Context, id uuid.UUID) error
	ListInvoices(ctx context.Context, filters *InvoiceDBFilter, pagination shared.Pagination) ([]*Invoice, error)
	GetTotalInvoiceAmount(ctx context.Context, customerID uuid.UUID) (float64, error)
	ListOverdueInvoices(ctx context.Context, limit, offset int) ([]Invoice, error)
}

type SQLRepository struct {
	db *gorm.DB
}

func (s *SQLRepository) CreateInvoice(ctx context.Context, invoice *DBInvoice) (*Invoice, error) {
	if invoice.ID == uuid.Nil {
		invoice.ID = uuid.New()
	}
	result := s.db.WithContext(ctx).Table(tableName).Create(invoice)

	if result.Error != nil {
		return nil, result.Error
	}
	return FromDBInvoice(invoice), nil
}

func (s *SQLRepository) GetInvoiceByID(ctx context.Context, id uuid.UUID) (*Invoice, error) {
	var invoice DBInvoice

	err := s.db.Clauses(dbresolver.Write).WithContext(ctx).Table(tableName).Where("id = ?", id).First(&invoice).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}

		return nil, err
	}

	return FromDBInvoice(&invoice), err
}

func (s *SQLRepository) UpdateInvoice(ctx context.Context, invoice *Invoice) error {
	return s.db.WithContext(ctx).Save(invoice).Error
}

func (s *SQLRepository) DeleteInvoice(ctx context.Context, id uuid.UUID) error {
	result := s.db.WithContext(ctx).Where("id = ?", id).Delete(&Invoice{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("no invoice found with the given ID")
	}
	return nil
}

func (s *SQLRepository) ListInvoices(ctx context.Context, filters *InvoiceDBFilter, pagination shared.Pagination) ([]*Invoice, error) {
	invoices := make([]*DBInvoice, 0)

	dataset, err := shared.BuildDataset(ctx, s.db, tableName, filters)
	if err != nil {
		return nil, err
	}

	paginatedDataset := shared.PaginateDataset(dataset, pagination)

	result := paginatedDataset.Find(&invoices)
	if result.Error != nil {
		return nil, result.Error
	}

	return FromDBInvoiceList(invoices), nil
}

func (s *SQLRepository) GetTotalInvoiceAmount(ctx context.Context, customerID uuid.UUID) (float64, error) {
	var total float64
	err := s.db.WithContext(ctx).
		Model(&Invoice{}).
		Where("customer_id = ?", customerID).
		Select("SUM(total_amount)").
		Scan(&total).Error
	return total, err
}

func (s *SQLRepository) ListOverdueInvoices(ctx context.Context, limit, offset int) ([]Invoice, error) {
	var invoices []Invoice
	now := time.Now()
	err := s.db.WithContext(ctx).
		Where("due_date < ? AND status != ?", now, "Paid").
		Order("due_date ASC").
		Limit(limit).
		Offset(offset).
		Find(&invoices).Error
	return invoices, err
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
