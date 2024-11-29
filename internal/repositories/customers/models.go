package customers

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type Customer struct {
	ID        uuid.UUID      `json:"ID" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Phone     string         `gorm:"type:varchar(20);unique" json:"phone"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}
