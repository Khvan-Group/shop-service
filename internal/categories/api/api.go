package api

import (
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	"net/http"
	"shop-service/internal/categories/service"
)

type API struct {
	categories service.Categories
}

func New() *API {
	return &API{
		categories: *service.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Save), constants.ADMIN)).Methods(http.MethodPost)
	r.Handle("/{code}", middlewares.AuthMiddleware(http.HandlerFunc(a.Delete), constants.ADMIN)).Methods(http.MethodDelete)

	r.HandleFunc("", a.FindAll).Methods(http.MethodGet)
}
