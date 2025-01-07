package rest

import (
	"net/http"
	"valerii/crudbananas/internal/domain"
)

type BananaItem interface {
	Create(banana domain.Banana) (int, error)
	GetAll() ([]domain.Banana, error)
	GetById(id int) (domain.Banana, error)
	Update(id int, banana domain.BananaUpdate) error
	Delete(id int) error
}

type Handler struct {
	bananaService BananaItem
}

func NewHandler(bananaService BananaItem) *Handler{
	return &Handler{bananaService: bananaService}
}

func (h *Handler) InitRoutes() http.Handler{
	mux := http.NewServeMux()

	mux.HandleFunc("/api/items/", h.routeHandler)

	return loggingMiddleware(mux)
}