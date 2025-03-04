package rest

import (
	"context"
	"net/http"
	"valerii/crudbananas/internal/domain"
)

type User interface {
	SignUp(ctx context.Context, inp domain.SignUpInput) error
	SignIn(ctx context.Context, inp domain.SignInInput) (string, string, error)
	ParseToken(ctx context.Context, token string) (int64, error)
	RefreshTokens(ctx context.Context, refreshToken string) (string, string, error)
}

type BananaItem interface {
	Create(ctx context.Context, banana domain.Banana) (int, error)
	GetAll(ctx context.Context, ) ([]domain.Banana, error)
	GetById(ctx context.Context, id int) (domain.Banana, error)
	Update(ctx context.Context, id int, banana domain.BananaUpdate) error
	Delete(ctx context.Context, id int) error
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

	//Узел на аунтитефикацию и регистрацию
	mux.HandleFunc("/api/auth/", h.routeAuth)
	mux.HandleFunc("/api/auth/refresh", h.refresh)

	//Защищенный узел на все методы взаимодействия с объектами
	//Требуется аунтетификация (токен)
	protectedMux := http.NewServeMux()
	protectedMux.HandleFunc("/api/items/", h.routeBananas)

	protectedRoutes := h.authMiddleware(protectedMux)

	mux.Handle("/api/items/", protectedRoutes)

	return loggingMiddleware(mux)
}