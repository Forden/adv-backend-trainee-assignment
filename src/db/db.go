package db

import "adv-backend-trainee-assignment/src/models"

type DatabaseConnection interface {
	NewAd(adData models.CreatingAd) (string, error)
	SelectAd(adID string) (*models.DbAd, error)
	GetAllAds(sortBy string, sortOrder string, page int, perPage int) ([]*models.DbAd, error)
	Close() error
}
