package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
)

type City struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllCities(c *fiber.Ctx) error {
	db := database.ConnectDb()
	defer db.Close()

	var cities []City

	rows, err := db.Query(`
		SELECT 
			id, 
			name
		FROM cities
	`)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var city City

		err := rows.Scan(&city.Id, &city.Name)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		cities = append(cities, city)
	}

	return c.Status(fiber.StatusOK).JSON(cities)
}
