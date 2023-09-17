package entity

import (
	"awesomeProject/pkg/entity"
	"errors"
	"time"
)

var (
	ErrIDIsRequired    = errors.New("id is required")
	ErrInvalidID       = errors.New("invalid id")
	ErrNameIsRequired  = errors.New("name is required")
	ErrPriceisRequired = errors.New("price is required")
	ErrInvalidPrice    = errors.New("invalid price")
)

type Product struct {
	ID       entity.ID `json:"id"`
	Name     string    `json:"name"`
	Price    float64   `json:"price"`
	CreateAt time.Time `gorm:"column:create_at;type:timestamp"`
}

func NewProduct(name string, price float64) (*Product, error) {
	product := &Product{
		ID:       entity.NewId(),
		Name:     name,
		Price:    price,
		CreateAt: time.Now(),
	}
	err := product.Validate()
	if err != nil {
		return nil, err
	}
	return product, nil
}

func (p *Product) Validate() error {
	if p.ID.String() == "" {
		return ErrIDIsRequired
	}
	if _, err := entity.ParseID(p.ID.String()); err != nil {
		return ErrInvalidID
	}
	if p.Name == "" {
		return ErrNameIsRequired
	}
	if p.Price == 0 {
		return ErrPriceisRequired
	}
	if p.Price < 0 {
		return ErrInvalidPrice
	}
	return nil
}
