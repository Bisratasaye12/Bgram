package controllers

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"

	"github.com/gin-gonic/gin"
)

type UserController struct {
	UserUseCase interfaces.UserUseCaseInterface
	env         *models.Env
}

func NewUserController(userUseCase interfaces.UserUseCaseInterface, env *models.Env) interfaces.UserControllerInterface {
	return &UserController{
		UserUseCase: userUseCase,
		env:         env,
	}
}

func (uc *UserController) GetUserByID(ctx *gin.Context) {
	id := ctx.Param("id")
	user, err := uc.UserUseCase.GetUserByID(id)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"user": user,
	})
}

func (uc *UserController) GetUserByEmailOrUsername(ctx *gin.Context) {
	user := &models.User{}
	if ctx.ShouldBindJSON(user) != nil {
		ctx.IndentedJSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	email := user.Email
	username := user.Username
	user, err := uc.UserUseCase.GetUserByEmailOrUsername(email, username)
	if err != nil {
		if err.Error() == "user does not exist" {
			ctx.IndentedJSON(404, gin.H{
				"error": "user does not exist",
			})
			return
		}
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"user": user,
	})
}

func (uc *UserController) RegisterUser(ctx *gin.Context) {
	user := &models.User{}
	if ctx.ShouldBindJSON(user) != nil {
		ctx.IndentedJSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	err := uc.UserUseCase.RegisterUser(user, uc.env)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"message": "check your email for verification",
	})
}

func (uc *UserController) ActivateUserAccount(ctx *gin.Context) {
	token := ctx.Param("token")

	user := &models.User{}
	if ctx.ShouldBindJSON(user) != nil {
		ctx.IndentedJSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}

	password := user.Password
	err := uc.UserUseCase.ActivateUserAccount(token, password)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"message": "user account activated",
	})
}

func (uc *UserController) LoginUser(ctx *gin.Context) {
	user := &models.User{}
	if ctx.ShouldBindJSON(user) != nil {
		ctx.IndentedJSON(400, gin.H{
			"error": "invalid request",
		})
		return
	}
	AccessToken, RefreshToken, err := uc.UserUseCase.LoginUser(user)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(), "message": "invalid credentials",
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"access token":  AccessToken,
		"refresh token": RefreshToken,
	})
}

func (uc *UserController) UpdateUserProfile(ctx *gin.Context) {
	id := ctx.Param("id")
	user := &models.User{}

	user.ID = id
	profilePhoto, _ := ctx.FormFile("profile_picture")
	user.Email = ctx.PostForm("email")
	user.Username = ctx.PostForm("username")
	user.Password = ctx.PostForm("password")
	user.Bio = ctx.PostForm("bio")

	user, err := uc.UserUseCase.UpdateUserProfile(user, profilePhoto)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"user": user,
	})
}

func (uc *UserController) DeleteUserAccount(ctx *gin.Context) {
	id := ctx.Param("id")
	err := uc.UserUseCase.DeleteUserAccount(id)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"message": "user account deleted",
	})
}


func (uc *UserController) Logout(ctx *gin.Context) {
	userID := ctx.Param("id")
	err := uc.UserUseCase.Logout(userID)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"message": "user logged out",
	})
}

func (uc *UserController) RefreshSession(ctx *gin.Context) {
	userID := ctx.Param("id")
	refreshToken := ctx.MustGet("token").(string)

	AccessToken, RefreshToken, err := uc.UserUseCase.RefreshSession(refreshToken, userID)
	if err != nil {
		ctx.IndentedJSON(500, gin.H{
			"error": err.Error(),
		})
		return
	}
	ctx.IndentedJSON(200, gin.H{
		"access token":  AccessToken,
		"refresh token": RefreshToken,
	})
}