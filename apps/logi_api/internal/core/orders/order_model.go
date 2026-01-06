package orders

import "time"

// Order defines the structure for an order in the system.
type Order struct {
	OrderID            string      `json:"order_id" gorm:"primaryKey;type:varchar(36)"`
	OrderNumber        string      `json:"order_number" gorm:"type:varchar(6);not null"`
	CreatedBy          string      `json:"created_by" gorm:"type:varchar(36);not null"`
	CreatedByUsername  string      `json:"created_by_username" gorm:"-"` // Not stored in orders table, fetched from users
	AssignedTo         *string     `json:"assigned_to" gorm:"type:varchar(36)"`
	AssignedToUsername *string     `json:"assigned_to_username" gorm:"-"` // Not stored in orders table, fetched from users
	OrderName          string      `json:"order_name" gorm:"type:varchar(150);not null"`
	OrderPhoneNumber   string      `json:"order_phone_number" gorm:"type:varchar(16);not null"`
	OrderEmail         *string     `json:"order_email" gorm:"type:text"`
	OrderCedula        *string     `json:"order_cedula" gorm:"type:varchar(10)"`
	DeliveryAddress    string      `json:"delivery_address" gorm:"type:text;not null"`
	Status             string      `json:"status" gorm:"type:order_status;not null;default:'pending'"`
	CreatedAt          time.Time   `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt          time.Time   `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	Items              []OrderItem `json:"items,omitempty" gorm:"-"` // Not stored in orders table, fetched separately
}

// OrderItem defines the structure for items within an order.
type OrderItem struct {
	ItemID        string    `json:"item_id" gorm:"primaryKey;type:varchar(36)"`
	OrderID       string    `json:"order_id" gorm:"type:varchar(36);not null"`
	ProductName   string    `json:"product_name" gorm:"type:text;not null"`
	Quantity      int       `json:"quantity" gorm:"type:int;not null"`
	RespondedForm bool      `json:"responded_form" gorm:"type:bool;not null;default:false"`
	CreatedAt     time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt     time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}

// OrderForm defines the structure for order satisfaction forms.
type OrderForm struct {
	FormID         string    `json:"form_id" gorm:"primaryKey;type:varchar(36)"`
	PublicID       string    `json:"public_id" gorm:"type:varchar(6);unique;not null"`
	OrderID        string    `json:"order_id" gorm:"type:varchar(36);not null"`
	DriverID       *string   `json:"driver_id" gorm:"type:varchar(36)"`
	DriverRating   *int      `json:"driver_rating" gorm:"type:int"`
	CargoCondition *string   `json:"cargo_condition" gorm:"type:varchar(50)"`
	Comments       *string   `json:"comments" gorm:"type:text"`
	IsFinished     bool      `json:"is_finished" gorm:"type:bool;not null;default:false"`
	DriverName     string    `json:"driver_name" gorm:"-"`
	DriverEmail    string    `json:"driver_email" gorm:"-"`
	CreatedAt      time.Time `json:"created_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"type:timestamp with time zone;default:current_timestamp"`
}
