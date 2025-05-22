package external

import (
	"event-management-system/internal/event"
	"event-management-system/internal/user"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func healthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status": "UP",
	})
}

func CheckToken() fiber.Handler {
	return func(c *fiber.Ctx) error {
		fmt.Println("Do Check Token")
		return c.Next()
	}
}

func PublicRoutes(r fiber.Router, db *gorm.DB) {
	apiV1NoGuard := r.Group("/api/v1")
	// apiV1Guard := r.Group("/api/v1", CheckToken())
	apiV1Guard := r.Group("/api/v1")
	apiV1NoGuard.Get("/health", healthCheck)

	event.Router(apiV1Guard, db)
	user.Router(apiV1Guard, db)
}
