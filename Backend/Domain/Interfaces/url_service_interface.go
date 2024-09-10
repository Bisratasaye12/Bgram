package interfaces

import models "BChat/Domain/Models"


type URLServiceInterface interface {
    GenerateVerificationURL(user *models.User, url_id string, env *models.Env) (string, error)
    SendVerificationEmail(email, verificationURL string) error
    VerifyUser(token string) (*models.CustomClaims, error)
}