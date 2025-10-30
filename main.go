package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/GLobyNew/GoogleTursoUserSync/internal/database"
	google "github.com/GLobyNew/GoogleTursoUserSync/internal/google"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	domain := os.Getenv("GOOGLE_DOMAIN")
	if domain == "" {
		log.Fatalf("GOOGLE_DOMAIN is not set")
	}

	customFieldMask := os.Getenv("GOOGLE_CUSTOM_FIELD_MASK")
	fmt.Printf("Using custom field mask: %s\n", customFieldMask)

	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	googleAdminService, err := google.NewGoogleAdminService(ctx, credentials)
	if err != nil {
		log.Fatalf("Unable to retrieve Google Admin service: %v", err)
	}
	userService := google.NewUserService(googleAdminService, domain)
	googleUsers, err := userService.GetAllUsers(ctx, customFieldMask)
	if err != nil {
		log.Fatalf("Error getting all users: %v", err)
	}

	db, err := database.NewDatabaseConnection(databaseURL)
	if err != nil {
		log.Fatalf("Error connecting to database: %v", err)
	}
	defer db.Close()
	log.Printf("Successfully connected to database.\n")

	tursoUsers, err := db.GetAllUsers()
	if err != nil {
		log.Fatalf("Error retrieving users from database: %v", err)
	}

	log.Printf("Total users retrieved from Google: %d\n", len(googleUsers))
	log.Printf("Total users retrieved from Turso database: %d\n", len(tursoUsers))

	updated, cretaed, upToDate := db.SyncUsers(googleUsers, tursoUsers)
	log.Printf("Sync complete. Updated: %d, Created: %d, Up-to-date: %d\n", updated, cretaed, upToDate)

}
