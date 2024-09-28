package usecases

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
)



type OAuthUseCase struct {
	oauthService interfaces.OAuthService
}

// NewOAuthUseCase creates a new OAuthUseCase
func NewOAuthUseCase(oauthService interfaces.OAuthService) interfaces.OAuthUseCase {
	return &OAuthUseCase{oauthService: oauthService}
}

// GetOAuthURL retrieves the URL where the user should be redirected for OAuth login
func (uc *OAuthUseCase) GetOAuthURL() (string, error) {
	return uc.oauthService.GetAuthURL()
}

// ExchangeCode exchanges the authorization code for an OAuth token
func (uc *OAuthUseCase) ExchangeCode(code string) (*models.OAuthToken, error) {
	return uc.oauthService.ExchangeCodeForToken(code)
}

// FetchUserInfo retrieves user information using the OAuth token
func (uc *OAuthUseCase) FetchUserInfo(token *models.OAuthToken) (*models.OAuthUser, error) {
	return uc.oauthService.GetUserInfo(token)
}