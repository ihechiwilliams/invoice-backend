package invoices

import (
	"github.com/google/uuid"
	"invoice-backend/internal/repositories/customers"
	"invoice-backend/internal/repositories/invoicesitems"
	"time"
)

type Invoice struct {
	ID            uuid.UUID                   `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID    uuid.UUID                   `gorm:"not null"`
	Customer      customers.Customer          `gorm:"foreignKey:CustomerID"`
	InvoiceNumber string                      `gorm:"not null;unique"`
	Status        string                      `gorm:"not null"` // e.g., Paid, Overdue, Draft
	TotalAmount   float64                     `gorm:"not null"`
	DueDate       time.Time                   `gorm:"not null"`
	IssueDate     time.Time                   `gorm:"not null"`
	Items         []invoicesitems.InvoiceItem `gorm:"foreignKey:InvoiceID"` // One-to-Many relationship
	CreatedAt     time.Time                   `gorm:"autoCreateTime"`
	UpdatedAt     time.Time                   `gorm:"autoUpdateTime"`
}
