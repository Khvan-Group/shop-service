package api

import (
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
	_ "shop-service/docs"
	basketApi "shop-service/internal/baskets/api"
	categoryApi "shop-service/internal/categories/api"
	itemsApi "shop-service/internal/items/api"
)

type API struct {
	itemsApi    itemsApi.API
	categoryApi categoryApi.API
	basketApi   basketApi.API
}

func New() *API {
	return &API{
		itemsApi:    *itemsApi.New(),
		categoryApi: *categoryApi.New(),
		basketApi:   *basketApi.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	r = r.PathPrefix("/api/v1").Subrouter()

	itemsRouter := r.PathPrefix("/items").Subrouter()
	a.itemsApi.AddRoutes(itemsRouter)

	categoryRouter := r.PathPrefix("/categories").Subrouter()
	a.categoryApi.AddRoutes(categoryRouter)

	basketRouter := r.PathPrefix("/basket").Subrouter()
	a.basketApi.AddRoutes(basketRouter)

	r.PathPrefix("/swagger").Handler(httpSwagger.WrapHandler)
}
