package interfaces

import (
	models "BChat/Domain/Models"
	"time"
)

type JWTServiceInterface interface {
    GenerateToken(user *models.User, expiration time.Duration, UrlID string) (string, error)
    GenerateRefreshToken(user *models.User, expiration time.Duration) (string, error)
    ValidateToken(token string) (*models.CustomClaims, error)
    HashPassword(password string) (string, error)
    CheckPasswordHash(password, hash string) bool
}
