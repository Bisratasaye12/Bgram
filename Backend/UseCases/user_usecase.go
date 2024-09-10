package usecases

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
	"fmt"
	"mime/multipart"

	"github.com/google/uuid"
)

type UserUseCase struct {
	UserRepo          interfaces.UserRepositoryInterface
	SessionRepo       interfaces.SessionRepositoryInterface
	JWTService        interfaces.JWTServiceInterface
	UrlService        interfaces.URLServiceInterface
	EPService         interfaces.EmailPasswordInterface
	CloudinaryService interfaces.CloudinaryServiceInterface
}

func NewUserUseCase(
	userRepo interfaces.UserRepositoryInterface,
	jwt interfaces.JWTServiceInterface,
	urlSrv interfaces.URLServiceInterface,
	cloudinray interfaces.CloudinaryServiceInterface,
	emailPasswordvalidation interfaces.EmailPasswordInterface,
	sessionRepo interfaces.SessionRepositoryInterface) interfaces.UserUseCaseInterface {
	return &UserUseCase{
		UserRepo:          userRepo,
		JWTService:        jwt,
		UrlService:        urlSrv,
		CloudinaryService: cloudinray,
		EPService:         emailPasswordvalidation,
		SessionRepo:       sessionRepo,
	}
}

func (u *UserUseCase) GetUserByID(id string) (*models.User, error) {
	user, err := u.UserRepo.GetUserByID(id)
	if err != nil {
		if err.Error() == "user does not exist" {
			return nil, err
		}
		return nil, fmt.Errorf("internal server error")
	}
	return user, nil
}

func (u *UserUseCase) GetUserByEmailOrUsername(email string, username string) (*models.User, error) {
	user, err := u.UserRepo.GetUserByEmailOrUsername(email, username)
	if err != nil {
		if err.Error() == "user does not exist" {
			return nil, err
		}
		return nil, fmt.Errorf("internal server error")
	}
	return user, nil
}

func (u *UserUseCase) RegisterUser(user *models.User, env *models.Env) error {
	// validate email
	if err := u.EPService.ValidateEmail(user.Email); err != nil {
		return fmt.Errorf("invalid email")
	}

	// validate user
	existingUser, _ := u.UserRepo.GetUserByEmailOrUsername(user.Email, user.Username)
	if existingUser != nil {
		return fmt.Errorf("user already exists")
	}

	// generete verificiation url
	url_id := uuid.New().String()
	url, err := u.UrlService.GenerateVerificationURL(user, url_id, env)
	if err != nil {
		return fmt.Errorf("internal server error")
	}

	// send email
	err = u.UrlService.SendVerificationEmail(user.Email, url)
	if err != nil {
		return fmt.Errorf("internal server error send")
	}
	return nil
}

func (u *UserUseCase) ActivateUserAccount(token string, password string) error {
	claims, err := u.UrlService.VerifyUser(token)
	if err != nil {
		return fmt.Errorf("error activating user account")
	}

	// verify password strength
	if err := u.EPService.ValidatePassword(password); err != nil {
		return err
	}

	hashedpassword, err := u.JWTService.HashPassword(password)
	if err != nil {
		return fmt.Errorf("internal server error")
	}
	user := &models.User{
		Password: hashedpassword,
		Username: claims.Username,
		Email:    claims.UserEmail,
		Role:     claims.Role,
	}
	u.UserRepo.CreateUser(user)
	return nil
}

func (u *UserUseCase) LoginUser(user *models.User) (string, string, error) {
	existingUser, _ := u.UserRepo.GetUserByEmailOrUsername(user.Email, user.Username)
	if existingUser == nil {
		return "", "", fmt.Errorf("user does not exist")
	}

	if u.JWTService.CheckPasswordHash(user.Password, existingUser.Password) {
		AccessToken, err := u.JWTService.GenerateToken(existingUser, 24, "")
		RefreshToken, rerr := u.JWTService.GenerateRefreshToken(existingUser, 24)

		if err != nil {
			return "", "", fmt.Errorf("internal server error")
		}

		if rerr != nil {
			return "", "", fmt.Errorf("internal server error")
		}

		u.SessionRepo.SaveTokens(existingUser.ID, AccessToken, RefreshToken)
		return AccessToken, RefreshToken, nil
	}
	return "", "", fmt.Errorf("invalid password")
}

func (u *UserUseCase) UpdateUserProfile(user *models.User, profilePhoto *multipart.FileHeader) (*models.User, error) {
	if profilePhoto != nil {
		file, err := profilePhoto.Open()
		if err != nil {
			return nil, fmt.Errorf("invalid picture file")
		}
		defer file.Close()

		photoURL, err := u.CloudinaryService.UploadImage(file, profilePhoto)
		if err != nil {
			return nil, err
		}
		user.ProfilePicture = photoURL
	}

	if user.Password != "" {
		hashedpassword, err := u.JWTService.HashPassword(user.Password)
		if err != nil {
			return nil, fmt.Errorf("internal server error")
		}
		user.Password = hashedpassword
	}

	updatedUser, err := u.UserRepo.UpdateUser(user)
	if err != nil {
		return nil, fmt.Errorf("internal server error")
	}
	return updatedUser, nil
}

func (u *UserUseCase) DeleteUserAccount(id string) error {
	err := u.UserRepo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("internal server error")
	}

	err = u.SessionRepo.DeleteTokens(id)
	if err != nil {
		return fmt.Errorf("internal server error")
	}
	return nil
}

func (u *UserUseCase) Logout(userID string) error {
	err := u.SessionRepo.DeleteTokens(userID)
	if err != nil {
		return fmt.Errorf("internal server error")
	}
	return nil
}

func (u *UserUseCase) RefreshSession(refreshToken string, userID string) (string, string, error) {
	claims, err := u.JWTService.ValidateToken(refreshToken)
	if err != nil {
		return "", "", fmt.Errorf("invalid token")
	}

	_, existingRefToken, err := u.SessionRepo.GetTokens(userID)
	if err != nil {
		return "", "", fmt.Errorf("internal server error")
	}

	if existingRefToken != refreshToken {
		return "", "", fmt.Errorf("invalid refresh token")
	}
	usr := &models.User{
		ID:       claims.UserID,
		Email:    claims.UserEmail,
		Username: claims.Username,
		Role:     claims.Role,
	}

	u.SessionRepo.DeleteTokens(usr.ID)

	AccessToken, err := u.JWTService.GenerateToken(usr, 24, "")
	if err != nil {
		return "", "", fmt.Errorf("internal server error")
	}

	simpleUsr := &models.User{
		ID:    claims.UserID,
		Email: claims.UserEmail,
	}
	RefreshToken, err := u.JWTService.GenerateRefreshToken(simpleUsr, 24)
	if err != nil {
		return "", "", fmt.Errorf("internal server error")
	}

	err = u.SessionRepo.SaveTokens(usr.ID, AccessToken, RefreshToken)
	if err != nil {
		return "", "", fmt.Errorf("internal server error")
	}
	return AccessToken, RefreshToken, nil
}
