package customers

import (
	"github.com/google/uuid"
	"time"

	"gorm.io/gorm"
)

type DBCustomer struct {
	ID        uuid.UUID      `json:"ID" gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID      `json:"user_id" gorm:"not null"`
	Name      string         `gorm:"type:varchar(255);not null" json:"name"`
	Email     string         `gorm:"type:varchar(255);unique;not null" json:"email"`
	Phone     string         `gorm:"type:varchar(20);unique" json:"phone"`
	Address   string         `gorm:"type:text" json:"address"`
	CreatedAt time.Time      `json:"created_at" gorm:"autoCreateTime"`
	UpdatedAt time.Time      `json:"updated_at" gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"deleted_at"`
}

type Customer struct {
	ID        uuid.UUID `json:"ID"`
	UserID    uuid.UUID `json:"user_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Phone     string    `json:"phone"`
	Address   string    `json:"address"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CustomerDBFilter struct {
	UserID []string `json:"user_id"`
}
