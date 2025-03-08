package domain

import "time"

// TODO Validator
type Banana struct {
	ID        int       `json:"id"`
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Length    float64   `json:"length"`
	CreatedAt time.Time `json:"created_at"`
}

type BananaUpdate struct {
	Name      string    `json:"name"`
	Color     string    `json:"color"`
	Length    float64   `json:"length"`
}