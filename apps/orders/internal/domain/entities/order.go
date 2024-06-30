package entities

import (
	"time"
)

type Order struct {
	ID        int
	Name      string
	Price     float64
	Quantity  int
	Status    string
	CreatedAt time.Time
	UpdatedAt time.Time
}
