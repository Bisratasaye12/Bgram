package models

import "github.com/dgrijalva/jwt-go"

type CustomClaims struct {
	UserID    string `json:"user_id"`
	UrlID     string `json:"url_id"`
	UserEmail string `json:"user_email"`
	Role      string `json:"role"`
	Username  string `json:"username"`
	jwt.StandardClaims
}
