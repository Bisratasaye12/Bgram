package controllers

import (
	interfaces "BChat/Domain/Interfaces"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type OAuthController struct {
	oauthUseCase interfaces.OAuthUseCase
}

// NewOAuthController creates a new OAuthController
func NewOAuthController(oauthUseCase interfaces.OAuthUseCase) interfaces.OAuthControllerInterface {
	return &OAuthController{oauthUseCase: oauthUseCase}
}

// RedirectToProvider redirects the user to the OAuth provider's authentication page
func (oc *OAuthController) RedirectToProvider(c *gin.Context) {
	authURL, err := oc.oauthUseCase.GetOAuthURL()
	if err != nil {
		log.Printf("Error getting OAuth URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error getting OAuth URL"})
		return
	}

	c.Redirect(http.StatusFound, authURL)
}

// HandleOAuthCallback handles the OAuth provider callback with the authorization code
func (oc *OAuthController) HandleOAuthCallback(c *gin.Context) {
	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No code in request"})
		return
	}

	token, err := oc.oauthUseCase.ExchangeCode(code)
	if err != nil {
		log.Printf("Error exchanging code for token: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error exchanging code for token"})
		return
	}

	user, err := oc.oauthUseCase.FetchUserInfo(token)
	if err != nil {
		log.Printf("Error fetching user info: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error fetching user info"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"user":  user,
		"token": token,
	})
}