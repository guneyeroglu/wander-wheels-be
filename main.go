package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
)

func helloWorld(c *fiber.Ctx) error {
	return c.SendString("Hello, world!")
}

func main() {
	app := fiber.New()

	app.Get("/api", helloWorld)
	log.Fatal(app.Listen(":3000"))
}
