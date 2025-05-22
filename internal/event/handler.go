package event

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
        Delete(c *fiber.Ctx) error
    }

    handler struct {
        service Service
		validator *utils.Validator
    }
)

func NewHandler(service Service, validator *utils.Validator) Handler {
    return &handler{service, validator}
}

func (h *handler) Create(c *fiber.Ctx) error {
    payload := entity.Event{}

    if err := c.BodyParser(&payload); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.Event, error]{
        Message: err.Error(),
        Errors:  err,
    })
	}

  errV := h.validator.Validate(payload)
	if errV != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.Event, []*utils.ErrorResponse]{
			Message: "Validation Error",
			Errors:  errV,
		})
	}

    err := h.service.Create(&payload)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response[entity.Event, error]{
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
    criteria := model.FindAllEventCriteria{}
    if err := c.QueryParser(&criteria); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.FindAllEventCriteria, error]{
            Message: err.Error(),
            Errors:  err,
        })
    }

    errV := h.validator.Validate(criteria)
    if errV != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.FindAllEventCriteria, []*utils.ErrorResponse]{
            Message: "Validation Error",
            Errors:  errV,
        })
    }

    result, count, err := h.service.FindAll(criteria)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(model.Response[model.FindAllEventCriteria, error]{
            Message: err.Error(),
            Success: false,
        })
    }

    return c.JSON(model.ResponseWithOffset[[]*entity.Event, error]{
		Message: "Success",
		Total:   count,
		Limit:   criteria.Limit,
		Offset:  criteria.Offset,
		Data:    result,
		Success: true,
	})

}

func (h *handler) FindByID(c *fiber.Ctx) error {
id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.Event, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	result, err := h.service.FindByID(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response[interface{}, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.JSON(model.Response[*entity.Event, error]{
		Message: "Success",
		Data:    result,
		Success: true,
	})
}

func (h *handler) Update(c *fiber.Ctx) error {
 payload := entity.Event{}
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.Event, error]{
            Message: err.Error(),
            Errors:  err,
        })
    }

    errV := h.validator.Validate(payload)
    if errV != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[entity.Event, []*utils.ErrorResponse]{
            Message: "Validation Error",
            Errors:  errV,
        })
    }

    id, err := c.ParamsInt("id")
    if err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[interface{}, error]{
            Message: err.Error(),
            Success: false,
        })
    }

    err = h.service.Update(uint(id), &payload)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(model.Response[entity.Event, error]{
            Message: err.Error(),
            Success: false,
        })
    }

    return c.JSON(model.Response[entity.Event, error]{
        Message: "Success",
        Success: true,
    })
}

func (h *handler) Delete(c *fiber.Ctx) error {
    id, err := c.ParamsInt("id")
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[interface{}, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	err = h.service.Delete(uint(id))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response[interface{}, error]{
			Message: err.Error(),
			Success: false,
		})
	}

	return c.JSON(model.Response[interface{}, error]{
		Message: "Success",
		Success: true,
	})
}

