package middleware

import (
	"github.com/gofiber/fiber/v3"
)

func Language(c fiber.Ctx) error {

	lang := c.Get("Accept-Language")

	if lang == "" {
		lang = "en_EN"
	}

	c.Locals("lang", lang)

	return c.Next()

}
