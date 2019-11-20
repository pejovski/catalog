package api

import (
	"github.com/pejovski/catalog/model"
)

type Mapper interface {
	mapDomainProductToProduct(dp *model.Product) *Product
	mapDomainProductsToProducts(dps []*model.Product) []*Product
}

type mapper struct {
}

func newMapper() Mapper {
	return mapper{}
}

func (m mapper) mapDomainProductToProduct(dp *model.Product) *Product {
	return &Product{
		Id:       dp.Id,
		Name:     dp.Name,
		Brand:    dp.Brand,
		Price:    dp.Price,
		Category: dp.Category,
		Image:    dp.Image,
		Rating: Rating{
			Stars:     dp.Stars,
			Customers: dp.Customers,
		},
	}
}

func (m mapper) mapDomainProductsToProducts(dps []*model.Product) []*Product {
	ps := []*Product{}
	for _, dp := range dps {
		ps = append(ps, m.mapDomainProductToProduct(dp))
	}
	return ps
}
