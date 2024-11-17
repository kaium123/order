package model

import (
	"github.com/kaium123/order/internal/utils"
	"github.com/uptrace/bun"
	"math"
	"regexp"
	"time"
)

// Order represents the structure for the order details.
type Order struct {
	bun.BaseModel `bun:"table:orders"`

	ID                 int64       `json:"id" bun:"id,pk,autoincrement"`
	StoreID            int64       `json:"store_id" bun:"store_id,notnull"`
	MerchantOrderID    string      `json:"merchant_order_id,omitempty" bun:"merchant_order_id"`
	RecipientName      string      `json:"recipient_name" bun:"recipient_name,notnull" `
	RecipientPhone     string      `json:"recipient_phone" bun:"recipient_phone,notnull" `
	RecipientAddress   string      `json:"recipient_address" bun:"recipient_address,notnull" `
	RecipientCity      int64       `json:"recipient_city" bun:"recipient_city,notnull" `
	RecipientZone      int64       `json:"recipient_zone" bun:"recipient_zone,notnull" `
	RecipientArea      int64       `json:"recipient_area" bun:"recipient_area,notnull" `
	DeliveryType       OrderType   `json:"delivery_type" bun:"delivery_type,notnull" `
	ItemType           ItemType    `json:"item_type" bun:"item_type,notnull" `
	SpecialInstruction string      `json:"special_instruction,omitempty" bun:"special_instruction"`
	ItemQuantity       int         `json:"item_quantity" bun:"item_quantity,notnull" `
	ItemWeight         float64     `json:"item_weight" bun:"item_weight,notnull" `
	AmountToCollect    float64     `json:"amount_to_collect" bun:"amount_to_collect,notnull" `
	ItemDescription    string      `json:"item_description,omitempty" bun:"item_description"`
	OrderConsignmentID string      `json:"order_consignment_id" bun:"order_consignment_id,notnull"`
	OrderTypeID        int         `json:"order_type_id" bun:"order_type_id"`
	CodFee             float64     `json:"cod_fee" bun:"cod_fee"`
	PromoDiscount      float64     `json:"promo_discount" bun:"promo_discount"`
	Discount           float64     `json:"discount" bun:"discount"`
	DeliveryFee        float64     `json:"delivery_fee" bun:"delivery_fee"`
	OrderStatus        OrderStatus `json:"order_status" bun:"order_status,notnull"`
	OrderType          OrderType   `json:"order_type" bun:"order_type"`
	OrderAmount        float64     `json:"order_amount" bun:"order_amount"`
	TotalFee           float64     `json:"total_fee" bun:"total_fee"`
	UserID             int64       `json:"user_id" bun:"user_id"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" bun:"created_at,default:current_timestamp,notnull"`                             // Created timestamp
	UpdatedAt time.Time `json:"updated_at" bun:"updated_at,default:current_timestamp,nullzero,onupdate:current_timestamp"` // Updated timestamp
	DeletedAt time.Time `json:"deleted_at" bun:"deleted_at,soft_delete,nullzero"`                                          // Soft delete timestamp
}

func (o *Order) CalculateDeliveryFee(baseDeliveryFee float64) {
	if o.ItemWeight <= 0.5 {
		o.DeliveryFee = baseDeliveryFee
	} else if o.ItemWeight > 0.5 && o.ItemWeight <= 1 {
		o.DeliveryFee = baseDeliveryFee + 10
	} else {
		extra := math.Ceil(o.ItemWeight - 1)
		o.DeliveryFee = baseDeliveryFee + 10 + (extra * 15)
	}
}

func (o *Order) CalculateCodFee() {
	o.CodFee = utils.CalculatePercentage(o.AmountToCollect, 1)
}

func (o *Order) CalculateTotalFee() {
	o.TotalFee = o.CodFee + o.DeliveryFee
}

// Validate validates the Order fields and returns errors in the required format.
func (o *Order) Validate() *utils.ResponseError {
	responseError := &utils.ResponseError{
		Code:    "422",
		Message: "Please fix the given errors",
		Type:    "error",
		Errors:  make(map[string][]string),
	}

	// Regex pattern for Bangladeshi phone numbers
	phoneRegex := regexp.MustCompile(`^(01)[3-9]{1}[0-9]{8}$`)

	// Manual validations
	if o.StoreID <= 0 {
		responseError.Errors["store_id"] = append(responseError.Errors["store_id"], "The store field is required.")
		responseError.Errors["store_id"] = append(responseError.Errors["store_id"], "Wrong Store selected.")
	}
	if o.RecipientName == "" {
		responseError.Errors["recipient_name"] = append(responseError.Errors["recipient_name"], "The recipient name field is required.")
	}
	if o.RecipientPhone == "" {
		responseError.Errors["recipient_phone"] = append(responseError.Errors["recipient_phone"], "The recipient phone field is required.")
	} else if !phoneRegex.MatchString(o.RecipientPhone) {
		responseError.Errors["recipient_phone"] = append(responseError.Errors["recipient_phone"], "The recipient phone number is invalid.")
	}
	if o.RecipientAddress == "" {
		responseError.Errors["recipient_address"] = append(responseError.Errors["recipient_address"], "The recipient address field is required.")
	}
	if o.DeliveryType <= 0 {
		responseError.Errors["delivery_type"] = append(responseError.Errors["delivery_type"], "The delivery type field is required.")
	}
	if o.AmountToCollect <= 0 {
		responseError.Errors["amount_to_collect"] = append(responseError.Errors["amount_to_collect"], "The amount to collect field is required.")
	}
	if o.ItemQuantity <= 0 {
		responseError.Errors["item_quantity"] = append(responseError.Errors["item_quantity"], "The item quantity field is required.")
	}
	if o.ItemWeight <= 0 {
		responseError.Errors["item_weight"] = append(responseError.Errors["item_weight"], "The item weight field is required.")
	}
	if o.ItemType <= 0 {
		responseError.Errors["item_type"] = append(responseError.Errors["item_type"], "The item type field is required.")
	}

	if len(responseError.Errors) > 0 {
		return responseError
	}

	return nil
}

// DeleteRequest is the request parameter for deleting a todo
type DeleteRequest struct {
	UserId int64 `param:"user_id" validate:"required"`
	ID     int   `param:"id" validate:"required"`
}

// FindRequest is the request parameter for finding a todo
type FindRequest struct {
	UserId int64 `param:"user_id" validate:"required"`
	ID     int   `param:"id" validate:"required"`
}

type OrderCancelRequest struct {
	UserId        int64 `param:"user_id" validate:"required"`
	ConsignmentID string
}

type FindAllRequest struct {
	UserId         int64 `param:"user_id" validate:"required"`
	TransferStatus string
	Archive        int
	Limit          int
	Offset         int
}

type CreateOrderResponse struct {
	ConsignmentID   string  `json:"consignment_id"`
	MerchantOrderID string  `json:"merchant_order_id"`
	OrderStatus     string  `json:"order_status"`
	DeliveryFee     float64 `json:"delivery_fee"`
}

type FindAllResponse struct {
	Orders      []*OrderResponse `json:"orders"`
	Total       int              `json:"total"`
	CurrentPage int              `json:"current_page"`
	PerPage     int              `json:"per_page"`
	TotalInPage int              `json:"total_in_page"`
	LastPage    int              `json:"last_page"`
}

// Order represents the structure of the order response.
type OrderResponse struct {
	OrderConsignmentID string    `json:"order_consignment_id"`
	OrderCreatedAt     time.Time `json:"order_created_at"`
	OrderDescription   string    `json:"order_description"`
	MerchantOrderID    string    `json:"merchant_order_id"`
	RecipientName      string    `json:"recipient_name"`
	RecipientAddress   string    `json:"recipient_address"`
	RecipientPhone     string    `json:"recipient_phone"`
	OrderAmount        float64   `json:"order_amount"`
	TotalFee           float64   `json:"total_fee"`
	Instruction        string    `json:"instruction"`
	OrderTypeID        int       `json:"order_type_id"`
	CODFee             float64   `json:"cod_fee"`
	PromoDiscount      float64   `json:"promo_discount"`
	Discount           float64   `json:"discount"`
	DeliveryFee        float64   `json:"delivery_fee"`
	OrderStatus        string    `json:"order_status"`
	OrderType          string    `json:"order_type"`
	ItemType           string    `json:"item_type"`
}

type OrderStatus int

const (
	Pending    OrderStatus = iota // 0
	Processing                    // 1
	Completed                     // 2
)

// String provides a string representation of the OrderStatus enum.
func (s OrderStatus) String() string {
	switch s {
	case Pending:
		return "Pending"
	case Processing:
		return "Processing"
	case Completed:
		return "Completed"
	default:
		return "Unknown"
	}
}

// OrderType represents the type of an order.
type OrderType int

const (
	UnknownOrderType OrderType = iota // 0
	Pickup
	Delivery // 1
)

// String provides a string representation of the OrderType enum.
func (o OrderType) String() string {
	switch o {
	case Delivery:
		return "Delivery"
	case Pickup:
		return "Pickup"
	default:
		return "Unknown"
	}
}

// ItemType represents the type of an item.
type ItemType int

const (
	UnknownItemType ItemType = iota // 0
	Document                        // 1
	Parcel
	Other // 2
)

// String provides a string representation of the ItemType enum.
func (i ItemType) String() string {
	switch i {
	case Parcel:
		return "Parcel"
	case Document:
		return "Document"
	case Other:
		return "Other"
	default:
		return "Unknown"
	}
}
