package invoices

import (
	"github.com/google/uuid"
	"invoice-backend/internal/repositories/invoices/enums"
	"invoice-backend/internal/repositories/invoicesitems"
	"time"
)

type DBInvoice struct {
	ID            uuid.UUID                    `json:"id" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	CustomerID    uuid.UUID                    `json:"customer_id" gorm:"not null"`
	InvoiceNumber string                       `json:"invoice_number" gorm:"not null;unique"`
	Status        enums.InvoiceStatus          `json:"status" gorm:"not null"` // e.g., Paid, Overdue, Draft
	TotalAmount   float64                      `json:"total_amount" gorm:"not null"`
	DueDate       time.Time                    `json:"due_date" gorm:"not null"`
	IssueDate     time.Time                    `json:"issue_date" gorm:"not null"`
	Items         []*invoicesitems.InvoiceItem `json:"items" gorm:"foreignKey:InvoiceID"` // One-to-Many relationship
	CreatedAt     time.Time                    `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt     time.Time                    `json:"updated_at" gorm:"autoUpdateTime"`
}

type Invoice struct {
	ID            uuid.UUID                    `json:"id"`
	CustomerID    uuid.UUID                    `json:"customer_id"`
	InvoiceNumber string                       `json:"invoice_number"`
	Status        enums.InvoiceStatus          `json:"status"` // e.g., Paid, Overdue, Draft
	TotalAmount   float64                      `json:"total_amount"`
	DueDate       time.Time                    `json:"due_date"`
	IssueDate     time.Time                    `json:"issue_date"`
	Items         []*invoicesitems.InvoiceItem `json:"items"` //One-to-Many relationship
	CreatedAt     time.Time                    `json:"created_at"`
	UpdatedAt     time.Time                    `json:"updated_at"`
}

type InvoiceDBFilter struct {
	CustomerID    []*uuid.UUID           `json:"customer_id,omitempty"`
	UserID        []*uuid.UUID           `json:"user_id,omitempty"`
	ID            []*uuid.UUID           `json:"id,omitempty"`
	InvoiceNumber []*string              `json:"invoice_number,omitempty"`
	Status        []*enums.InvoiceStatus `json:"status,omitempty"`
}

type FindAllInvoicesResult struct {
	Accounts   []*Invoice `json:"accounts"`
	Page       int64      `json:"page"`
	PageSize   int64      `json:"page_size"`
	PageCount  int64      `json:"page_count"`
	TotalCount int64      `json:"total_count"`
}
