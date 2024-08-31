package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

func UnknowMethod(c *fiber.Ctx) error {
	return c.Status(fiber.StatusMethodNotAllowed).JSON(fiber.Map{
		"code":    fiber.StatusMethodNotAllowed,
		"status":  false,
		"data":    "",
		"message": "Method Not Allowed",
	})
}
