package interfaces




type SessionRepositoryInterface interface {
	SaveTokens(userID string, accessToken string, refreshToken string) error
	DeleteTokens(userID string) error
	GetTokens(userID string) (string, string, error)
}