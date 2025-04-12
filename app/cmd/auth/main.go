package main

import (
	"auth-service/app/internal/config"
	"auth-service/app/internal/http/handler"
	"auth-service/app/internal/http/middleware"
	"auth-service/app/internal/lib"
	"auth-service/app/internal/repository/postgres"
	"auth-service/app/internal/service"
	"fmt"
	"log/slog"
	"os"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jmoiron/sqlx"
)

func setupLogger() *slog.Logger {
	handler := slog.NewTextHandler(
		os.Stdout,
		&slog.HandlerOptions{Level: slog.LevelInfo, AddSource: true},
	)
	return slog.New(handler)
}

func main() {
	logger := setupLogger()
	config := config.Get()

	validate := validator.New(validator.WithRequiredStructEnabled())
	err := validate.RegisterValidation("login", lib.ValidateLogin)
	if err != nil {
		panic(err)
	}

	storagePath := fmt.Sprintf(
		"host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
		config.Database.Host,
		config.Database.Port,
		config.Database.User,
		config.Database.Password,
		config.Database.Name,
		"disable",
	)

	db, err := sqlx.Connect("postgres", storagePath)
	if err != nil {
		logger.Error("Database connect error", "error", err)
		return
	}

	uow := postgres.NewUnitOfWork(db, logger)
	authService := service.NewAuthService(uow, logger)
	authHandler := handler.NewAuth(authService, validate)

	app := fiber.New()
	api := app.Group(
		"/api",
		middleware.HandleGlobalErrors(logger),
		middleware.HandleTimeOut(config),
	)

	auth := api.Group("/auth", middleware.HandleAuthErrors)
	auth.Post("/registration", authHandler.RegisterUser)
	auth.Get("/login", authHandler.LoginUser)
	auth.Patch("/refresh", authHandler.RefreshTokens)

	app.Listen("0.0.0.0:8000")
}
