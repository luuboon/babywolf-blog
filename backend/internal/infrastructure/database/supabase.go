package database

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nedpals/supabase-go"
)

var SupabaseClient *supabase.Client

func InitSupabase() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: Error loading .env file, relying on system environment variables")
	}

	supabaseURL := os.Getenv("SUPABASE_URL")
	supabaseKey := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")

	if supabaseURL == "" || supabaseKey == "" {
		log.Fatal("SUPABASE_URL and SUPABASE_SERVICE_ROLE_KEY must be set in .env")
	}

	SupabaseClient = supabase.CreateClient(supabaseURL, supabaseKey)
	log.Println("Supabase client initialized")
}
