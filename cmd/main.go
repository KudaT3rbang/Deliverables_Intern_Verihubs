package main

import (
	"lendbook/internal/delivery/http/handler"
	"lendbook/internal/infrastructure/postgres"
	"lendbook/internal/usecase"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/labstack/echo/v4"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	jwtsecret := os.Getenv("JWT_SECRET")
	if jwtsecret == "" {
		log.Fatal("JWT_SECRET environment variable not set")
	}

	dbPool, err := postgres.InitDb()
	if err != nil {
		log.Fatal("Error connecting to database")
	}

	defer dbPool.Close()

	userRepo := postgres.NewUserRepository(dbPool)
	userUseCase := usecase.NewUserUsecase(userRepo, jwtsecret)
	userHandler := handler.NewUserHandler(userUseCase)

	e := echo.New()
	e.POST("/register", userHandler.Register)
	e.POST("/login", userHandler.Login)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	if err := e.Start(":" + port); err != nil {
		e.Logger.Fatal(err)
	}

}
