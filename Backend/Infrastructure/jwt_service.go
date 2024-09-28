package infrastructure

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
)

type JWTService struct {
	secretKey string
}

// NewJWTService creates a new JWT service
func NewJWTService(secretKey string) interfaces.JWTServiceInterface {
	return &JWTService{
		secretKey: secretKey,
	}
}

// GenerateToken generates a JWT token with a given expiration time
func (s *JWTService) GenerateToken(user *models.User, expiration time.Duration, UrlID string) (string, error) {
	claims := &models.CustomClaims{
		UserID:    user.ID,
		UserEmail: user.Email,
		Role:      user.Role,
		Username:  user.Username,
		UrlID:     UrlID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expiration)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken generates a JWT refresh token with a given expiration time
func (s *JWTService) GenerateRefreshToken(user *models.User, expiration time.Duration) (string, error) {
	claims := &models.CustomClaims{
		UserID:    user.ID,
		UserEmail: user.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Local().Add(time.Hour * time.Duration(expiration)).Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ValidateToken validates a given JWT token and returns the claims if valid
func (s *JWTService) ValidateToken(tokenString string) (*models.CustomClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &models.CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.secretKey), nil
	})

	if err != nil {
		return nil, fmt.Errorf("invalid token: %v", err)
	}

	if claims, ok := token.Claims.(*models.CustomClaims); ok {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

// HashPassword hashes a given password using bcrypt
func (s *JWTService) HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

// CheckPasswordHash compares a password with a bcrypt hash
func (s *JWTService) CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
