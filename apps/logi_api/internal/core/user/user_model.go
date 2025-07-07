package user

import "time"

// User defines the structure for a user in the system.
type User struct {
	UserID       string    `json:"user_id" gorm:"primaryKey;type:varchar(36)"`
	Email        string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"` // Store hashed password
	Role         string    `json:"role" gorm:"type:varchar(50);not null;default:'client'"`
	Name         string    `json:"name" gorm:"type:varchar(100);not null"`
	Phone        string    `json:"phone" gorm:"type:varchar(15);not null"`
	CreatedAt    time.Time `json:"created_at" gorm:"autoCreateTime;not null"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"autoUpdateTime;not null"`
}
