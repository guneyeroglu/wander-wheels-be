package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
)

type Color struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllColors(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var colors []Color
	query := `
		SELECT 
			id,
			CASE 
				WHEN $1 = 'tr_TR' then name_tr
				ELSE name_en
			END as name
		FROM colors
	`

	rows, err := db.Query(query, lang)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var color Color

		err := rows.Scan(&color.Id, &color.Name)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		colors = append(colors, color)
	}

	return c.Status(fiber.StatusOK).JSON(colors)
}
