package orders

import "time"

// Order defines the structure for an order in the system.
type Order struct {
	OrderID         string    `json:"order_id" gorm:"primaryKey;type:varchar(36)"`
	CreatedBy       string    `json:"created_by" gorm:"type:varchar(36);not null"`
	AssignedTo      string    `json:"assigned_to" gorm:"type:varchar(36)"`
	DeliveryAddress string    `json:"delivery_address" gorm:"type:text;not null"`
	Status          string    `json:"status" gorm:"type:order_status;not null;default:'pending'"`
	CreatedAt       time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt       time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}
