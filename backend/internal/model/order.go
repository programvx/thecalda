package model

import (
	"time"

	"github.com/google/uuid"
)

// Order maps to the public.orders table. A row is either a shopping cart
// (Type "cart") or a placed order (Type "order"); see migration 003_orders.
type Order struct {
	Base
	UserID            int64       `json:"-"`
	Type              string      `json:"type"`
	Status            string      `json:"status"`
	CurrencyID        int64       `json:"-"`
	BillingAddressID  *int64      `json:"-"`
	ShippingAddressID *int64      `json:"-"`
	PaymentMethodID   *int64      `json:"-"`
	TotalAmount       float64     `json:"totalAmount" gorm:"->"`
	OrderNumber       *string     `json:"orderNumber"`
	Notes             *string     `json:"notes"`
	PlacedAt          *time.Time  `json:"placedAt"`
	Items             []OrderItem `json:"items" gorm:"foreignKey:OrderID"`
}

// TableName tells GORM which table this model maps to.
func (Order) TableName() string {
	return "orders"
}

// OrderItem maps to the public.order_items table — one line of an order.
// ItemName and the unit_* prices are snapshots taken when the line is created;
// UnitPriceDiscounted and LineTotal are computed by the database.
type OrderItem struct {
	Base
	OrderID             int64   `json:"-"`
	ItemID              int64   `json:"-"`
	ItemName            string  `json:"itemName"`
	Quantity            int     `json:"quantity"`
	UnitPrice           float64 `json:"unitPrice"`
	UnitDiscount        float64 `json:"unitDiscount"`
	UnitPriceDiscounted float64 `json:"unitPriceDiscounted" gorm:"->"`
	LineTotal           float64 `json:"lineTotal" gorm:"->"`
}

// TableName tells GORM which table this model maps to.
func (OrderItem) TableName() string {
	return "order_items"
}

// OrderItemInput is a single requested line in an OrderCreate request.
type OrderItemInput struct {
	ItemUID  uuid.UUID `json:"itemUid"`
	Quantity int       `json:"quantity" binding:"required,gt=0"`
}

// OrderCreate is the request body for creating a cart or order. The server
// snapshots each item's name and price from the catalog.
type OrderCreate struct {
	Type  string           `json:"type"` // "cart" (default) or "order"
	Notes *string          `json:"notes"`
	Items []OrderItemInput `json:"items" binding:"required,min=1,dive"`
}

// OrderUpdate is the request body for updating an order. Fields are optional;
// only those provided are applied.
type OrderUpdate struct {
	Status *string `json:"status"`
	Notes  *string `json:"notes"`
}

// OrderItemCreate is the request body for adding an item to an order. The
// server snapshots the catalog item's name and price onto the new line.
type OrderItemCreate struct {
	OrderUID uuid.UUID `json:"orderUid"`
	ItemUID  uuid.UUID `json:"itemUid"`
	Quantity int       `json:"quantity" binding:"required,gt=0"`
}

// OrderItemUpdate is the request body for updating an order item. Only the
// quantity of an existing line can change; the price snapshot is fixed.
type OrderItemUpdate struct {
	Quantity int `json:"quantity" binding:"required,gt=0"`
}
