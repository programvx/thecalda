package model

// Item maps to the public.items table — an e-commerce product.
type Item struct {
	Base
	Slug            string         `json:"slug"`
	SKU             *string        `json:"sku"`
	Name            string         `json:"name"`
	Description     *string        `json:"description"`
	Price           float64        `json:"price"`
	Discount        float64        `json:"discount"`
	PriceDiscounted float64        `json:"priceDiscounted" gorm:"->"`
	CurrencyID      int64          `json:"-"`
	StockID         *int64         `json:"-"`
	Stock           *ItemStock     `json:"stock" gorm:"foreignKey:StockID;references:ID"`
	IsActive        bool           `json:"isActive"`
	Medias          []ItemMedia    `json:"medias" gorm:"foreignKey:ItemID"`
	Properties      []ItemProperty `json:"properties" gorm:"foreignKey:ItemID"`
}

// TableName tells GORM which table this model maps to.
func (Item) TableName() string {
	return "items"
}

// ItemStock maps to the public.item_stock table — the stock level and
// availability status for an item, linked from items.stock_id.
type ItemStock struct {
	Base
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

// TableName tells GORM which table this model maps to.
func (ItemStock) TableName() string {
	return "item_stock"
}

// ItemMedia maps to the public.item_medias table — media attached to an item.
// The lowest position is the item's primary media.
type ItemMedia struct {
	Base
	MediaType string  `json:"mediaType"`
	URL       string  `json:"url"`
	Alt       *string `json:"alt"`
	Position  int     `json:"position"`
	ItemID    int64   `json:"-"`
}

// TableName tells GORM which table this model maps to.
func (ItemMedia) TableName() string {
	return "item_medias"
}

// ItemProperty maps to the public.item_properties table — a label/value
// attribute attached to an item.
type ItemProperty struct {
	Base
	Label  string `json:"label"`
	Value  string `json:"value"`
	ItemID int64  `json:"-"`
}

// TableName tells GORM which table this model maps to.
func (ItemProperty) TableName() string {
	return "item_properties"
}
