package main

import (
	config "BChat/Config"
	controllers "BChat/Delivery/Controllers"
	routers "BChat/Delivery/Routers"
	infrastructure "BChat/Infrastructure"
	repository "BChat/Repository"
	usecases "BChat/UseCases"
	"context"

	"github.com/gin-gonic/gin"
)

func main() {
	router := gin.Default()

	env := config.InitEnv()
	
	// Setup database
	db := config.InitDB(env)
	defer db.Client().Disconnect(context.Background())

	// Setup repositories
	userRepo := repository.NewUserRepository(db)
	urlRepo := repository.NewVerificationURLRepository(db)
	sessionRepo := repository.NewSessionRepository(db)

	// setup infrastructures
	cloudinaryService, err := infrastructure.NewCloudinaryService(env)
	if err != nil {
		panic(err)
	}
	jwtService := infrastructure.NewJWTService(env.JWT_SECRET_KEY)
	urlService := infrastructure.NewURLService(jwtService, urlRepo, env)
	emailPasswordService := infrastructure.NewEmailPasswordService()
	oauthService := infrastructure.NewOAuthService(env)

	// Setup use cases
	userUsecase := usecases.NewUserUseCase(userRepo, jwtService, urlService, cloudinaryService,emailPasswordService, sessionRepo)
	oauthUseCase := usecases.NewOAuthUseCase(oauthService)
	// Setup controllers
	userController := controllers.NewUserController(userUsecase, env)
	oauthController := controllers.NewOAuthController(oauthUseCase)

	// Setup routes
	routers.SetupRoutes(router, userController, oauthController)

	router.Run(env.APP_BASE_URL)

}
