package service

import (
	"context"
	"valerii/crudbananas/internal/domain"
	"time"
	//"valerii/crudbananas/internal/repository/pdb"
)

type PasswordHasher interface {
	Hash(password string) (string, error)
}

type UsersRepository interface {
	Create(ctx context.Context, user domain.User) error
	GetByCredentials(ctx context.Context, email, password string) (domain.User, error)
}

type Users struct {
	repo UsersRepository
	hasher PasswordHasher
}

func NewUsers(repo UsersRepository, hasher PasswordHasher) *Users{
	return &Users{
		repo: repo,
		hasher: hasher,
	}
}

func (s *Users) SignUp(ctx context.Context, inp domain.SignUpInput) error {
	//Хешируем пароль
	password, err := s.hasher.Hash(inp.Password)
	if err != nil {
		return err
	}

	//Создаем объект структуры domain.User
	user := domain.User{
		Name:         inp.Name,
		Email:        inp.Email,
		Password:     password,
		RegisteredAt: time.Now(),
	}

	//Передаем на уровень репозитория
	return s.repo.Create(ctx, user)
}