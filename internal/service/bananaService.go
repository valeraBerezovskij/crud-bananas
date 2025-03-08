package service

import (
	"context"
	"github.com/sirupsen/logrus"
	audit "github.com/valeraBerezovskij/logger-mongo/pkg/domain"
	"time"
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
	//Достаем ID
	bananaID, err := p.repo.Create(ctx, banana)
	if err != nil {
		return 0, err
	}

	//Логирование
	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    "CREATE",
		Entity:    "BANANA",
		EntityID:  int64(bananaID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "Create",
		}).Error(err)
	}

	return bananaID, nil
}

func (p *Bananas) GetAll(ctx context.Context) ([]domain.Banana, error) {
	bananas, err := p.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	//Логирование
	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    "GET",
		Entity:    "BANANA",
		EntityID:  0,
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "GetAll",
		}).Error(err)
	}

	return bananas, nil
}

func (p *Bananas) GetById(ctx context.Context, id int) (domain.Banana, error) {
	banana, err := p.repo.GetById(ctx, id)
	if err != nil {
		return domain.Banana{}, err
	}

	//Логирование
	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    "GET",
		Entity:    "BANANA",
		EntityID:  int64(banana.ID),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "GetById",
		}).Error(err)
	}

	return p.repo.GetById(ctx, id)
}

func (p *Bananas) Update(ctx context.Context, id int, banana domain.BananaUpdate) error {
	err := p.repo.Update(ctx, id, banana)
	if err != nil {
		return err
	}

	//Логирование
	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    "UPDATE",
		Entity:    "BANANA",
		EntityID:  int64(id),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "Update",
		}).Error(err)
	}

	return nil
}

func (p *Bananas) Delete(ctx context.Context, id int) error {
	err := p.repo.Delete(ctx, id)
	if err != nil {
		return err
	}

	//Логирование
	if err := p.auditClient.SendLogRequest(ctx, audit.LogItem{
		Action:    "DELETE",
		Entity:    "BANANA",
		EntityID:  int64(id),
		Timestamp: time.Now(),
	}); err != nil {
		logrus.WithFields(logrus.Fields{
			"handler": "Delete",
		}).Error(err)
	}

	return nil
}
