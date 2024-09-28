package interfaces

import "github.com/gin-gonic/gin"

type OAuthControllerInterface interface {
	RedirectToProvider(c *gin.Context)
	HandleOAuthCallback(c *gin.Context)
}