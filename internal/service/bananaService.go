package service

import (
	"valerii/crudbananas/internal/domain"
)

type BananaRepository interface {
	Create(banana domain.Banana) (int, error)
	GetAll() ([]domain.Banana, error)
	GetById(id int) (domain.Banana, error)
	Update(id int, banana domain.BananaUpdate) error
	Delete(id int) error
}

type Bananas struct {
	repo BananaRepository
}

func NewBananas(repo BananaRepository) *Bananas {
	return &Bananas{
		repo: repo,
	}
}

func (p *Bananas) Create(banana domain.Banana) (int, error) {
	return p.repo.Create(banana)
}

func (p *Bananas) GetAll() ([]domain.Banana, error) {
	return p.repo.GetAll()
}

func (p *Bananas) GetById(id int) (domain.Banana, error) {
	return p.repo.GetById(id)
}

func (p *Bananas) Update(id int, banana domain.BananaUpdate) error {
	return p.repo.Update(id, banana)
}

func (p *Bananas) Delete(id int) error {
	return p.repo.Delete(id)
}