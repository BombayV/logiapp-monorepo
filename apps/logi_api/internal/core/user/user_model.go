package user

import "time"

// User defines the structure for a user in the system.
type User struct {
	UserID       string    `json:"user_id" gorm:"primaryKey;type:varchar(36)"`
	Email        string    `json:"email" gorm:"type:varchar(100);unique;not null"`
	PasswordHash string    `json:"-" gorm:"type:varchar(255);not null"` // Store hashed password
	Role         string    `json:"role" gorm:"type:user_role;not null;default:'sales'"`
	CreatedAt    time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt    time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}

type UserData struct {
	UserID         string    `json:"user_id" gorm:"primaryKey;type:varchar(36)"`
	FirstName      string    `json:"first_name" gorm:"type:varchar(100);not null"`
	LastName       string    `json:"last_name" gorm:"type:varchar(100);not null"`
	PhoneNumber    string    `json:"phone_number" gorm:"type:varchar(16);unique;not null;"`
	LastConnection time.Time `json:"last_connection" gorm:"type:timestamp with time zone;default:current_timestamp"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}

type UserLocation struct {
	UserID      string    `json:"user_id" gorm:"primaryKey;type:varchar(36)"`
	Location    string    `json:"location" gorm:"type:geography(Point, 4326);not null"`
	LastUpdated time.Time `json:"last_updated" gorm:"type:timestamp with time zone;default:current_timestamp"`
}
