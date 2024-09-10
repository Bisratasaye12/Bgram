package interfaces

import (
	models "BChat/Domain/Models"
	"mime/multipart"
)



type UserUseCaseInterface interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByEmailOrUsername(email string, username string) (*models.User, error)
	RegisterUser(user *models.User, env *models.Env) error
	ActivateUserAccount(token string, password string) error
	LoginUser(user *models.User) (string, string, error)
	UpdateUserProfile(user *models.User, profilePhoto *multipart.FileHeader) (*models.User, error)
	DeleteUserAccount(id string) error
	Logout(userID string) error
	RefreshSession(refreshToken string, userID string) (string, string, error)
}
