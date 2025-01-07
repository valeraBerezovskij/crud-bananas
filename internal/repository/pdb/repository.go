package pdb

import (
	"database/sql"
	"valerii/crudbananas/internal/domain"
)

type BananaItem interface {
	Create(banana domain.Banana) (int, error)
	GetAll() ([]domain.Banana, error)
	GetById(id int) (domain.Banana, error)
	Update(id int, banana domain.BananaUpdate) error
	Delete(id int) error
}

type Repository struct{
	BananaItem
}

func NewRepository(db *sql.DB) *Repository{
	return &Repository{
		BananaItem: NewItemPostgres(db),
	}
}