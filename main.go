package main

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/guneyeroglu/wander-wheels-be/middleware"
	"github.com/guneyeroglu/wander-wheels-be/models"
)

func envPortOr(port string) string {
	if envPort := os.Getenv("PORT"); envPort != "" {
		return ":" + envPort
	}

	return ":" + port
}

func main() {
	app := fiber.New()
	app.Use(cors.New(), middleware.Language, middleware.Jwt)
	api := app.Group("/api")

	api.Post("/login", models.Login)
	api.Post("/sign-up", models.SignUp)
	api.Get("/user-info", models.GetUserInfo)
	api.Get("/fuels", models.GetAllFuels)
	api.Get("/transmissions", models.GetAllTransmissions)
	api.Get("/colors", models.GetAllColors)
	api.Get("/brands", models.GetAllBrands)
	api.Get("/models", models.GetAllModels)
	api.Get("/cities", models.GetAllCities)
	api.Get("/cities/:id", models.GetCityById)
	api.Post("/cars", models.GetAllCars)
	api.Get("/cars/:id", models.GetCarById)
	api.Post("/rent-car", models.CreateRental)
	api.Get("/price-range", models.GetPriceRange)
	api.Get("/year-range", models.GetYearRange)
	api.Get("/seats", models.GetSeats)

	app.Listen(envPortOr("3000"))
}
