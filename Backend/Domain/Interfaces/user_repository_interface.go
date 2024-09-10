package interfaces

import models "BChat/Domain/Models"



type UserRepositoryInterface interface {
	GetUserByID(id string) (*models.User, error)
	GetUserByEmailOrUsername(email string, username string) (*models.User, error)
	CreateUser(user *models.User) (*models.User, error)
	UpdateUser(user *models.User) (*models.User, error)
	DeleteUser(id string) error
}