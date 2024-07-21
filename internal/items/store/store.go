package store

import (
	"fmt"
	"github.com/Khvan-Group/common-library/errors"
	commonModels "github.com/Khvan-Group/common-library/models"
	"github.com/jmoiron/sqlx"
	"shop-service/internal/db"
	"shop-service/internal/items/models"
	"strings"
)

type ItemStore struct {
	db *sqlx.DB
}

func New() *ItemStore {
	return &ItemStore{
		db: db.DB,
	}
}

func (s *ItemStore) Create(input models.ItemDto) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		input.Code = strings.ToUpper(input.Code)

		var existsItem bool
		tx.Get(&existsItem, "select exists(select 1 from t_items where upper(code) = $1)", input.Code)
		if existsItem {
			return errors.NewBadRequest("Такой товар уже существует.")
		}

		if err := validateEditing(tx, input); err != nil {
			return err
		}

		_, execErr := tx.NamedExec("insert into t_items values (:code, :name, :description, :count, :price, :category)", input)
		if execErr != nil {
			panic(execErr)
		}

		return nil
	})
}

func (s *ItemStore) Update(input models.ItemDto) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		input.Code = strings.ToUpper(input.Code)

		var existsItem bool
		tx.Get(&existsItem, "select exists(select 1 from t_items where upper(code) = $1)", input.Code)
		if !existsItem {
			return errors.NewBadRequest("Данный товар не существует.")
		}

		if err := validateEditing(tx, input); err != nil {
			return err
		}

		query := `
			update t_items 
				set name = :name, 
				    description = :description, 
				    count = :count, 
				    price = :price, 
				    category = :category
			where code = :code
		`
		_, execErr := tx.NamedExec(query, input)
		if execErr != nil {
			panic(execErr)
		}

		return nil
	})
}

func (s *ItemStore) FindAll(input models.ItemSearch) commonModels.Page {
	var items []models.Item
	var totalElements int
	query := buildQuery(input)
	if err := s.db.Select(&items, query, input.Size, input.Page*input.Size); err != nil || items == nil {
		items = make([]models.Item, 0)
	}

	query = "select count(*) from t_items "
	query += addQueryByValidate(input)
	err := s.db.Get(&totalElements, query)
	if err != nil {
		totalElements = len(items)
	}

	if totalElements == 0 {
		input.Size = 0
	}

	return commonModels.Page{
		Result:        items,
		Page:          input.Page,
		Size:          input.Size,
		TotalElements: totalElements,
	}
}

func (s *ItemStore) FindByCode(code string) (*models.Item, *errors.CustomError) {
	var item models.Item
	query := `
		select i.code, i.name, i.description, i.count, i.price, i.category as "category.code", c.name as "category.name"
		from t_items i
		inner join t_categories c on c.code = i.category
		where upper(i.category) = upper(:code)
	`
	err := s.db.Get(&item, query, code)

	if err != nil {
		return nil, errors.NewBadRequest("Товар с данным кодом не найден.")
	}

	return &item, nil
}

func (s *ItemStore) Delete(code string) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var existsItem bool
		tx.Get(&existsItem, "select exists(select 1 from t_items where upper(code) = $1)", code)
		if !existsItem {
			return errors.NewBadRequest("Данный товар не существует.")
		}

		_, err := tx.Exec("delete from t_items where code = $1", code)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func (s *ItemStore) Buy(code string, input models.ItemBuy) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var item models.Item
		err := tx.Get(&item, "select * from t_items where upper(code) = $1", code)
		if err != nil {
			return errors.NewBadRequest("Данный товар не существует.")
		}

		if item.Count < input.Count {
			return errors.NewBadRequest("Нельзя купить товара больше, чем имеется.")
		}

		return nil
	})
}

func buildQuery(input models.ItemSearch) string {
	query := `
		select i.code, i.name, i.description, i.count, i.price, i.category as "category.code", c.name as "category.name"
		from t_items i
		inner join t_categories c on c.code = i.category
	`

	query += addQueryByValidate(input)

	if len(input.SortBy) > 0 {
		query += " order by "
		sortJoins := make([]string, 0)

		for _, field := range input.SortBy {
			sortJoins = append(sortJoins, fmt.Sprintf("%s %s ", field.SortBy, field.Direction))
		}

		query += strings.Join(sortJoins, ", ")
	}

	return query + "limit $1 offset $2"
}

func addQueryByValidate(input models.ItemSearch) string {
	query := " where true "
	if input.Code != nil && len(*input.Code) > 0 {
		query += "and upper(i.code) = upper(" + *input.Code + ") "
	}

	if input.Name != nil && len(*input.Name) > 0 {
		query += "and i.name like '%" + *input.Name + "%' "
	}

	if input.Category != nil && len(*input.Category) > 0 {
		query += "and upper(i.category) = upper(" + *input.Category + ") "
	}

	return query
}

func validateEditing(tx *sqlx.Tx, input models.ItemDto) *errors.CustomError {
	var existsCategory bool
	tx.Get(&existsCategory, "select exists(select 1 from t_categories where upper(code) = $1)", input.Category)
	if !existsCategory {
		return errors.NewBadRequest("Данной категории товаров не существует.")
	}

	if len(input.Name) == 0 {
		return errors.NewBadRequest("Название товара не может быть пустым.")
	}

	if input.Count < 0 {
		return errors.NewBadRequest("Количество товара не может быть отрицательным числом.")
	}

	if input.Price < 0 {
		return errors.NewBadRequest("Цена товара не может быть отрицательным числом.")
	}

	return nil
}
