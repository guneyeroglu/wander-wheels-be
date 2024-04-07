package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/guneyeroglu/wander-wheels-be/database"
	"github.com/guneyeroglu/wander-wheels-be/utils"
)

type Model struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Brand Brand  `json:"brand"`
}

func GetAllModels(c *fiber.Ctx) error {
	lang := c.Locals("lang").(string)
	db := database.ConnectDb()
	defer db.Close()

	var models []Model

	rows, err := db.Query(`
		SELECT 
			M.id AS model_id, 
			M.name AS model_name,
			M.brand_id AS brand_id,
			B.name AS brand_name
		FROM models AS M
		JOIN brands AS B ON b.id = M.brand_id
		ORDER BY model_id ASC
		`)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var model Model
		var brand Brand

		err := rows.Scan(&model.Id, &model.Name, &brand.Id, &brand.Name)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		model.Brand = brand
		models = append(models, model)
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  fiber.StatusOK,
		"data":    models,
		"message": utils.GetTranslation(lang, "success"),
	})
}
