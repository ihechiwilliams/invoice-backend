package invoicesitems

import "github.com/google/uuid"

type InvoiceItem struct {
	ID          uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	InvoiceID   uuid.UUID `json:"invoice_id" gorm:"type:uuid;not null"`
	Description string    `gorm:"not null"` // Item description
	Quantity    int       `gorm:"not null"` // Number of items
	UnitPrice   float64   `gorm:"not null"` // Price per item
	TotalPrice  float64   `gorm:"not null"` // Quantity * UnitPrice
}
