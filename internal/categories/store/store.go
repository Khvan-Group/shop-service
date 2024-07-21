package store

import (
	"github.com/Khvan-Group/common-library/errors"
	"github.com/jmoiron/sqlx"
	"shop-service/internal/categories/models"
	"shop-service/internal/db"
	"strings"
)

type CategoryStore struct {
	db *sqlx.DB
}

func New() *CategoryStore {
	return &CategoryStore{
		db: db.DB,
	}
}

func (s *CategoryStore) Save(input models.Category) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		input.Code = strings.ToUpper(input.Code)

		var existsCategory bool
		tx.Get(&existsCategory, "select exists(select 1 from t_categories where upper(code) = $1)", input.Code)
		if !existsCategory {
			_, err := tx.NamedExec("insert into t_categories values(:code, :name)", input)
			if err != nil {
				panic(err)
			}
		} else {
			_, err := tx.NamedExec("update t_categories set name = :name where upper(code) = :code", input)
			if err != nil {
				panic(err)
			}
		}

		return nil
	})
}

func (s *CategoryStore) FindAll() []models.Category {
	var response []models.Category
	query := "select * from t_categories"
	if err := s.db.Select(&response, query); err != nil || response == nil {
		return make([]models.Category, 0)
	}

	return response
}

func (s *CategoryStore) Delete(code string) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var existsCategory bool
		tx.Get(&existsCategory, "select exists(select 1 from t_categories where upper(code) = $1)", code)
		if !existsCategory {
			return errors.NewBadRequest("Данной категории не существует.")
		}

		_, err := tx.Exec("delete from t_categories where code = $1", code)
		if err != nil {
			panic(err)
		}

		return nil
	})
}
