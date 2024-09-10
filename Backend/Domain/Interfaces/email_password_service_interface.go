package interfaces



type EmailPasswordInterface interface {
	ValidateEmail(email string) error
	ValidatePassword(password string) error
}