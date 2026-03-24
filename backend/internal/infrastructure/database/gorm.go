package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// InitGORM connects to Postgres using Supabase connection string
func InitGORM() {
	// e.g. "host=db.xxxxx.supabase.co user=postgres password=xxxxx dbname=postgres port=5432 sslmode=require"
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Println("Warning: DATABASE_URL not set, GORM initialization deferred or failed")
		return
	}

	var err error
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database using GORM: %v", err)
	}

	log.Println("Database connection established via GORM")
}
