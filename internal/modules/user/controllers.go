package user

import "github.com/gofiber/fiber/v2"

func (um *UserModule) getAllUsers(c *fiber.Ctx) error {
	// Placeholder implementation
	return c.JSON(fiber.Map{
		"message": "Get all users - not implemented",
	})
}

func (um *UserModule) createUser(c *fiber.Ctx) error {
	// Placeholder implementation
	return c.JSON(fiber.Map{
		"message": "Create user - not implemented",
	})
}
