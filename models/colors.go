package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/guneyeroglu/wander-wheels-be/utils"
)

type Color struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
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
			END AS name,
			code
		FROM colors
		ORDER BY id ASC
	`

	rows, err := db.Query(
		query, // query
		lang,  // $1
	)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var color Color

		err := rows.Scan(&color.Id, &color.Name, &color.Code)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		colors = append(colors, color)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    colors,
		"message": utils.GetTranslation(lang, "success"),
	})
}
