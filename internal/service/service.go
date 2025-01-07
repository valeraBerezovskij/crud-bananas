package service

import (
	"valerii/crudbananas/internal/domain"
	"valerii/crudbananas/internal/repository/pdb"
)

type BananaItem interface {
	Create(banana domain.Banana) (int, error)
	GetAll() ([]domain.Banana, error)
	GetById(id int) (domain.Banana, error)
	Update(id int, banana domain.BananaUpdate) error
	Delete(id int) error
}

type Service struct {
	BananaItem
}

func NewService(repos *pdb.Repository) *Service {
	return &Service{
		BananaItem: NewBananaService(repos.BananaItem),
	}
}
