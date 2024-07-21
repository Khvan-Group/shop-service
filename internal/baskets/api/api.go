package api

import (
	"github.com/Khvan-Group/common-library/middlewares"
	"github.com/gorilla/mux"
	"net/http"
	"shop-service/internal/baskets/service"
)

type API struct {
	baskets service.Baskets
}

func New() *API {
	return &API{
		baskets: *service.New(),
	}
}

func (a *API) AddRoutes(r *mux.Router) {
	r.Handle("", middlewares.AuthMiddleware(http.HandlerFunc(a.Save))).Methods(http.MethodPost)
	r.Handle("/{username}", middlewares.AuthMiddleware(http.HandlerFunc(a.FindByUser))).Methods(http.MethodGet)
	r.Handle("/{username}", middlewares.AuthMiddleware(http.HandlerFunc(a.Remove))).Methods(http.MethodDelete)
	r.Handle("/{username}", middlewares.AuthMiddleware(http.HandlerFunc(a.Pay))).Methods(http.MethodPost)
	r.Handle("/{username}/history", middlewares.AuthMiddleware(http.HandlerFunc(a.GetHistoryByUser))).Methods(http.MethodGet)
}
