package user

import (
	"event-management-system/dto/entity"
	"event-management-system/dto/model"
	"event-management-system/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type (
	Handler interface {
		Create(c *fiber.Ctx) error
		FindAll(c *fiber.Ctx) error
		FindByID(c *fiber.Ctx) error
		Update(c *fiber.Ctx) error
	}

	handler struct {
		service   Service
		validator *utils.Validator
	}
)

func NewHandler(service Service, validator *utils.Validator) Handler {
	return &handler{service, validator}
}

func (h *handler) Create(c *fiber.Ctx) error {
	payload := entity.User{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Errors:  err,
		})
	}

	errV := h.validator.Validate(payload)
	if errV != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, []*utils.ErrorResponse]{
			Message: "Validation Error",
			Errors:  errV,
		})
	}

	validRole := payload.CheckRole()
	if !validRole {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, error]{
			Message: "Role is not valid",
			Errors:  nil,
		})
	}

	err := h.service.Create(&payload)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.JSON(model.Response[interface{}, error]{
		Message: "Success",
		Success: true,
	})

}

func (h *handler) FindAll(c *fiber.Ctx) error {
	criteria := model.FindAllUserCriteria{}
	if err := c.QueryParser(&criteria); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Errors:  err,
		})
	}

	result, count, err := h.service.FindAll(criteria)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.JSON(model.ResponseWithOffset[[]*entity.User, error]{
		Message: "Success",
		Success: true,
		Data:    result,
		Limit:   criteria.Limit,
		Offset:  criteria.Offset,
		Total:   count,
	})
}

func (h *handler) FindByID(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Errors:  err,
		})
	}

	user, err := h.service.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.JSON(model.Response[*entity.User, error]{
		Message: "Success",
		Success: true,
		Data:    user,
	})
}

func (h *handler) Update(c *fiber.Ctx) error {
	id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Errors:  err,
		})
	}
	payload := entity.User{}

	if err := c.BodyParser(&payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Errors:  err,
		})
	}

	errV := h.validator.Validate(payload)
	if errV != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.User, []*utils.ErrorResponse]{
			Message: "Validation Error",
			Errors:  errV,
		})
	}

	err = h.service.Update(uint(id), &payload)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response[entity.User, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.JSON(model.Response[interface{}, error]{
		Message: "Success",
		Success: true,
	})
}
