package api

import (
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	"net/http"
	"shop-service/internal/items/service"
)

type API struct {
	items service.Items
}

func New() *API {
	return &API{
		items: *service.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Create), constants.ADMIN)).Methods(http.MethodPost)
	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Update), constants.ADMIN)).Methods(http.MethodPut)
	r.Handle("/{code}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete), constants.ADMIN)).Methods(http.MethodDelete)

	r.HandleFunc("", a.FindAll).Methods(http.MethodGet)
	r.HandleFunc("/{code}", a.FindByCode).Methods(http.MethodGet)
}
