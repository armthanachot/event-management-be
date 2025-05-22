package event
import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"event-management-system/pkg/utils"

)

func Router(r fiber.Router, db *gorm.DB) {
	repository := NewRepository(db)
	service := NewService(repository)

	validator := utils.NewValidator()

	handler := NewHandler(service, validator)

	groupRoute := r.Group("/event")
	groupRoute.Post("/", handler.Create)
    groupRoute.Get("/", handler.FindAll)
    groupRoute.Get("/:id", handler.FindByID)
    groupRoute.Patch("/:id", handler.Update)
    groupRoute.Delete("/:id", handler.Delete)
}
