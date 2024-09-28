package infrastructure

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// OAuthService struct holds the OAuth2 configuration
type OAuthService struct {
	config *oauth2.Config
	env    *models.Env
}

// NewOAuthService creates a new instance of OAuthService
func NewOAuthService(env *models.Env) interfaces.OAuthService {
	clientID := env.G_CLIENT_ID
	clientSecret := env.G_CLIENT_SECRET
	redirectURL := env.G_REDIRECT_URL

	log.Println(clientID, "id", clientSecret, "sec", redirectURL, "EENVV from oauth service")

	config := &oauth2.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  redirectURL,
		Scopes:       []string{"openid", "email", "profile"},
		Endpoint:     google.Endpoint,
	}

	return &OAuthService{
		config: config,
		env:    env,
	}
}

// GetAuthURL returns the URL to redirect the user for authentication
func (s *OAuthService) GetAuthURL() (string, error) {
	authURL := s.config.AuthCodeURL("state", oauth2.AccessTypeOffline)
	log.Println(s.config, authURL, "from oauth service")
	return authURL, nil
}

// ExchangeCodeForToken exchanges the authorization code for an OAuth token
func (s *OAuthService) ExchangeCodeForToken(code string) (*models.OAuthToken, error) {
	token, err := s.config.Exchange(context.Background(), code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange token: %v", err)
	}

	return &models.OAuthToken{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		ExpiresIn:    int(token.Expiry.Sub(token.Expiry).Seconds()),
		TokenType:    token.TokenType,
	}, nil
}

// GetUserInfo retrieves the user's profile information
func (s *OAuthService) GetUserInfo(token *models.OAuthToken) (*models.OAuthUser, error) {
	client := s.config.Client(context.Background(), &oauth2.Token{
		AccessToken: token.AccessToken,
		TokenType:   token.TokenType,
	})

	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %v", err)
	}
	defer resp.Body.Close()

	var user models.OAuthUser
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return nil, errors.New("failed to parse user info")
	}

	return &user, nil
}
