package service

import (
	"context"
	"valerii/crudbananas/internal/domain"
)

type BananaRepository interface {
	Create(ctx context.Context, banana domain.Banana) (int, error)
	GetAll(ctx context.Context) ([]domain.Banana, error)
	GetById(ctx context.Context, id int) (domain.Banana, error)
	Update(ctx context.Context, id int, banana domain.BananaUpdate) error
	Delete(ctx context.Context, id int) error
}

type Bananas struct {
	repo        BananaRepository
	auditClient AuditClient
}

func NewBananas(repo BananaRepository, auditClient AuditClient) *Bananas {
	return &Bananas{
		repo:        repo,
		auditClient: auditClient,
	}
}

func (p *Bananas) Create(ctx context.Context, banana domain.Banana) (int, error) {
	return p.repo.Create(ctx, banana)
}

func (p *Bananas) GetAll(ctx context.Context) ([]domain.Banana, error) {
	return p.repo.GetAll(ctx)
}

func (p *Bananas) GetById(ctx context.Context, id int) (domain.Banana, error) {
	return p.repo.GetById(ctx, id)
}

func (p *Bananas) Update(ctx context.Context, id int, banana domain.BananaUpdate) error {
	return p.repo.Update(ctx, id, banana)
}

func (p *Bananas) Delete(ctx context.Context, id int) error {
	return p.repo.Delete(ctx, id)
}
