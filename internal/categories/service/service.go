package service

import (
	"github.com/Khvan-Group/common-library/errors"
	"shop-service/internal/categories/models"
	"shop-service/internal/categories/store"
)

type CategoryService interface {
	Save(input models.Category) *errors.CustomError
	FindAll() []models.Category
	Delete(code string) *errors.CustomError
}

type Categories struct {
	Service CategoryService
}

func New() *Categories {
	return &Categories{
		Service: store.New(),
	}
}
