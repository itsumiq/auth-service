package middleware

import (
	"auth-service/app/internal/config"
	"auth-service/app/internal/exception"
	"auth-service/app/internal/http"
	"context"
	"errors"
	"log/slog"
	"time"

	"github.com/gofiber/fiber/v2"
)

func HandleGlobalErrors(logger *slog.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		err := c.Next()
		if err != nil {
			if errors.Is(err, exception.ErrInternal) {
				return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
					"msg": "Internal server error",
				})
			}

			logger.Error("Unexpected error", "error", err)
			return c.Status(fiber.ErrInternalServerError.Code).JSON(fiber.Map{
				"msg": "Internal server error",
			})
		}

		return err

	}
}

func HandleTimeOut(config *config.Config) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		c.SetUserContext(ctx)

		err := c.Next()
		if err != nil && errors.Is(err, context.DeadlineExceeded) {
			return c.Status(fiber.StatusGatewayTimeout).JSON(fiber.Map{
				"msg": " Query time limit exceeded",
			})
		}
		return err
	}
}

func HandleAuthErrors(c *fiber.Ctx) error {
	err := c.Next()
	if err != nil {
		if errors.Is(err, http.ErrInvalidBody) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"msg": err.Error(),
			})
		}
	}
	return err
}
