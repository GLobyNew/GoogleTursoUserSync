package main

import (
	"context"
	"log"
	"os"

	google "github.com/GLobyNew/GoogleTursoUserSync/internal/google"
	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	godotenv.Load()

	DatabaseURL := os.Getenv("DATABASE_URL")
	if DatabaseURL == "" {
		log.Fatalf("DATABASE_URL is not set")
	}

	// config := &config.Config{
	// 	DatabaseURL: DatabaseURL,
	// }

	credentials, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	googleAdminService, err := google.NewGoogleAdminService(ctx, credentials)
	if err != nil {
		log.Fatalf("Unable to retrieve Google Admin service: %v", err)
	}
	userService := google.NewUserService(googleAdminService)

	err = userService.ListAllUsers(ctx)
	if err != nil {
		log.Fatalf("Error listing users: %v", err)
	}

}
