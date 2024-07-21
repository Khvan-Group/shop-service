package api

import (
	"encoding/json"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/gorilla/context"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"shop-service/internal/categories/models"
	"shop-service/internal/common"
)

// Save
// @Summary Создать/Обновить категорию товаров
// @ID save-category
// @Accept json
// @Produce json
// @Param input body models.Category true "Информация о категории"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /categories [post]
// @Security ApiKeyAuth
func (a *API) Save(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)
	var input models.Category

	data, err := io.ReadAll(r.Body)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(data, &input); err != nil {
		panic(err)
	}

	if customErr := a.categories.Service.Save(input); customErr != nil {
		errors.HandleError(w, customErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

// FindAll
// @Summary Получить список категорий товаров
// @ID find-all-categories
// @Accept json
// @Produce json
// @Success 200 {array} models.Category
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /categories [get]
func (a *API) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Set(constants.CONTENT_TYPE, constants.APPLICATION_JSON)

	data := a.categories.Service.FindAll()
	response, err := json.Marshal(data)
	if err != nil {
		panic(err)
	}

	w.Write(response)
	w.WriteHeader(http.StatusOK)
}

// Delete
// @Summary Удалить категорию товаров
// @ID delete-category
// @Accept json
// @Produce json
// @Param code path string true "Информация о создаваемой категории"
// @Success 200
// @Failure 404 {object} string
// @Failure 400 {object} string
// @Router /categories/{code} [delete]
// @Security ApiKeyAuth
func (a *API) Delete(w http.ResponseWriter, r *http.Request) {
	code := mux.Vars(r)["code"]

	if customErr := a.categories.Service.Delete(code); customErr != nil {
		errors.HandleError(w, customErr)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func getJwtUser(r *http.Request) common.JwtUser {
	return common.JwtUser{
		Login: utils.ToString(context.Get(r, "login")),
		Role:  utils.ToString(context.Get(r, "role")),
	}
}
