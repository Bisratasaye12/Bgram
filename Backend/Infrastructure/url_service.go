package infrastructure

import (
	interfaces "BChat/Domain/Interfaces"
	models "BChat/Domain/Models"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"gopkg.in/gomail.v2"
)

type urlService struct {
	jwtSvc       interfaces.JWTServiceInterface
	urlRepo      interfaces.VerificationURLRepositoryInterface
	smtpHost     string
	smtpPort     string
	smtpUser     string
	smtpPass     string
	emailFrom    string
	appBaseURL   string
	emailSubject string
}

// NewURLService creates a new instance of URLService
func NewURLService(jwtSvc interfaces.JWTServiceInterface, urlRepo interfaces.VerificationURLRepositoryInterface, env *models.Env) interfaces.URLServiceInterface {
	return &urlService{
		jwtSvc:       jwtSvc,
		urlRepo:      urlRepo,
		smtpHost:     env.SMTP_HOST,
		smtpPort:     env.SMTP_PORT,
		smtpUser:     env.SMTP_USERNAME,
		smtpPass:     env.SMTP_PASSWORD,
		emailFrom:    env.SMTP_EMAIL_FROM,
		appBaseURL:   env.APP_BASE_URL,
		emailSubject: env.EMAIL_SUBJECT,
	}
}

// GenerateVerificationURL generates a verification URL for a user
func (s *urlService) GenerateVerificationURL(user *models.User, url_id string, env *models.Env) (string, error) {
	// Generate a unique token with JWT
	token, err := s.jwtSvc.GenerateToken(user, time.Duration(env.JWT_EXPIRATION_TIME)*time.Hour, url_id)
	if err != nil {
		return "", err
	}

	// Create the verification URL
	verificationURL := fmt.Sprintf("%s/verify?token=%s", s.appBaseURL, token)

	// Save the URL to the repository
	urlModel := &models.VerificationURL{
		UrlID:     url_id,
		URL:       verificationURL,
        StandardClaims: jwt.StandardClaims{
            ExpiresAt: time.Now().Add(time.Duration(1) * time.Hour).Unix(),
        },
	}

	
	_, err = s.urlRepo.SaveUrl(urlModel)
	if err != nil {
		return "", err
	}

	return verificationURL, nil
}

// SendVerificationEmail sends a verification email to the user
func (s *urlService) SendVerificationEmail(email, verificationURL string) error {
	// Compose the email body
	body := fmt.Sprintf("Hello,\n\nPlease verify your email by clicking the link below:\n%s\n\nThank you!", verificationURL)

	m := gomail.NewMessage()
	m.SetHeader("From", s.smtpUser)
	m.SetHeader("To", email)
	m.SetHeader("Subject", s.emailSubject)
	m.SetBody("text/plain", body)

	port, err := strconv.Atoi(s.smtpPort)
	if err != nil {
		return fmt.Errorf("invalid SMTP port configuration")
	}

	d := gomail.NewDialer(s.smtpHost, port, s.smtpUser, s.smtpPass)

	if err := d.DialAndSend(m); err != nil {
		return fmt.Errorf("error occurred while sending email")
	}

	log.Println("Email sent successfully")

	return nil
}

// VerifyUser verifies the user's email based on the token
func (s *urlService) VerifyUser(token string) (*models.CustomClaims, error) {
	// Validate the token
	claims, err := s.jwtSvc.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// Retrieve the verification URL from the repository
	urlModel, err := s.urlRepo.GetUrlByID(claims.UrlID)
    
	if err != nil {
		return nil, err
	}

	if urlModel.URL != fmt.Sprintf("%s/verify?token=%s", s.appBaseURL, token) {
		return nil, fmt.Errorf("invalid verification URL")
	}

	// Delete the verification URL after successful verification
	err = s.urlRepo.DeleteUrlByID(claims.UrlID)
	if err != nil {
		return nil, err
	}

	return claims, nil
}
