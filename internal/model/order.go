package model

import (
	"github.com/uptrace/bun"
	"time"
)

// Order represents the structure for the order details.
type Order struct {
	// BaseModel includes fields like ID, CreatedAt, UpdatedAt, and DeletedAt.
	bun.BaseModel `bun:"table:orders"` // Indicates the table name in the database

	// The fields added in your SQL statement
	ID                 int64   `json:"id" bun:"id,pk,autoincrement"`                        // Primary Key
	StoreID            int64   `json:"store_id" bun:"store_id,notnull"`                     // Store ID
	MerchantOrderID    string  `json:"merchant_order_id,omitempty" bun:"merchant_order_id"` // Optional field
	RecipientName      string  `json:"recipient_name" bun:"recipient_name,notnull" validate:"required"`
	RecipientPhone     string  `json:"recipient_phone" bun:"recipient_phone,notnull" validate:"required"`
	RecipientAddress   string  `json:"recipient_address" bun:"recipient_address,notnull" validate:"required"`
	RecipientCity      int64   `json:"recipient_city" bun:"recipient_city,notnull" validate:"required"`
	RecipientZone      int64   `json:"recipient_zone" bun:"recipient_zone,notnull" validate:"required"`
	RecipientArea      int64   `json:"recipient_area" bun:"recipient_area,notnull" validate:"required"`
	DeliveryType       int64   `json:"delivery_type" bun:"delivery_type,notnull" validate:"required"`
	ItemType           int64   `json:"item_type" bun:"item_type,notnull" validate:"required"`
	SpecialInstruction string  `json:"special_instruction,omitempty" bun:"special_instruction"` // Optional field
	ItemQuantity       int     `json:"item_quantity" bun:"item_quantity,notnull" validate:"required"`
	ItemWeight         float64 `json:"item_weight" bun:"item_weight,notnull" validate:"required"`
	AmountToCollect    float64 `json:"amount_to_collect" bun:"amount_to_collect,notnull" validate:"required"`
	ItemDescription    string  `json:"item_description,omitempty" bun:"item_description"` // Optional field

	// Newly added fields
	OrderConsignmentID string `json:"order_consignment_id" bun:"order_consignment_id,notnull"` // Consignment ID
	OrderTypeID        int    `json:"order_type_id" bun:"order_type_id"`                       // Order type ID
	CodFee             int    `json:"cod_fee" bun:"cod_fee"`                                   // COD fee
	PromoDiscount      int    `json:"promo_discount" bun:"promo_discount"`                     // Promo discount
	Discount           int    `json:"discount" bun:"discount"`                                 // Discount
	DeliveryFee        int    `json:"delivery_fee" bun:"delivery_fee"`                         // Delivery fee
	OrderStatus        string `json:"order_status" bun:"order_status,notnull"`                 // Order status (Pending)
	OrderType          string `json:"order_type" bun:"order_type,notnull"`                     // Order type (Delivery)
	OrderAmount        int    `json:"order_amount" bun:"order_amount"`                         // Order amount
	TotalFee           int    `json:"total_fee" bun:"total_fee"`                               // Total fee

	// Timestamps
	CreatedAt time.Time `json:"created_at" bun:"created_at,default:current_timestamp,notnull"`                             // Created timestamp
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp,nullzero,onupdate:current_timestamp"` // Updated timestamp
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`                                          // Soft delete timestamp
}

// DeleteRequest is the request parameter for deleting a todo
type DeleteRequest struct {
	ID int `param:"id" validate:"required"`
}

// FindRequest is the request parameter for finding a todo
type FindRequest struct {
	ID int `param:"id" validate:"required"`
}

type OrderCancelRequest struct {
	ConsignmentID string
}

type FindAllRequest struct {
	TransferStatus string
	Archive        int
	Limit          int
	Offset         int
}

type CreateOrderResponse struct {
	ConsignmentID   string `json:"consignment_id"`
	MerchantOrderID string `json:"merchant_order_id"`
	OrderStatus     string `json:"order_status"`
	DeliveryFee     int    `json:"delivery_fee"`
}

type FindAllResponse struct {
	Orders      []*Order `json:"orders"`
	Total       int      `json:"total"`
	CurrentPage int      `json:"current_page"`
	PerPage     int      `json:"per_page"`
	TotalInPage int      `json:"total_in_page"`
	LastPage    int      `json:"last_page"`
}
