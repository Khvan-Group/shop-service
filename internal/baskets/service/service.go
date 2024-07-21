package service

import (
	"github.com/Khvan-Group/common-library/errors"
	"shop-service/internal/baskets/models"
	"shop-service/internal/baskets/store"
)

type BasketService interface {
	Save(input models.BasketSave) *errors.CustomError
	FindByUser(username string) models.BasketView
	Remove(username, itemCode string) *errors.CustomError
	Pay(username string) *errors.CustomError
	GetHistoryByUser(username string) []models.BasketHistoryView
}

type Baskets struct {
	Service BasketService
}

func New() *Baskets {
	return &Baskets{
		Service: store.New(),
	}
}
