package api

import (
	"encoding/json"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"shop-service/internal/baskets/models"
	"shop-service/internal/common"
)

// Save
// @Summary Добавить/Удалить один товар из корзины
// @ID save-to-basket
// @Accept json
// @Produce json
// @Param input body models.BasketSave true "Сохранить или удалить товар из корзины"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /basket [post]
// @Security ApiKeyAuth
func (a *API) Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	var input models.BasketSave

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	currentUser := common.GetJwtUser(r)
	if currentUser.Role != constants.ADMIN && currentUser.Login != input.Username {
		errors.HandleError(w, errors.NewForbidden("Доступ запрещен."))
		return
	}

	saveErr := a.baskets.Service.Save(input)
	if saveErr != nil {
		errors.HandleError(w, saveErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FindByUser
// @Summary Добавить/Удалить один товар из корзины
// @ID find-basket-by-user
// @Accept json
// @Produce json
// @Param username path string true "Логин пользователя"
// @Success 200 {array} models.BasketView
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /basket/{username} [get]
// @Security ApiKeyAuth
func (a *API) FindByUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if len(username) == 0 {
		errors.HandleError(w, errors.NewBadRequest("Логин пользователя обязателен к заполнению."))
		return
	}

	currentUser := common.GetJwtUser(r)
	if currentUser.Role != constants.ADMIN && currentUser.Login != username {
		errors.HandleError(w, errors.NewForbidden("Доступ запрещен."))
		return
	}

	baskets := a.baskets.Service.FindByUser(username)
	response, err := json.Marshal(baskets)

	if err != nil {
		panic(err)
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

// Remove
// @Summary Полностью удалить товар из корзины
// @ID remove-full-item-from-basket
// @Accept json
// @Produce json
// @Param username path string true "Логин пользователя"
// @Param itemCode query string true "Код товара"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /basket/{username} [delete]
// @Security ApiKeyAuth
func (a *API) Remove(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]
	itemCode := r.URL.Query().Get("itemCode")

	if len(username) == 0 {
		errors.HandleError(w, errors.NewBadRequest("Логин пользователя обязателен к заполнению."))
		return
	}

	if len(itemCode) == 0 {
		errors.HandleError(w, errors.NewBadRequest("Код товара обязателен к заполнению."))
		return
	}

	currentUser := common.GetJwtUser(r)
	if currentUser.Role != constants.ADMIN && currentUser.Login != username {
		errors.HandleError(w, errors.NewForbidden("Доступ запрещен."))
		return
	}

	removeErr := a.baskets.Service.Remove(username, itemCode)
	if removeErr != nil {
		errors.HandleError(w, removeErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// Pay
// @Summary Оплатить корзину товаров
// @ID pay-basket
// @Accept json
// @Produce json
// @Param username path string true "Логин пользователя"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /basket/{username} [post]
// @Security ApiKeyAuth
func (a *API) Pay(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if len(username) == 0 {
		errors.HandleError(w, errors.NewBadRequest("Логин пользователя обязателен к заполнению."))
		return
	}

	currentUser := common.GetJwtUser(r)
	if currentUser.Role != constants.ADMIN && currentUser.Login != username {
		errors.HandleError(w, errors.NewForbidden("Доступ запрещен."))
		return
	}

	removeErr := a.baskets.Service.Pay(username)
	if removeErr != nil {
		errors.HandleError(w, removeErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// GetHistoryByUser
// @Summary Получить историю оплаты товаров корзин
// @ID get-history-by-user
// @Accept json
// @Produce json
// @Param username path string true "Логин пользователя"
// @Success 200 {array} models.BasketHistoryView
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /basket/{username}/history [get]
// @Security ApiKeyAuth
func (a *API) GetHistoryByUser(w http.ResponseWriter, r *http.Request) {
	username := mux.Vars(r)["username"]

	if len(username) == 0 {
		errors.HandleError(w, errors.NewBadRequest("Логин пользователя обязателен к заполнению."))
		return
	}

	currentUser := common.GetJwtUser(r)
	if currentUser.Role != constants.ADMIN && currentUser.Login != username {
		errors.HandleError(w, errors.NewForbidden("Доступ запрещен."))
		return
	}

	response := a.baskets.Service.GetHistoryByUser(username)

	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}
