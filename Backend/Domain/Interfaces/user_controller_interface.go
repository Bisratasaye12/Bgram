package interfaces

import (

	"github.com/gin-gonic/gin"
)

type UserControllerInterface interface {
	GetUserByID(ctx *gin.Context)
	GetUserByEmailOrUsername(ctx *gin.Context)
	RegisterUser(ctx *gin.Context)
	ActivateUserAccount(ctx *gin.Context)
	LoginUser(ctx *gin.Context)
	UpdateUserProfile(ctx *gin.Context)
	DeleteUserAccount(ctx *gin.Context)
	Logout(ctx *gin.Context)
	RefreshSession(ctx *gin.Context)
}