
if [ -z "$1" ]
then
    echo "Please provide a module name"
    exit 1
fi

pwd

# Read the go.mod file and extract the module name
module_name=$(grep '^module' go.mod | awk '{print $2}')

# Check if module name was found
if [ -n "$module_name" ]; then
  # Write the module name to a new file
  echo "Module name '$module_name' copied to module_name.txt"
else
  echo "Module name not found in go.mod"
  exit 1
fi

public_internal_module_name=$(echo "$1" | awk '{print toupper(substr($0, 1, 1)) tolower(substr($0, 2))}')
criteria_str="Criteria"

cat > dto/entity/$1.go <<EOF
package entity

import (
    "gorm.io/gorm"
	"$module_name/dto/model"
)

type $public_internal_module_name struct {
    gorm.Model
}

func (e *$public_internal_module_name) ToModel() model.$public_internal_module_name {
    return model.$public_internal_module_name{
    }
}

EOF

cat > dto/model/$1.go <<EOF
package model

type $public_internal_module_name struct {
}

type FindAll$public_internal_module_name$criteria_str struct {
    Limit  int   \`json:"limit" validate:"gte=0"\`
    Offset int   \`json:"offset" validate:"gte=0"\`
    Search string \`json:"search"\`
}

type CreateUpdate$public_internal_module_name struct {
}

EOF

cd internal


mkdir -p $1
cd $1
cat > repository.go <<EOF
// you must be migrate entity first (put entity in the pkg/db/db.go in function Migrate under db.AutoMigrate and g.ApplyInterface)
// after that run script file script/migrate.sh
package $1

import (
	"$module_name/dto/entity"
	"$module_name/dto/model"
	"$module_name/pkg/utils"
	"$module_name/query"

	"gorm.io/gorm"
)

type (
    Repository interface {
        Create(payload *entity.$public_internal_module_name) error
        FindAll(criteria model.FindAll$public_internal_module_name$criteria_str) ([]*entity.$public_internal_module_name, int64, error)
        FindByID(id uint) (*entity.$public_internal_module_name, error)
        Update(id uint, payload *entity.$public_internal_module_name) error
        Delete(id uint) error
    }
    repository struct {
        db *gorm.DB
    }
)

func NewRepository(db *gorm.DB) Repository {
    query.SetDefault(db)
    return &repository{db}
}

func (r *repository) Create(payload *entity.$public_internal_module_name) error {
    instance := query.$public_internal_module_name
    err := instance.Create(payload)
    if err != nil {
        return err
    }
    return nil
}

func (r *repository) FindAll(criteria model.FindAll$public_internal_module_name$criteria_str) ([]*entity.$public_internal_module_name, int64, error) {
    instance := query.$public_internal_module_name
    limit := criteria.Limit
    offset := criteria.Offset
    search := "%" + criteria.Search + "%"
    result, err := instance.Limit(limit).Offset(offset).Find()
    if err != nil {
        return nil, 0, err
    }
    count, err := instance.Count()
    if err != nil {
        return nil, 0, err
    }
    return  result, count, nil
}

func (r *repository) FindByID(id uint) (*entity.$public_internal_module_name, error) {
    instance := query.$public_internal_module_name
    result, err := instance.Where(instance.ID.Eq(id)).First()
    if err != nil {
        return nil, err
    }
    return result, nil
}

func (r *repository) Update(id uint, payload *entity.$public_internal_module_name) error {
    instance := query.$public_internal_module_name
    _, err := instance.Where(instance.ID.Eq(id)).Updates(payload)
    if err != nil {
        return err
    }
    return nil
}

func (r *repository) Delete(id uint) error {
    instance := query.$public_internal_module_name
    // remove Unscoped if you want to use soft delete
    _, err := instance.Unscoped().Where(instance.ID.Eq(id)).Delete()
    if err != nil {
        return err
    }
    return nil
}
EOF

cat > service.go <<EOF
package $1

import (
    "$module_name/dto/model"
)

type (
    Service interface {
        Create(payload model.CreateUpdate$public_internal_module_name) error
        FindAll(criteria model.FindAll$public_internal_module_name$criteria_str) ([]model.$public_internal_module_name, int64, error)
        FindByID(id uint) (model.$public_internal_module_name, error)
        Update(id uint, payload model.CreateUpdate$public_internal_module_name) error
        Delete(id uint) error
    }


    service struct {
        repository Repository
    }
)

func NewService(repository Repository) Service {
    return &service{repository}
}

func (s *service) Create(payload model.CreateUpdate$public_internal_module_name) error {
    err := s.repository.Create(ToEntity(payload))
    if err != nil {
        return err
    }
    return nil
}

func (s *service) FindAll(criteria model.FindAll$public_internal_module_name$criteria_str) ([]model.$public_internal_module_name, int64, error) {
    result, count, err := s.repository.FindAll(criteria)
    if err != nil {
        return nil, 0, err
    }

    resp := []model.$public_internal_module_name{}
    for _, v := range result {
        resp = append(resp, v.ToModel())
    }

    return resp, count, nil
}

func (s *service) FindByID(id uint) (model.$public_internal_module_name, error) {
    result, err := s.repository.FindByID(id)
    if err != nil {
        return model.$public_internal_module_name{}, err
    }

    return result.ToModel(), nil
}

func (s *service) Update(id uint, payload model.CreateUpdate$public_internal_module_name) error {
    err := s.repository.Update(id, ToEntity(payload))
    if err != nil {
        return err
    }
    return nil
}

func (s *service) Delete(id uint) error {
    err := s.repository.Delete(id)
    if err != nil {
        return err
    }
    return nil
}

EOF

cat > handler.go <<EOF
package $1
import (
	"github.com/gofiber/fiber/v2"
	"$module_name/pkg/utils"
    "$module_name/dto/model"

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
    payload := model.CreateUpdate$public_internal_module_name{}

    if err := c.BodyParser(&payload); err != nil {
    return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.$public_internal_module_name, error]{
        Message: err.Error(),
        Errors:  err,
    })
	}

  errV := h.validator.Validate(payload)
	if errV != nil {
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.$public_internal_module_name, []*utils.ErrorResponse]{
			Message: "Validation Error",
			Errors:  errV,
		})
	}

    err := h.service.Create(payload)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(model.Response[model.$public_internal_module_name, error]{
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
    criteria := model.FindAll$public_internal_module_name$criteria_str{}
    if err := c.QueryParser(&criteria); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.FindAll$public_internal_module_name$criteria_str, error]{
            Message: err.Error(),
            Errors:  err,
        })
    }

    errV := h.validator.Validate(criteria)
    if errV != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.FindAll$public_internal_module_name$criteria_str, []*utils.ErrorResponse]{
            Message: "Validation Error",
            Errors:  errV,
        })
    }

    result, count, err := h.service.FindAll(criteria)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(model.Response[model.FindAll$public_internal_module_name$criteria_str, error]{
            Message: err.Error(),
            Success: false,
        })
    }

    return c.JSON(model.ResponseWithOffset[[]model.$public_internal_module_name, error]{
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
		return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.$public_internal_module_name, error]{
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

	return c.JSON(model.Response[model.$public_internal_module_name, error]{
		Message: "Success",
		Data:    result,
		Success: true,
	})
}

func (h *handler) Update(c *fiber.Ctx) error {
 payload := model.CreateUpdate$public_internal_module_name{}
    if err := c.BodyParser(&payload); err != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.CreateUpdate$public_internal_module_name, error]{
            Message: err.Error(),
            Errors:  err,
        })
    }

    errV := h.validator.Validate(payload)
    if errV != nil {
        return c.Status(fiber.StatusBadRequest).JSON(model.Response[model.CreateUpdate$public_internal_module_name, []*utils.ErrorResponse]{
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

    err = h.service.Update(uint(id), payload)
    if err != nil {
        return c.Status(fiber.StatusInternalServerError).JSON(model.Response[model.CreateUpdate$public_internal_module_name, error]{
            Message: err.Error(),
            Success: false,
        })
    }

    return c.JSON(model.Response[model.CreateUpdate$public_internal_module_name, error]{
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

EOF

cat > router.go <<EOF
package $1
import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"

	"$module_name/pkg/utils"

)

func Router(r fiber.Router, db *gorm.DB) {
	repository := NewRepository(db)
	service := NewService(repository)

	validator := utils.NewValidator()

	handler := NewHandler(service, validator)

	groupRoute := r.Group("/$1")
	groupRoute.Post("/", handler.Create)
    groupRoute.Get("/", handler.FindAll)
    groupRoute.Get("/:id", handler.FindByID)
    groupRoute.Patch("/:id", handler.Update)
    groupRoute.Delete("/:id", handler.Delete)
}
EOF


# echo "package $1" > transform.go
cat > transform.go <<EOF
package $1

import (
    "$module_name/dto/entity"
    "$module_name/dto/model"
)

func ToEntity(payload model.CreateUpdate$public_internal_module_name) *entity.$public_internal_module_name {
    return &entity.$public_internal_module_name{
    }
}

func ToModel(payload *entity.$public_internal_module_name) model.$public_internal_module_name {
    return model.$public_internal_module_name{
    }
}

EOF


# prepared CRUD