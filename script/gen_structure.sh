
if [ -z "$1" ]
then
    echo "Please provide a module name"
    exit 1
fi

go mod init $1

cat > go.mod <<EOF
module $1

go 1.20

require (
	github.com/go-playground/validator/v10 v10.22.0
	github.com/gofiber/fiber/v2 v2.52.5
	github.com/google/uuid v1.5.0
	github.com/joho/godotenv v1.5.1
	github.com/rs/zerolog v1.33.0
	github.com/samber/lo v1.47.0
	github.com/xuri/excelize/v2 v2.8.1
	gorm.io/driver/postgres v1.5.9
	gorm.io/gen v0.3.26
	gorm.io/gorm v1.25.11
	gorm.io/plugin/dbresolver v1.5.0
)

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/gabriel-vasile/mimetype v1.4.3 // indirect
	github.com/go-playground/locales v0.14.1 // indirect
	github.com/go-playground/universal-translator v0.18.1 // indirect
	github.com/go-sql-driver/mysql v1.7.0 // indirect
	github.com/jackc/pgpassfile v1.0.0 // indirect
	github.com/jackc/pgservicefile v0.0.0-20221227161230-091c0ba34f0a // indirect
	github.com/jackc/pgx/v5 v5.5.5 // indirect
	github.com/jackc/puddle/v2 v2.2.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/leodido/go-urn v1.4.0 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/richardlehane/mscfb v1.0.4 // indirect
	github.com/richardlehane/msoleps v1.0.3 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.55.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xuri/efp v0.0.0-20231025114914-d1ff6096ae53 // indirect
	github.com/xuri/nfp v0.0.0-20230919160717-d98342af3f05 // indirect
	golang.org/x/crypto v0.24.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/net v0.26.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/sys v0.24.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	gorm.io/datatypes v1.1.1-0.20230130040222-c43177d3cf8c // indirect
	gorm.io/driver/mysql v1.4.4 // indirect
	gorm.io/hints v1.1.0 // indirect
)

EOF

mkdir dto
mkdir dto/entity
mkdir dto/model
mkdir example
mkdir external
echo package external > external/routers.go
cat > external/routers.go <<EOF
package external

import (
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
	apiV1NoGuard.Get("/health", healthCheck)
}
EOF
mkdir internal
mkdir pkg
mkdir pkg/db
cat > pkg/db/db.go <<EOF
package db

import (
	"context"
	"fmt"
	"log"
	"$1/pkg/env"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gen"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type Querier interface {
	// SELECT * FROM @@table WHERE name = @name{{if role !=""}} AND role = @role{{end}}
	FilterWithNameAndRole(name, role string) ([]gen.T, error)
}

func GetDB() (*gorm.DB, error) {
	// Define the connection string
	//=> go to supabase project -> settings -> project settings -> configuration -> Database
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", env.Env().DB_HOST, env.Env().DB_USER, env.Env().DB_PASSWORD, env.Env().DB_NAME, env.Env().DB_PORT)

	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		DryRun: false,
	})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	// Access the underlying *sql.DB to configure connection pool settings
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database connection: %v", err)
	}

	// Set maximum number of open connections
	sqlDB.SetMaxOpenConns(10)

	// Set maximum number of idle connections
	sqlDB.SetMaxIdleConns(5)

	// Set maximum connection lifetime
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	log.Println("Successfully connected to Supabase!")
	return db, nil
}

func Migrate() {
	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	db, _ := GetDB()
	g.UseDB(db) // reuse your gorm db

	db.AutoMigrate(
	)

	g.ApplyInterface(func(Querier) { context.TODO() },
	)

	// Generate the code
	g.Execute()
}

func DeAllocate() {
	db, _ := GetDB()
	db.Exec("DEALLOCATE ALL")
}

EOF
mkdir pkg/env
cat > pkg/env/env.go <<EOF
package env

import (
	"os"
	"reflect"
	"sync"

	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

// This constant list should be in sync with the .env file
type env struct {
	APP_ENV      string
	APP_PORT     string
	PROJECT_PATH string
	TEMP_PATH    string

	DB_HOST     string
	DB_USER     string
	DB_PASSWORD string
	DB_NAME     string
	DB_PORT     string
}

// The single instance of env struct
var instance *env

// To make sure one goroutine access this at a time
var lock = &sync.Mutex{}

// Returns the same instance of env struct
func Env() *env {
	if instance == nil {
		lock.Lock()
		defer lock.Unlock()
		if instance == nil {
			log.Debug().Msg("No 'env' instance, create one and initialize its values")
			newInst := load()
			instance = &newInst
		}
	}
	return instance
}

// Loads and initialize values from .env file
func load() env {
	// Read .env
	if _, err := os.Stat("/configs/.env"); err == nil {
		log.Debug().Msg("Loading config from /configs/.env")
		err := godotenv.Load("/configs/.env")
		if err != nil {
			log.Error().Err(err).Msg("")
			log.Fatal().Msg("Error loading .env file from /configs/.env\n")
		}
	} else {
		log.Debug().Msg("Loading config from default location .")
		err := godotenv.Load()
		if err != nil {
			log.Debug().Msgf("Error loading .env from default location: %s\n", err)
		}
	}
	newInstance := env{}
	fields := reflect.VisibleFields(reflect.TypeOf(newInstance))
	ps := reflect.ValueOf(&newInstance)
	for i := 0; i < len(fields); i++ {
		fieldname := fields[i].Name
		value := os.Getenv(fieldname)
		ps.Elem().FieldByName(fieldname).SetString(value)
	}
	return newInstance
}
EOF

mkdir pkg/utils
cat > pkg/utils/utils.go <<EOF
package utils
EOF

mkdir query
touch .env
touch .env.test
touch .env.prod

cat > main.go <<EOF
package main
import (
	"$1/pkg/db"
	"$1/pkg/env"

	routers "$1/external"

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
	db.Exec("DEALLOCATE ALL")
	defer func() {
		dbInstance, _ := db.DB()
		_ = dbInstance.Close()
	}()

	routers.PublicRoutes(app, db)

	app.Listen(":" + port)

}
EOF

mkdir gen
mkdir gen/db_option
cat > gen/db_option/db_option.go <<EOF
package main

import (
	"$1/pkg/db"
)

func main() {
	db.DeAllocate()
}
EOF

cat > gen/db_migrate.go <<EOF
package main

import "$1/pkg/db"

func main() {
	db.Migrate()
}
EOF



go mod tidy