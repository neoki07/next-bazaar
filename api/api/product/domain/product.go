package domain

import "github.com/google/uuid"

type Product struct {
	ID   uuid.UUID
	Name string
}

func NewProduct(id uuid.UUID, name string) *Product {
	return &Product{
		ID:   id,
		Name: name,
	}
}
