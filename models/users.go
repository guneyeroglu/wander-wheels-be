package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guneyeroglu/wander-wheels-be/database"
)

type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
}

func GetAllUsers(c *fiber.Ctx) error {
	db := database.ConnectDb()
	defer db.Close()

	var users []User

	rows, err := db.Query(`
		SELECT 
			id, 
			username, 
			password 
		FROM users
	`)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var user User

		err := rows.Scan(&user.Id, &user.Username, &user.Password)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
