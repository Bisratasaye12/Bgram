package middlewares

import (
	infrastructure "BChat/Infrastructure"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)



func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		clientToken := c.Request.Header.Get("Authorization")
		if clientToken == "" {
		 // If the Authorization header is not present, return a 403 status code
		 c.JSON(403, gin.H{"message": "No Authorization header provided"})
		 c.Abort()
		 return
		}

		// Split the Authorization header to get the token
		extractedToken := strings.Split(clientToken, "Bearer ")
		if len(extractedToken) == 2 {
		 // Trim the token
		 clientToken = strings.TrimSpace(extractedToken[1])
		} else {
		 // If the token is not in the correct format, return a 400 status code
		 c.JSON(400, gin.H{"message": "Malformed token"})
		 c.Abort()
		 return
		}

		viper.SetConfigFile("config.json")
		viper.ReadInConfig()
		secretKey := viper.GetString("JWT_SECRET_KEY")

		// Create a JwtWrapper with the secret key
		jwtWrapper := infrastructure.NewJWTService(secretKey)

		// Validate the token
		claims, err := jwtWrapper.ValidateToken(clientToken)
		if err != nil {
		 // If the token is not valid, return a 401 status code
		 c.JSON(401, gin.H{"error": err.Error()})
		 c.Abort()
		 return
		}
		// Set the claims in the context
		c.Set("claims", claims)
		c.Set("token", clientToken)
		// Continue to the next handler
		c.Next()
	}
}