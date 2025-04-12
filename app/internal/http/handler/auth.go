package handler

import (
	"auth-service/app/internal/exception"
	"auth-service/app/internal/http"
	"auth-service/app/internal/model"
	"auth-service/app/internal/service"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

type authHandler struct {
	authService service.Auth
	validate    *validator.Validate
}

func NewAuth(authService service.Auth, validate *validator.Validate) *authHandler {
	return &authHandler{authService: authService, validate: validate}
}

func (h *authHandler) RegisterUser(c *fiber.Ctx) error {
	registrationRequest := &model.UserRegisterRequest{}
	if err := http.ValidateBodyWithData(c, h.validate, registrationRequest); err != nil {
		return err
	}

	ctx := c.UserContext()
	registrationResponse, err := h.authService.RegisterUser(ctx, registrationRequest)
	if err != nil {
		if errors.Is(err, exception.ErrDuplicateEntry) {
			return c.Status(fiber.StatusConflict).JSON(fiber.Map{
				"msg": "User already exists",
			})
		}
		return err
	}

	c.Set("Authorization", fmt.Sprintf("Bearer %s", registrationResponse.AccessToken))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"refresh_token": registrationResponse.RefreshToken,
	})

}

func (h *authHandler) LoginUser(c *fiber.Ctx) error {
	loginRequest := &model.UserLoginRequest{}
	if err := http.ValidateBodyWithData(c, h.validate, loginRequest); err != nil {
		return err
	}

	ctx := c.UserContext()
	loginResponse, err := h.authService.LoginUser(ctx, loginRequest)
	if err != nil {
		if errors.Is(err, exception.ErrNotFound) || errors.Is(err, exception.ErrVerify) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"msg": "Invalid login data",
			})
		}
		return err
	}

	c.Set("Authorization", fmt.Sprintf("Bearer %s", loginResponse.AccessToken))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"refresh_token": loginResponse.RefreshToken,
	})
}

func (h *authHandler) RefreshTokens(c *fiber.Ctx) error {
	refreshTokensRequest := &model.RefreshTokensRequest{}
	if err := http.ValidateBody(c, refreshTokensRequest); err != nil {
		return err
	}

	ctx := c.UserContext()
	refreshResponse, err := h.authService.RefreshTokens(ctx, refreshTokensRequest.RefreshToken)
	if err != nil {
		if errors.Is(err, exception.ErrTokenExpired) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "Refresh token expired",
			})
		}

		if errors.Is(err, exception.ErrNotFound) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"msg": "Invalid token",
			})
		}
		return err
	}

	c.Set("Authorization", fmt.Sprintf("Bearer %s", refreshResponse.AccessToken))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"refresh_token": refreshResponse.RefreshToken,
	})
}
