package service

import (
	"github.com/Khvan-Group/common-library/errors"
	commonModels "github.com/Khvan-Group/common-library/models"
	"shop-service/internal/items/models"
	"shop-service/internal/items/store"
)

type ItemService interface {
	Create(input models.ItemDto) *errors.CustomError
	Update(input models.ItemDto) *errors.CustomError
	FindAll(input models.ItemSearch) commonModels.Page
	FindByCode(code string) (*models.Item, *errors.CustomError)
	Delete(code string) *errors.CustomError
	Buy(code string, input models.ItemBuy) *errors.CustomError
}

type Items struct {
	Service ItemService
}

func New() *Items {
	return &Items{
		Service: store.New(),
	}
}
