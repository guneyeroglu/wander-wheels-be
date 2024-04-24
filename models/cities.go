package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/guneyeroglu/wander-wheels-be/utils"
)

type City struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllCities(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var cities []City

	query := `
		SELECT 
			id, 
			name
		FROM cities
		ORDER BY id ASC
	`

	rows, err := db.Query(query)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var city City

		err := rows.Scan(&city.Id, &city.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		cities = append(cities, city)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    cities,
		"message": utils.GetTranslation(lang, "success"),
	})
}

func GetCityById(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	cityId := c.Params("id")
	db := database.ConnectDb()
	defer db.Close()

	var cities []City

	query := `
		SELECT 
			id, 
			name
		FROM cities
		WHERE id = $1
		ORDER BY id ASC
	`

	rows, err := db.Query(
		query,  // query
		cityId, // $1

	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var city City

		err := rows.Scan(&city.Id, &city.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		cities = append(cities, city)
	}

	if cities == nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"status":  fiber.StatusNotFound,
			"data":    nil,
			"message": utils.GetTranslation(lang, "notFound"),
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    cities[0],
		"message": utils.GetTranslation(lang, "success"),
	})
}
