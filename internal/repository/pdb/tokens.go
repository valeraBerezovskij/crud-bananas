package pdb

import (
	"context"
	"database/sql"
	"valerii/crudbananas/internal/domain"
)

type Tokens struct {
	db *sql.DB
}

func NewTokens(repo *sql.DB) *Tokens {
	return &Tokens{db: repo}
}

func (t *Tokens) Create(ctx context.Context, token domain.RefreshSession) error {
	_, err := t.db.Exec("INSERT INTO tokens (user_id, token, expires_at) VALUES ($1, $2, $3)",
		token.UserID, token.Token, token.ExpiresAt)

	return err
}

func (t *Tokens) Get(ctx context.Context, token string) (domain.RefreshSession, error) {
	var inp domain.RefreshSession
	err := t.db.QueryRow("SELECT id, user_id, token, expires_at FROM tokens WHERE token=$1",
		token).Scan(&inp.ID, &inp.UserID, &inp.Token, &inp.ExpiresAt)
	if err != nil {
		return inp, err
	}

	_, err = t.db.Exec("DELETE FROM refresh_tokens WHERE user_id=$1", inp.UserID)
	return inp, err
}
