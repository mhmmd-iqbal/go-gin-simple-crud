package models

type Product struct {
	ID          int64   `gorm:"primaryKey" json:"id"`
	Name        string  `json:"name"`
	Description string  `gorm:"type:text" json:"description"`
	Price       float64 `json:"price"`
	SKU         string  `json:"sku"`
}

type ListProductResponse struct {
	SKU   string  `json:"sku"`
	Name  string  `json:"name"`
	Price float64 `json:"price"`
}

type DetailProductResponse struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
}
