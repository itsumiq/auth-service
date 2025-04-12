package http

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

func ValidateBody(c *fiber.Ctx, model any) error {
	if err := c.BodyParser(model); err != nil {
		return ErrInvalidBody
	}
	return nil
}

func ValidateBodyWithData(c *fiber.Ctx, validate *validator.Validate, model any) error {
	if err := ValidateBody(c, model); err != nil {
		return err
	}
	if err := validate.Struct(model); err != nil {
		return ErrInvalidBody
	}

	return nil
}
