package middleware

import (
	"os"

	jwtware "github.com/gofiber/contrib/jwt"
	"github.com/gofiber/fiber/v2"
)

func Jwt(c *fiber.Ctx) error {
	jwtSecretCode := os.Getenv("JWT_SECRET_CODE")

	jwtware.New(jwtware.Config{
		SigningKey: jwtware.SigningKey{
			Key: []byte(jwtSecretCode),
		},
	})

	return c.Next()
}
