package reviewing

import "github.com/pejovski/catalog/domain"

func (g Gateway) mapRatingToDomainRating(r Rating) *domain.Rating {
	return &domain.Rating{
		Stars:     r.Stars,
		Customers: r.Customers,
	}
}
