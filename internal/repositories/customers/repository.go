package customers

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Repository interface {
	// CreateCustomer Create a new customer
	CreateCustomer(customer *Customer) error

	// ListCustomers Retrieve all customers
	ListCustomers(limit, offset int) ([]Customer, error)

	// GetCustomerByID Retrieve a customer by ID
	GetCustomerByID(customerID uuid.UUID) (*Customer, error)

	// UpdateCustomer Update customer details (optional based on use case)
	UpdateCustomer(customerID uuid.UUID, updatedData *Customer) error

	// DeleteCustomer Delete a customer (optional, soft-delete recommended)
	DeleteCustomer(customerID uuid.UUID) error
}

type SQLRepository struct {
	db *gorm.DB
}

func NewSQLRepository(db *gorm.DB) *SQLRepository {
	return &SQLRepository{
		db: db,
	}
}
