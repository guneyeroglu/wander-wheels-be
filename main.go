package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/middleware"
	"github.com/guneyeroglu/wander-wheels-be/models"
)

func main() {
	app := fiber.New()
	app.Use(middleware.Language)
	api := app.Group("/api")

	api.Get("/users", models.GetAllUsers)
	api.Get("/fuels", models.GetAllFuels)
	api.Get("/transmissions", models.GetAllTransmissions)
	api.Get("/colors", models.GetAllColors)

	log.Fatal(app.Listen(":3000"))
}
