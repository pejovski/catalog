package reviewing

import (
	"github.com/pejovski/catalog/model"
)

func (g gateway) mapRatingToDomainRating(r Rating) *model.Rating {
	return &model.Rating{
		Stars:     r.Stars,
		Customers: r.Customers,
	}
}
