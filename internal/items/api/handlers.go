package api

import (
	"encoding/json"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	commonModels "github.com/Khvan-Group/common-library/models"
	"github.com/gorilla/mux"
	"golang.org/x/tools/container/intsets"
	"io"
	"net/http"
	"shop-service/internal/items/models"
	"strconv"
	"strings"
)

// Create
// @Summary Создать товар
// @ID create-item
// @Accept json
// @Produce json
// @Param input body models.ItemDto true "Информация о создаваемом товаре"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /items [post]
// @Security ApiKeyAuth
func (a *API) Create(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	var input models.ItemDto

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	createErr := a.items.Service.Create(input)
	if createErr != nil {
		errors.HandleError(w, createErr)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

// Update
// @Summary Обновить товар
// @ID update-item
// @Accept json
// @Produce json
// @Param input body models.ItemDto true "Информация об изменяемом товаре"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /items [put]
// @Security ApiKeyAuth
func (a *API) Update(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	var input models.ItemDto

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	updateErr := a.items.Service.Update(input)
	if updateErr != nil {
		errors.HandleError(w, updateErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FindAll
// @Summary Получить список товаров
// @ID find-all-items
// @Accept  json
// @Produce  json
// @Param page query int false "Номер страницы"
// @Param size query int false "Количество элементов за раз"
// @Param sortFields query string false "Сортировка по полям"
// @Param name query string false "Название товара"
// @Param code query string false "Код товара"
// @Param category query string false "Код категории товара"
// @Success 200 {array} commonModels.Page
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /items [get]
func (a *API) FindAll(w http.ResponseWriter, r *http.Request) {
	var input models.ItemSearch
	queryParams := r.URL.Query()
	page, err := strconv.Atoi(queryParams.Get("page"))
	if err != nil {
		page = 0
	}

	size, err := strconv.Atoi(queryParams.Get("size"))
	if err != nil {
		size = intsets.MaxInt
	}

	sortFields := queryParams["sortFields"]
	for _, field := range sortFields {
		parts := strings.Split(field, ":")

		if len(parts) == 2 {
			input.SortBy = append(input.SortBy, commonModels.SortField{SortBy: parts[0], Direction: parts[1]})
		}
	}

	name := queryParams.Get("name")
	if len(name) > 0 {
		input.Name = &name
	}

	code := queryParams.Get("code")
	if len(code) > 0 {
		input.Code = &code
	}
	category := queryParams.Get("category")
	if len(category) > 0 {
		input.Category = &category
	}

	input.Page = page
	input.Size = size
	response := a.items.Service.FindAll(input)

	data, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}

	w.Write(data)
	w.WriteHeader(http.StatusOK)
}

// FindByCode
// @Summary Получить товар по коду
// @ID find-item-by-code
// @Accept json
// @Produce json
// @Param code path string true "Код товара"
// @Success 200 {object} models.Item
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /items/{code} [get]
func (a *API) FindByCode(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	item, err := a.items.Service.FindByCode(code)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	response, marshalErr := json.Marshal(item)
	if marshalErr != nil {
		panic(err)
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

// Delete
// @Summary Удалить товар
// @ID delete-item
// @Accept json
// @Produce json
// @Param code path string true "Код товара"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /items/{code} [delete]
// @Security ApiKeyAuth
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	err := a.items.Service.Delete(code)
	if err != nil {
		errors.HandleError(w, err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
