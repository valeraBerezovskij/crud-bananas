package pdb

import (
	"database/sql"
	"fmt"
	"valerii/crudbananas/internal/domain"
	"valerii/crudbananas/pkg/database"
)

type ItemPostgres struct {
	db *sql.DB
}

func NewItemPostgres(db *sql.DB) *ItemPostgres {
	return &ItemPostgres{
		db: db,
	}
}

func (p *ItemPostgres) Create(banana domain.Banana) (int, error) {
	var id int
	query := fmt.Sprintf("insert into %s(name, color, length) values($1, $2, $3) RETURNING id", database.BananaTable)
	row := p.db.QueryRow(query, banana.Name, banana.Color, banana.Length)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *ItemPostgres) GetAll() ([]domain.Banana, error) {
	bananas := make([]domain.Banana, 0)
	query := fmt.Sprintf("select * from %s", database.BananaTable)
	rows, err := p.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var banana domain.Banana
		if err := rows.Scan(&banana.ID, &banana.Name, &banana.Color, &banana.Length, &banana.CreatedAt); err != nil {
			return nil, err
		}
		bananas = append(bananas, banana)
	}

	return bananas, nil
}

func (p *ItemPostgres) GetById(id int) (domain.Banana, error) {
	var banana domain.Banana
	query := fmt.Sprintf("select * from %s where id = $1", database.BananaTable)
	row := p.db.QueryRow(query, id)
	if err := row.Scan(&banana.ID, &banana.Name, &banana.Color, &banana.Length, &banana.CreatedAt); err != nil {
		return domain.Banana{}, nil
	}

	return banana, nil
}

func (p *ItemPostgres) Update(id int, banana domain.BananaUpdate) error {
	query := fmt.Sprintf("UPDATE %s SET %s = $1, %s = $2, %s = $3 WHERE id = $4",
		database.BananaTable, "name", "color", "length")
	result, err := p.db.Exec(query, &banana.Name, banana.Color, banana.Length, id)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return fmt.Errorf("record not found")
	}

	return nil
}

func (p *ItemPostgres) Delete(id int) error {
	query := fmt.Sprintf("DELETE FROM %s WHERE id = $1", database.BananaTable)
	result, err := p.db.Exec(query, id)
	if err != nil {
		fmt.Printf("Error executing query: %v\n", err)
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		fmt.Printf("Error getting rows affected: %v\n", err)
		return err
	}

	if rowsAffected == 0 {
		fmt.Printf("No rows affected. Record with id=%d not found.\n", id)
		return fmt.Errorf("record not found")
	}

	return nil
}
