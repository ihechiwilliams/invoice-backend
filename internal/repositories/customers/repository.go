package customers

import (
	"context"
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

func (S SQLRepository) CreateCustomer(ctx context.Context, customer *Customer) error {
	//TODO implement me
	panic("implement me")
}

func (S SQLRepository) ListCustomers(ctx context.Context, filters *CustomerDBFilter) ([]*Customer, error) {
	//TODO implement me
	panic("implement me")
}

func (S SQLRepository) GetCustomerByID(ctx context.Context, customerID uuid.UUID) (*Customer, error) {
	//TODO implement me
	panic("implement me")
}

func (S SQLRepository) UpdateCustomer(ctx context.Context, customerID uuid.UUID, updatedData *Customer) error {
	//TODO implement me
	panic("implement me")
}

func (S SQLRepository) DeleteCustomer(ctx context.Context, customerID uuid.UUID) error {
	//TODO implement me
	panic("implement me")
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
