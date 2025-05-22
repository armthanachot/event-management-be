package main

import (
	routers "event-management-system/external"
	"event-management-system/pkg/db"
	"event-management-system/pkg/env"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

func main() {
	port := env.Env().APP_PORT

	app := fiber.New(fiber.Config{
		// Override default error handler
		ErrorHandler: func(ctx *fiber.Ctx, err error) error {
			// Retrieve the custom status code if it's an fiber.*Error
			if _, ok := err.(*fiber.Error); !ok {
				if env.Env().APP_ENV == "develop" {
					// microsoft.NewAppinsights().Error("Error log ::" + err.Error())
				}
				return ctx.Status(fiber.StatusInternalServerError).SendString("Internal Server Error")
			}

			// Return from handler
			return nil
		},
		BodyLimit: 100 * 1024 * 1024,
	})

	app.Use(recover.New())
	app.Use(cors.New())

	app.Use(requestid.New())

	db, _ := db.GetDB()
	// db.Exec("DEALLOCATE ALL")
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	routers.PublicRoutes(app, db)
	app.Listen(":" + port)

}
