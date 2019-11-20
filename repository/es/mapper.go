package es

import (
	"github.com/pejovski/catalog/model"
)

func mapHitToProduct(h *Hit) *model.Product {
	s := h.Source
	return &model.Product{Id: h.Id, Name: s.Name, Brand: s.Brand, Price: s.Price, Category: s.Category, Image: s.Image}
}

func mapProductToDocument(p *model.Product) *Document {
	return &Document{Name: p.Name, Brand: p.Brand, Price: p.Price, Category: p.Category, Image: p.Image}
}
