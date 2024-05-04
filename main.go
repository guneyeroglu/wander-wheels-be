package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/guneyeroglu/wander-wheels-be/middleware"
	"github.com/guneyeroglu/wander-wheels-be/models"
)

func main() {
	app := fiber.New()
	app.Use(cors.New(), middleware.Language, middleware.Jwt)
	// app := app.Group("/app")
	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, WanderWheels")
	})
	app.Post("/login", models.Login)
	app.Post("/sign-up", models.SignUp)
	app.Get("/user-info", models.GetUserInfo)
	app.Get("/fuels", models.GetAllFuels)
	app.Get("/transmissions", models.GetAllTransmissions)
	app.Get("/colors", models.GetAllColors)
	app.Get("/brands", models.GetAllBrands)
	app.Get("/models", models.GetAllModels)
	app.Get("/cities", models.GetAllCities)
	app.Get("/cities/:id", models.GetCityById)
	app.Post("/cars", models.GetAllCars)
	app.Get("/cars/:id", models.GetCarById)
	app.Post("/rent-car", models.CreateRental)
	app.Get("/price-range", models.GetPriceRange)
	app.Get("/year-range", models.GetYearRange)
	app.Get("/seats", models.GetSeats)

	port := os.Getenv("PORT")

	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen("0.0.0.0:" + port))
}
