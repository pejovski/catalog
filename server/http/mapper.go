package http

import "github.com/pejovski/catalog/domain"

func mapDomainProductToProduct(dp *domain.Product) *Product {
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

func mapDomainProductsToProducts(dps []*domain.Product) []*Product {
	ps := []*Product{}
	for _, dp := range dps {
		ps = append(ps, mapDomainProductToProduct(dp))
	}
	return ps
}
