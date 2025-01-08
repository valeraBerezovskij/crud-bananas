package rest

import (
	"context"
	"net/http"
	"valerii/crudbananas/internal/domain"
)

type User interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	//SignIn(ctx context.Context, inp domain.SignInInput) (string, error)
	//ParseToken(ctx context.Context, token string) (int64, error)
}

type BananaItem interface {
	Create(banana domain.Banana) (int, error)
	GetAll() ([]domain.Banana, error)
	GetById(id int) (domain.Banana, error)
	Update(id int, banana domain.BananaUpdate) error
	Delete(id int) error
}

type Handler struct {
	bananaService BananaItem
	userService   User
}

func NewHandler(bananaService BananaItem, userService User) *Handler {
	return &Handler{
		bananaService: bananaService,
		userService:   userService,
	}
}

func (h *Handler) InitRoutes() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/api/auth/sign-up", h.SignUp) //POST
	//mux.HandleFunc("/api/auth/sign-up", h.SignIn) //GET

	mux.HandleFunc("/api/items/", h.routeHandler)

	return loggingMiddleware(mux)
}
