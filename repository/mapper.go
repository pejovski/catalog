package repository

import "github.com/pejovski/catalog/domain"

func mapHitToProduct(h *Hit) *domain.Product {
	s := h.Source
	return &domain.Product{Id: h.Id, Name: s.Name, Brand: s.Brand, Price: s.Price, Category: s.Category, Image: s.Image}
}

func mapProductToDocument(p *domain.Product) *Document {
	return &Document{Name: p.Name, Brand: p.Brand, Price: p.Price, Category: p.Category, Image: p.Image}
}
