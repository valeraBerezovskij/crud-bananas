package service

import (
	"valerii/crudbananas/internal/domain"
	"valerii/crudbananas/internal/repository/pdb"
)

type BananaService struct {
	repo pdb.BananaItem
}

func NewBananaService(repo pdb.BananaItem) *BananaService{
	return &BananaService{repo: repo}
}

func (p *BananaService) Create(banana domain.Banana) (int, error) {
	return p.repo.Create(banana)
}

func (p *BananaService) GetAll() ([]domain.Banana, error) {
	return p.repo.GetAll()
}

func (p *BananaService) GetById(id int) (domain.Banana, error) {
	return p.repo.GetById(id)
}

func (p *BananaService) Update(id int, banana domain.BananaUpdate) error {
	return p.repo.Update(id, banana)
}

func (p *BananaService) Delete(id int) error {
	return p.repo.Delete(id)
}