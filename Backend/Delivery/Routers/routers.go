package routers

import (
	interfaces "BChat/Domain/Interfaces"
	"BChat/Delivery/Middlewares" 

	"github.com/gin-gonic/gin"
)

func SetupRoutes(router *gin.Engine, userControllers interfaces.UserControllerInterface) {
	// API v1 group
	v1 := router.Group("/api/v1")
	{
		userGroup := v1.Group("/users")
		{
			userGroup.POST("/register", userControllers.RegisterUser)
			userGroup.POST("/verify/:token", userControllers.ActivateUserAccount)
			userGroup.POST("/login", userControllers.LoginUser)
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
