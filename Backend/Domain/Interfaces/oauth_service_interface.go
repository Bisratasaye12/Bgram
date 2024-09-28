package interfaces

import models "BChat/Domain/Models"

// OAuthService defines the interface for handling OAuth operations
type OAuthService interface {
	GetAuthURL() (string, error)               // Get the OAuth URL for user redirection
	ExchangeCodeForToken(code string) (*models.OAuthToken, error) // Exchange authorization code for a token
	GetUserInfo(token *models.OAuthToken) (*models.OAuthUser, error)  // Get user info using the token
}