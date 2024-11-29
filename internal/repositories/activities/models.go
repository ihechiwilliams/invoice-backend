package activities

import (
	"github.com/google/uuid"
	"invoice-backend/internal/repositories/invoices"
	"time"
)

type Activity struct {
	ID          uuid.UUID        `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Description string           `gorm:"not null"`
	InvoiceID   uuid.UUID        `gorm:"not null"`
	Invoice     invoices.Invoice `gorm:"foreignKey:InvoiceID"`
	CreatedAt   time.Time        `gorm:"autoCreateTime"`
}
