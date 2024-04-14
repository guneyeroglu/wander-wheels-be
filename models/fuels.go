package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/guneyeroglu/wander-wheels-be/utils"
)

type Fuel struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

func GetAllFuels(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var fuels []Fuel
	query := `
		SELECT 
			id,
			CASE 
				WHEN $1 = 'tr_TR' then name_tr
				ELSE name_en
			END as name
		FROM fuels
		ORDER BY id ASC
	`

	rows, err := db.Query(query, lang)

	if err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"status":  fiber.StatusUnprocessableEntity,
			"message": err.Error(),
		})
	}

	for rows.Next() {
		var fuel Fuel

		err := rows.Scan(&fuel.Id, &fuel.Name)

		if err != nil {
			return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
				"status":  fiber.StatusUnprocessableEntity,
				"message": err.Error(),
			})
		}

		fuels = append(fuels, fuel)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    fuels,
		"message": utils.GetTranslation(lang, "success"),
	})
}
