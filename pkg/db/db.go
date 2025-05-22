package db

import (
	"context"
	"event-management-system/dto/entity"
	"event-management-system/pkg/env"
	"fmt"
	"log"
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
	//=> pgbouncer=true is disable stmt cache
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s pgbouncer=true", env.Env().DB_HOST, env.Env().DB_USER, env.Env().DB_PASSWORD, env.Env().DB_NAME, env.Env().DB_PORT)
	// Open a connection to the database
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger:      logger.Default.LogMode(logger.Info),
		DryRun:      false,
		PrepareStmt: false,
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
	db, _ := GetDB()

	g := gen.NewGenerator(gen.Config{
		OutPath: "./query",
		Mode:    gen.WithoutContext | gen.WithDefaultQuery | gen.WithQueryInterface, // generate mode
	})

	g.UseDB(db) // reuse your gorm db

	db.AutoMigrate(
		&entity.User{},
		&entity.Event{},
		&entity.EventParticipant{},
	)

	// Generate Type Safe API with Dynamic SQL defined on Querier interface for `model.User` and `model.Company`
	// g.ApplyInterface(func(Querier) {},
	g.ApplyInterface(func(Querier) { context.TODO() },
		&entity.User{},
		&entity.Event{},
		&entity.EventParticipant{},
	)

	// Generate the code
	g.Execute()
}

func DeAllocate() {
	db, _ := GetDB()
	db.Exec("DEALLOCATE ALL")
}
