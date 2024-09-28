package interfaces

import models "BChat/Domain/Models"


type OAuthUseCase interface{
	GetOAuthURL() (string, error)
	ExchangeCode(code string) (*models.OAuthToken, error)
	FetchUserInfo(token *models.OAuthToken) (*models.OAuthUser, error)
}