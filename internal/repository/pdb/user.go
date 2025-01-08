package pdb

import (
	"context"
	"database/sql"
	"fmt"
	"valerii/crudbananas/internal/domain"
	"valerii/crudbananas/pkg/database"
)

type Users struct {
	db *sql.DB
}

func NewUsers(db *sql.DB) *Users{
	return &Users{db: db}
}

func (r *Users) Create(ctx context.Context, user domain.User) error {
	query := fmt.Sprintf("INSERT INTO %s (name, email, password, registered_at) values ($1, $2, $3, $4)", database.UsersTable)
	_, err := r.db.Exec(query, user.Name, user.Email, user.Password, user.RegisteredAt)

	return err
}

func (r *Users) GetByCredentials(ctx context.Context, email, password string) (domain.User, error) {
	var user domain.User
	query := fmt.Sprintf("SELECT id, name, email, registered_at FROM %s WHERE email=$1 AND password=$2", database.UsersTable)
	err := r.db.QueryRow(query, email, password).
		Scan(&user.ID, &user.Name, &user.Email, &user.RegisteredAt)

	return user, err
}