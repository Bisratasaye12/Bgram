package config

import (
	models "BChat/Domain/Models"
	"log"

	viper "github.com/spf13/viper"
)

func InitEnv() *models.Env {
	// Set the environment variables
	env := &models.Env{}

	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading .env file: %v", err)
	}

	env.DBURL = viper.GetString("DBURL")
	env.DBNAME = viper.GetString("DBNAME")
	env.JWT_SECRET_KEY = viper.GetString("JWT_SECRET_KEY")
	env.JWT_EXPIRATION_TIME = viper.GetInt("JWT_EXPIRATION_TIME")
	env.SMTP_HOST = viper.GetString("SMTP_HOST")
	env.SMTP_PORT = viper.GetString("SMTP_PORT")
	env.SMTP_EMAIL_FROM = viper.GetString("SMTP_EMAIL_FROM")
	env.SMTP_USERNAME = viper.GetString("SMTP_USERNAME")
	env.SMTP_PASSWORD = viper.GetString("SMTP_PASSWORD")
	env.SMTP_FROM = viper.GetString("SMTP_FROM")
	env.APP_BASE_URL = viper.GetString("APP_BASE_URL")
	env.EMAIL_SUBJECT = viper.GetString("EMAIL_SUBJECT")
	env.CLOUDINARY_API_KEY = viper.GetString("CLOUDINARY_API_KEY")
	env.CLOUDINARY_CLOUD_NAME = viper.GetString("CLOUDINARY_CLOUD_NAME")
	env.CLOUDINARY_API_SECRET = viper.GetString("CLOUDINARY_API_SECRET")

	return env
}
