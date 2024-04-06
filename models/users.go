package models

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/guneyeroglu/wander-wheels-be/database"
)

type Role struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type User struct {
	Id       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Password string    `json:"password"`
	Role     Role      `json:"role"`
}

func GetAllUsers(c *fiber.Ctx) error {
	db := database.ConnectDb()
	defer db.Close()

	var users []User

	rows, err := db.Query(`
		SELECT 
			U.id, 
			U.username, 
			U.password,
			R.id,
			R.name 
		FROM users AS U
		JOIN roles AS R ON r.id = u.role_id
	`)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	for rows.Next() {
		var user User
		var role Role

		err := rows.Scan(&user.Id, &user.Username, &user.Password, &role.Id, &role.Name)

		if err != nil {
			return c.Status(fiber.StatusInternalServerError).SendString(err.Error())
		}

		user.Role = role
		users = append(users, user)
	}

	return c.Status(fiber.StatusOK).JSON(users)
}
