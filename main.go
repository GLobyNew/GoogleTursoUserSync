package main

import (
	"context"
	"fmt"
	"log"
	"os"

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

	err = userService.PrintAllUsers(ctx, customFieldMask)
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}

}
