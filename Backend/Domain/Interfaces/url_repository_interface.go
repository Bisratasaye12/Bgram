package interfaces

import models "BChat/Domain/Models"

type VerificationURLRepositoryInterface interface {
	SaveUrl(url *models.VerificationURL) (*models.VerificationURL, error)
	GetUrlByID(urlID string) (*models.VerificationURL, error)
	DeleteUrlByID(urlID string) error
}
