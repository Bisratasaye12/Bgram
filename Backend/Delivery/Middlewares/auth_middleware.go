package middlewares

import (
	config "BChat/Config"
	models "BChat/Domain/Models"
	infrastructure "BChat/Infrastructure"
	"strings"

	"github.com/gin-gonic/gin"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")

		if clientToken == "" {
			c.JSON(403, gin.H{"message": "No Authorization header provided"})
			c.Abort()
			return
		}

		// Split the Authorization header to get the token
		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) == 2 {
			clientToken = strings.TrimSpace(extractedToken[1])
		} else {
			c.JSON(400, gin.H{"message": "Malformed token"})
			c.Abort()
			return
		}

		env := config.InitEnv()
		oauthService := infrastructure.NewOAuthService(env)

		// Try to validate the token as a JWT first
		jwtWrapper := infrastructure.NewJWTService(env.JWT_SECRET_KEY)
		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err == nil {
			c.Set("claims", claims)
			c.Set("token", clientToken)
			c.Next()
			return
		}

		// If JWT validation fails, try to validate it as an OAuth token
		userInfo, err := oauthService.GetUserInfo(&models.OAuthToken{AccessToken: clientToken})
		if err != nil {
			c.JSON(401, gin.H{"error": "Invalid access token"})
			c.Abort()
			return
		}

		// Set user info in the context
		c.Set("user", userInfo)
		c.Next()
	}
}
