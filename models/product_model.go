package models

type Product struct {
	ID         int    `json:"id"`
	Name       string `json:"name"`
	Price      int    `json:"price"`
	Stock      int    `json:"stock"`
	CategoryID int    `json:"category_id"`
}

type ProductDTO struct {
	ID         int       `json:"id"`
	Name       string    `json:"name"`
	Price      int       `json:"price"`
	Stock      int       `json:"stock"`
	CategoryID int       `json:"-"`
	Category   *Category `json:"category,omitempty"` // Eager loaded category, omitempty means it can be nil and will be omitted in JSON
}
