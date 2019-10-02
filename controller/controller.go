package controller

import (
	"github.com/pejovski/catalog/domain"
	"github.com/sirupsen/logrus"
)

type Catalog struct {
	repository domain.CatalogRepository
	emitter    domain.AmqpEmitter
	reviewing  domain.ReviewingGateway
}

func NewCatalog(r domain.CatalogRepository, e domain.AmqpEmitter, rev domain.ReviewingGateway) Catalog {
	return Catalog{repository: r, emitter: e, reviewing: rev}
}

func (c Catalog) GetProduct(id string) (*domain.Product, error) {
	p, err := c.repository.Get(id)
	if err != nil {
		logrus.Errorf("Failed to get product %s; Error: %s", id, err)
		return nil, err
	}

	return p, nil
}

func (c Catalog) GetProducts(category string) ([]*domain.Product, error) {
	ps, err := c.repository.GetByCategory(category)
	if err != nil {
		logrus.Errorf("Failed to get products for category %s; Error: %s", category, err)
		return nil, err
	}

	return ps, nil
}

func (c Catalog) CreateProduct(p *domain.Product) (id string, err error) {
	return c.repository.Create(p)
}

func (c Catalog) UpdateProduct(p *domain.Product) (err error) {
	err = c.repository.Update(p)
	if err != nil {
		logrus.Errorf("Failed to update product %s; Error: %s", p.Id, err)
		return err
	}

	go c.emitter.ProductUpdated(p.Id)

	return err
}

func (c Catalog) UpdateProductPrice(id string, price float32) (err error) {
	err = c.repository.UpdatePrice(id, price)
	if err != nil {
		logrus.Errorf("Failed to update price of product %s; Error: %s", id, err)
		return err
	}

	go c.emitter.ProductPriceUpdated(id, price)

	return err
}

func (c Catalog) DeleteProduct(id string) (err error) {
	err = c.repository.Delete(id)
	if err != nil {
		logrus.Errorf("Failed to delete product %s; Error: %s", id, err)
		return
	}

	go c.emitter.ProductDeleted(id)

	return
}

func (c Catalog) UpdateRating(id string) error {
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
