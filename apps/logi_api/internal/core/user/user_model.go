package user

import "time"

// User defines the structure for a user in the system.
type User struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Password  string    `json:"-"` // The '-' tag prevents this field from being serialized to JSON
	CreatedAt time.Time `json:"created_at"`
}
