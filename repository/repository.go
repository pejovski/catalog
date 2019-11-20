package repository

import "github.com/pejovski/catalog/model"

type Repository interface {
	Get(id string) (*model.Product, error)
	Create(p *model.Product) (id string, err error)
	Update(p *model.Product) error
	Delete(id string) error
	GetByCategory(category string) ([]*model.Product, error)
	UpdatePrice(id string, price float32) error
	UpdateRating(id string, r *model.Rating) error
}
