package store

import (
	"github.com/Khvan-Group/common-library/errors"
	"github.com/go-resty/resty/v2"
	"github.com/jmoiron/sqlx"
	"shop-service/internal/baskets/models"
	"shop-service/internal/clients"
	"shop-service/internal/common"
	"shop-service/internal/core/rabbitmq"
	"shop-service/internal/db"
	itemModel "shop-service/internal/items/models"
)

type BasketStore struct {
	db     *sqlx.DB
	client *resty.Client
}

func New() *BasketStore {
	return &BasketStore{
		db:     db.DB,
		client: resty.New(),
	}
}

func (s *BasketStore) Save(input models.BasketSave) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var basket *models.Basket
		var itemCount int
		err := tx.Get(&basket, "select * from t_basket where username = $1 and item_code = $2", input.Username, input.ItemCode)
		if err != nil {
			basket = nil
		}

		err = tx.Get(&itemCount, "select count from t_items where code = $1", input.ItemCode)
		if err != nil {
			return errors.NewBadRequest("Данного товара не существует.")
		}

		if input.Action == "ADD" {
			if basket == nil {
				_, err = tx.NamedExec("insert into t_basket (username, item_code) values (:username, :item_code)", input)
				if err != nil {
					panic(err)
				}
			} else {
				if basket.Count+1 > itemCount {
					return errors.NewBadRequest("Нельзя добавить в корзину количетсво товаров превышающее количество товара на складе.")
				}

				_, err = tx.NamedExec("update t_basket set count = count+1 where username = :username and item_code = :item_code", input)
				if err != nil {
					panic(err)
				}
			}
		} else {
			if basket == nil {
				return errors.NewBadRequest("Данного товара нет в корзине пользователя.")
			} else {
				if itemCount-1 == 0 {
					_, err = tx.Exec("delete from t_basket where username = :username and item_code = :item_code", input)
					if err != nil {
						panic(err)
					}
				} else {
					_, err = tx.Exec("update t_basket set count = count-1 where username = :username and item_code = :item_code", input)
					if err != nil {
						panic(err)
					}
				}
			}
		}

		return nil
	})
}

func (s *BasketStore) FindByUser(username string) models.BasketView {
	var items []itemModel.ItemBasketView
	var totalSum int

	query := `
		select i.code, i.name, b.count,
		       i.price, i.price*b.count as total, i.category as "category.code", c.name as "category.name"
		from t_basket b
			inner join t_items i on i.code = b.item_code
			inner join t_categories c on c.code = i.category
		where b.username = $1
	`
	err := s.db.Select(&items, query, username)
	if err != nil || items == nil {
		return models.BasketView{
			Items: make([]itemModel.ItemBasketView, 0),
			Total: 0,
		}
	}

	for _, i := range items {
		totalSum += i.Total
	}

	return models.BasketView{
		Items: items,
		Total: totalSum,
	}
}

func (s *BasketStore) Remove(username, itemCode string) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var existsItemBasket bool
		tx.Get(&existsItemBasket, "select exists(select 1 from t_basket where username = $1 and item_code = $2)")

		if !existsItemBasket {
			return errors.NewBadRequest("Данного товара нет в корзине пользователя.")
		}

		_, err := tx.Exec("delete from t_basket where username = $1 and item_code = $2", username, itemCode)
		if err != nil {
			panic(err)
		}

		return nil
	})
}

func (s *BasketStore) Pay(username string) *errors.CustomError {
	return db.StartTransaction(func(tx *sqlx.Tx) *errors.CustomError {
		var totalSum int
		var items []itemModel.ItemBasketView
		query := `
			select i.code, i.name, b.count,
				   i.price, i.price*b.count as total, i.category as "category.code", c.name as "category.name"
			from t_basket b
				inner join t_items i on i.code = b.item_code
				inner join t_categories c on c.code = i.category
			where b.username = $1
		`
		err := tx.Select(&items, query, username)
		if err != nil {
			return errors.NewBadRequest("Корзина пользователя пустая.")
		}

		for _, i := range items {
			totalSum += i.Total
		}

		wallet, customErr := clients.GetWalletByUser(username, s.client)
		if customErr != nil {
			return customErr
		}

		if wallet.Total < totalSum {
			return errors.NewBadRequest("Недостаточно средств для оплаты.")
		}

		query = `
			do $$
			begin
			    insert into t_basket_history
				select username, item_code, count, now() from t_basket where username = $1;
			end;
			$$ language plpgsql
		`
		_, err = tx.Exec(query, username)
		if err != nil {
			panic(err)
		}

		_, err = tx.Exec("delete from t_basket where username = $1", username)
		if err != nil {
			panic(err)
		}

		msg := common.WalletUpdate{
			Username: username,
			Total:    100,
			Action:   common.WALLET_TOTAL_SUBSTRUCT,
		}

		if customErr = rabbitmq.SendToUpdateWallet(msg); err != nil {
			return customErr
		}

		return nil
	})
}

func (s *BasketStore) GetHistoryByUser(username string) []models.BasketHistoryView {
	var items []itemModel.ItemBasketHistoryView

	query := `
		select i.code, i.name, b.count,
		       i.price, i.price*b.count as total, i.category as "category.code", c.name as "category.name", b.payed_at
		from t_basket_history b
			inner join t_items i on i.code = b.item_code
			inner join t_categories c on c.code = i.category
		where b.username = $1
		order by payed_at
	`
	err := s.db.Select(&items, query, username)
	if err != nil {
		return make([]models.BasketHistoryView, 0)
	}

	groupsMap := make(map[string][]itemModel.ItemBasketHistoryView)
	response := make([]models.BasketHistoryView, 0)
	payedAtSet := make(map[string]struct{})
	for _, item := range items {
		payedAtSet[item.PayedAt] = struct{}{}
		groupsMap[item.PayedAt] = append(groupsMap[item.PayedAt], item)
	}

	for payedAt := range payedAtSet {
		totalSum := 0
		for _, item := range groupsMap[payedAt] {
			totalSum += item.Total
		}

		basket := models.BasketHistoryView{
			Items:   groupsMap[payedAt],
			Total:   totalSum,
			PayedAt: payedAt,
		}

		response = append(response, basket)
	}

	return response
}
