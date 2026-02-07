package entities

import (
	"time"

	"gorm.io/gorm"
)

type UserEntity struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`

	Username     string `gorm:"not null;unique"`
	PasswordHash string `gorm:"not null"`
}
