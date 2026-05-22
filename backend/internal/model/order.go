package model

import (
	"time"

	"github.com/google/uuid"
)

// Order maps to the public.orders table. A row is either a shopping cart
// (Type "cart") or a placed order (Type "order"); see migration 003_orders.
type Order struct {
	Base
	UserID            int64          `json:"-"`
	Type              string         `json:"type"`
	Status            string         `json:"status"`
	CurrencyID        int64          `json:"-"`
	BillingAddressID  *int64         `json:"-"`
	ShippingAddressID *int64         `json:"-"`
	PaymentMethodID   *int64         `json:"-"`
	TotalAmount       float64        `json:"totalAmount" gorm:"->"`
	OrderNumber       *string        `json:"orderNumber"`
	Notes             *string        `json:"notes"`
	PlacedAt          *time.Time     `json:"placedAt"`
	Items             []OrderItem    `json:"items" gorm:"foreignKey:OrderID"`
	BillingAddress    *Address       `json:"billingAddress" gorm:"foreignKey:BillingAddressID;references:ID"`
	ShippingAddress   *Address       `json:"shippingAddress" gorm:"foreignKey:ShippingAddressID;references:ID"`
	PaymentMethod     *PaymentMethod `json:"paymentMethod" gorm:"foreignKey:PaymentMethodID;references:ID"`
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
	Item                *Item   `json:"item" gorm:"foreignKey:ItemID;references:ID"`
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

// Address maps to the public.addresses table — a postal address an order
// references for billing and/or shipping.
type Address struct {
	Base
	FirstName    string  `json:"firstName"`
	LastName     string  `json:"lastName"`
	Email        string  `json:"email"`
	Phone        *string `json:"phone"`
	AddressLine1 string  `json:"addressLine1"`
	AddressLine2 *string `json:"addressLine2"`
	PostalCode   string  `json:"postalCode"`
	City         string  `json:"city"`
	Country      string  `json:"country"`
}

// TableName tells GORM which table this model maps to.
func (Address) TableName() string {
	return "addresses"
}

// PaymentMethod maps to the public.payment_methods reference table.
type PaymentMethod struct {
	Base
	Code     string `json:"code"`
	Name     string `json:"name"`
	IsActive bool   `json:"isActive"`
}

// TableName tells GORM which table this model maps to.
func (PaymentMethod) TableName() string {
	return "payment_methods"
}

// AddressInput is one postal address in a checkout request. The contact email
// is shared across addresses and carried on OrderCheckout.
type AddressInput struct {
	FirstName    string  `json:"firstName" binding:"required"`
	LastName     string  `json:"lastName" binding:"required"`
	Phone        *string `json:"phone"`
	AddressLine1 string  `json:"addressLine1" binding:"required"`
	AddressLine2 *string `json:"addressLine2"`
	City         string  `json:"city" binding:"required"`
	PostalCode   string  `json:"postalCode" binding:"required"`
	Country      string  `json:"country" binding:"required"`
}

// OrderCheckout is the request body for checking out a cart — placing the
// order. BillingAddress is optional; when omitted, billing reuses shipping.
type OrderCheckout struct {
	Email             string        `json:"email" binding:"required,email"`
	Note              *string       `json:"note"`
	PaymentMethodCode string        `json:"paymentMethodCode" binding:"required"`
	ShippingAddress   AddressInput  `json:"shippingAddress"`
	BillingAddress    *AddressInput `json:"billingAddress"`
}
