package models

type Session struct {
	ID           int    `json:"_id", bson:"_id"`
	UserID       int    `json:"user_id", bson:"user_id"`
	AccessToken  string `json:"access_token", bson:"access_token"`
	RefreshToken string `json:"refresh_token", bson:"refresh_token"`
	ExpiredAt    string `json:"expired_at", bson:"expired_at"`
}
