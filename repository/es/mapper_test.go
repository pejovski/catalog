package es

import (
	"testing"
)

func TestMapHitToProduct(t *testing.T) {

	hit := &Hit{
		Id: "111",
		Source: Document{
			Name:     "Galaxy",
			Brand:    "Samsung",
			Price:    800,
			Category: "555",
			Image:    "galaxy.jpg",
		},
	}

	p := mapHitToProduct(hit)

	if p.Id != hit.Id {
		t.Error("Expected ids to be equal")
	}
}
