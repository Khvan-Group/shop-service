package models

import "shop-service/internal/items/models"

type Basket struct {
	Username string `json:"username" db:"username"`
	ItemCode string `json:"item_code" db:"item_code"`
	Count    int    `json:"count" db:"count"`
}

type BasketView struct {
	Items []models.ItemBasketView `json:"items"`
	Total int                     `json:"total"`
}

type BasketSave struct {
	Username string `json:"username" db:"username"`
	ItemCode string `json:"item_code" db:"item_code"`
	Action   string `json:"action"`
}

type BasketHistory struct {
	Username string `json:"username" db:"username"`
	ItemCode string `json:"item_code" db:"item_code"`
	Count    int    `json:"count" db:"count"`
	PayedAt  string `json:"payed_at" db:"payed_at"`
}

type BasketHistoryView struct {
	Items   []models.ItemBasketHistoryView `json:"items"`
	Total   int                            `json:"total"`
	PayedAt string                         `json:"payed_at"`
}

type GroupedBasketHistoryView struct {
	Basket  BasketHistoryView `json:"items"`
	Total   int               `json:"total"`
	PayedAt string            `json:"payed_at"`
}
