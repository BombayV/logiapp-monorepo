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
	UserID    string    `json:"user_id" gorm:"primaryKey;type:varchar(36)"`
	Location  string    `json:"location" gorm:"type:geometry(Point, 4326);not null"`
	Latitude  float64   `json:"latitude" gorm:"-"`  // Not stored in DB, computed from Location
	Longitude float64   `json:"longitude" gorm:"-"` // Not stored in DB, computed from Location
	UpdatedAt time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}

// UserWithData combines User and UserData for operations that need both
type UserWithData struct {
	User     *User     `json:"user"`
	UserData *UserData `json:"user_data"`
}

// DriverLocation represents a driver's basic info with location
type DriverLocation struct {
	UserID            string    `json:"user_id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"first_name"`
	LastName          string    `json:"last_name"`
	PhoneNumber       string    `json:"phone_number"`
	Latitude          float64   `json:"latitude"`
	Longitude         float64   `json:"longitude"`
	LastConnection    time.Time `json:"last_connection"`
	LocationUpdatedAt time.Time `json:"location_updated_at"`
}

// Driver represents a driver's basic info without location
type Driver struct {
	UserID         string    `json:"user_id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	PhoneNumber    string    `json:"phone_number"`
	Role           string    `json:"role"`
	LastConnection time.Time `json:"last_connection"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}
