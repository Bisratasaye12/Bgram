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
