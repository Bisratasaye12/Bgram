package models

type Env struct {
	DBURL               string
	DBNAME              string
	JWT_SECRET_KEY      string
	JWT_EXPIRATION_TIME int

	SMTP_HOST             string
	SMTP_PORT             string
	SMTP_EMAIL_FROM       string
	SMTP_USERNAME         string
	SMTP_PASSWORD         string
	SMTP_FROM             string
	APP_BASE_URL          string
	EMAIL_SUBJECT         string
	CLOUDINARY_API_KEY    string
	CLOUDINARY_CLOUD_NAME string
	CLOUDINARY_API_SECRET string
}
