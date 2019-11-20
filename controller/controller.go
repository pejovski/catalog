package controller

import (
	emitter "github.com/pejovski/catalog/emitter/amqp"
	"github.com/pejovski/catalog/gateway/reviewing"
	"github.com/pejovski/catalog/model"
	"github.com/pejovski/catalog/repository"
	"github.com/sirupsen/logrus"
)

type Controller interface {
	GetProduct(id string) (*model.Product, error)
	GetProducts(category string) ([]*model.Product, error)
	CreateProduct(p *model.Product) (id string, err error)
	UpdateProduct(p *model.Product) error
	DeleteProduct(id string) error
	UpdateProductPrice(id string, price float32) error
	UpdateRating(id string) error
}

type controller struct {
	repository repository.Repository
	emitter    emitter.Emitter
	reviewing  reviewing.Gateway
}

func New(r repository.Repository, e emitter.Emitter, rev reviewing.Gateway) Controller {
	return controller{repository: r, emitter: e, reviewing: rev}
}

func (c controller) GetProduct(id string) (*model.Product, error) {
	p, err := c.repository.Get(id)
	if err != nil {
		logrus.Errorf("Failed to get product %s; Error: %s", id, err)
		return nil, err
	}

	return p, nil
}

func (c controller) GetProducts(category string) ([]*model.Product, error) {
	ps, err := c.repository.GetByCategory(category)
	if err != nil {
		logrus.Errorf("Failed to get products for category %s; Error: %s", category, err)
		return nil, err
	}

	return ps, nil
}

func (c controller) CreateProduct(p *model.Product) (id string, err error) {
	return c.repository.Create(p)
}

func (c controller) UpdateProduct(p *model.Product) (err error) {
	err = c.repository.Update(p)
	if err != nil {
		logrus.Errorf("Failed to update product %s; Error: %s", p.Id, err)
		return err
	}

	go c.emitter.ProductUpdated(p.Id)

	return err
}

func (c controller) UpdateProductPrice(id string, price float32) (err error) {
	err = c.repository.UpdatePrice(id, price)
	if err != nil {
		logrus.Errorf("Failed to update price of product %s; Error: %s", id, err)
		return err
	}

	go c.emitter.ProductPriceUpdated(id, price)

	return err
}

func (c controller) DeleteProduct(id string) (err error) {
	err = c.repository.Delete(id)
	if err != nil {
		logrus.Errorf("Failed to delete product %s; Error: %s", id, err)
		return
	}

	go c.emitter.ProductDeleted(id)

	return
}

func (c controller) UpdateRating(id string) error {
	rating, err := c.reviewing.Rating(id)
	if err != nil {
		logrus.Errorf("Failed to get rating for product %s, Error: %s", id, err)
		return err
	}

	if err = c.repository.UpdateRating(id, rating); err != nil {
		logrus.Errorf("Failed to update rating for product %s, Error: %s", id, err)
		return err
	}

	return nil
}
