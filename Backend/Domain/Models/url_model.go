package models

import (
	"github.com/dgrijalva/jwt-go"
)

type VerificationURL struct {
	ID    string `json:"_id", bson:"_id"`
	UrlID string `json:"url_id", bson:"url_id"`
	URL   string `json:"url", bson:"url"`
	jwt.StandardClaims
}
