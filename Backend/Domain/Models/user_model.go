package models

type User struct {
	ID             string `json:"_id", bson:"_id"`
	Username       string `json:"username", bson:"username", validate: "min=4,max=20"`
	Password       string `json:"password", bson:"password", validate:"required"`
	Email          string `json:"email", bson:"email", validate:"required,email"`
	Role           string `json:"role", bson:"role"`
	ProfilePicture string `json:"profile_picture", bson:"profile_picture"`
	Bio            string `json:"bio", bson:"bio"`
}


// OAuthUser represents user information from the OAuth provider
type OAuthUser struct {
	ID    string
	Email string
	Name  string
}

// OAuthToken represents the token returned from the OAuth provider
type OAuthToken struct {
	AccessToken  string
	RefreshToken string
	ExpiresIn    int
	TokenType    string
}

// Google Certs
type GoogleCerts struct {
	Keys []struct {
		Kid string `json:"kid"`
		N   string `json:"n"`
		E   string `json:"e"`
	} `json:"keys"`
}
