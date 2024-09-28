package routers

import (
	middlewares "BChat/Delivery/Middlewares"
	interfaces "BChat/Domain/Interfaces"

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userControllers interfaces.UserControllerInterface, oauthControllers interfaces.OAuthControllerInterface) {
	
	router.GET("/oauth-redirect", oauthControllers.HandleOAuthCallback)

	// API v1 group
	v1 := router.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			userGroup.POST("/register", userControllers.RegisterUser)
			userGroup.POST("/verify/:token", userControllers.ActivateUserAccount)
			userGroup.POST("/login", userControllers.LoginUser)

			oauthGroup := userGroup.Group("/oauth")
			{
				oauthGroup.GET("/login", oauthControllers.RedirectToProvider)
			}
		}

		protectedGroup := v1.Group("/users")
		protectedGroup.Use(middlewares.AuthMiddleware())
		{
			protectedGroup.GET("/:id", userControllers.GetUserByID)
			protectedGroup.GET("/get-user", userControllers.GetUserByEmailOrUsername)
			protectedGroup.PUT("/update-profile/:id", userControllers.UpdateUserProfile)
			protectedGroup.DELETE("/delete-account/:id", userControllers.DeleteUserAccount)
			protectedGroup.POST("/logout/:id", userControllers.Logout)
			protectedGroup.POST("/refresh/:id", userControllers.RefreshSession)
		}
	}

	// Additional groups or routes can be added here
}
