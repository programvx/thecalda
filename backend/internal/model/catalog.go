package model

// Catalog maps to the public.catalogs table — a named collection of items.
type Catalog struct {
	Base
	Slug        string  `json:"slug"`
	Name        string  `json:"name"`
	Description *string `json:"description"`
	IsActive    bool    `json:"isActive"`
}

// TableName tells GORM which table this model maps to.
func (Catalog) TableName() string {
	return "catalogs"
}

// CatalogCreate is the request body for creating a catalog.
type CatalogCreate struct {
	Slug        string  `json:"slug" binding:"required"`
	Name        string  `json:"name" binding:"required"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"isActive"`
}

// CatalogUpdate is the request body for updating a catalog. Every field is
// optional; only the fields provided are applied.
type CatalogUpdate struct {
	Slug        *string `json:"slug"`
	Name        *string `json:"name"`
	Description *string `json:"description"`
	IsActive    *bool   `json:"isActive"`
}

// CatalogWithItems is a catalog together with the items it contains.
type CatalogWithItems struct {
	Catalog
	Items []Item `json:"items"`
}
