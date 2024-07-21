package models

import (
	"github.com/Khvan-Group/common-library/models"
	categoryModels "shop-service/internal/categories/models"
)

// swagger:model
type Item struct {
	Code        string                  `json:"code" db:"code"`
	Name        string                  `json:"name" db:"name"`
	Description string                  `json:"description" db:"description"`
	Count       int                     `json:"count" db:"count"`
	Price       int                     `json:"price" db:"price"`
	Category    categoryModels.Category `json:"category" db:"category"`
}

// DTOs
// swagger:model
type ItemDto struct {
	Code        string `json:"code" db:"code"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Count       int    `json:"count" db:"count"`
	Price       int    `json:"price" db:"price"`
	Category    string `json:"category" db:"category"`
}

// swagger:model
type ItemSearch struct {
	Page     int
	Size     int
	SortBy   []models.SortField
	Name     *string `json:"name"`
	Code     *string `json:"code"`
	Category *string `json:"category"`
}

// swagger:model
type ItemBuy struct {
	Count    int    `json:"count" db:"count"`
	Username string `json:"username" db:"username"`
}

// swagger:model
type ItemBasketView struct {
	Code     string                  `json:"code" db:"code"`
	Name     string                  `json:"name" db:"name"`
	Count    int                     `json:"count" db:"count"`
	Price    int                     `json:"price" db:"price"`
	Total    int                     `json:"total" db:"total"`
	Category categoryModels.Category `json:"category" db:"category"`
}

// swagger:model
type ItemBasketHistoryView struct {
	Code     string                  `json:"code" db:"code"`
	Name     string                  `json:"name" db:"name"`
	Count    int                     `json:"count" db:"count"`
	Price    int                     `json:"price" db:"price"`
	Total    int                     `json:"total" db:"total"`
	Category categoryModels.Category `json:"category" db:"category"`
	PayedAt  string                  `json:"payed_at" db:"payed_at"`
}
