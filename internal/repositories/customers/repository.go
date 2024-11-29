package customers

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// CreateCustomer Create a new customer
	CreateCustomer(ctx context.Context, customer *Customer) error

	// ListCustomers Retrieve all customers
	ListCustomers(ctx context.Context, filters *CustomerDBFilter) ([]*Customer, error)

	// GetCustomerByID Retrieve a customer by ID
	GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*Customer, error)

	// UpdateCustomer Update customer details (optional based on use case)
	UpdateCustomer(ctx context.Context, customerID uuid.UUID, updatedData *Customer) error

	// DeleteCustomer Delete a customer (optional, soft-delete recommended)
	DeleteCustomer(ctx context.Context, customerID uuid.UUID) error
}

type SQLRepository struct {
	db *gorm.DB
}

func (s SQLRepository) CreateCustomer(ctx context.Context, customer *Customer) error {
	if customer.ID == uuid.Nil {
		customer.ID = uuid.New()
	}

	return s.db.WithContext(ctx).Create(customer).Error
}

func (s SQLRepository) ListCustomers(ctx context.Context, filters *CustomerDBFilter) ([]*Customer, error) {
	var customers []*Customer
	query := s.db.WithContext(ctx)

	if filters != nil && len(filters.UserID) > 0 {
		query = query.Where("id IN ?", filters.UserID)
	}

	if err := query.Find(&customers).Error; err != nil {
		return nil, err
	}

	return customers, nil
}

func (s SQLRepository) GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*Customer, error) {
	var customer Customer
	err := s.db.WithContext(ctx).First(&customer, "id = ?", customerID).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	} else if err != nil {
		return nil, err
	}

	return &customer, nil
}

func (s SQLRepository) UpdateCustomer(ctx context.Context, customerID uuid.UUID, updatedData *Customer) error {
	return s.db.WithContext(ctx).
		Model(&Customer{}).
		Where("id = ?", customerID).
		Updates(updatedData).
		Error
}

func (s SQLRepository) DeleteCustomer(ctx context.Context, customerID uuid.UUID) error {
	return s.db.WithContext(ctx).
		Where("id = ?", customerID).
		Delete(&Customer{}).
		Error
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
