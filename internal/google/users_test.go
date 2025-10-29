package google

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func newTestGoogleAdminService(ctx context.Context, credentials []byte) (*admin.Service, error) {

	googleConfig, err := google.ConfigFromJSON(credentials, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return nil, err
	}

	client := getTestClient(googleConfig)

	srv, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil

}

func TestGetAllUsers(t *testing.T) {
	godotenv.Load()
	credentials, err := os.ReadFile("test-credentials.json")
	if err != nil {
		t.Fatalf("Unable to read client secret file: %v", err)
	}
	domain := os.Getenv("GOOGLE_TEST_DOMAIN")
	if domain == "" {
		t.Fatalf("GOOGLE_TEST_DOMAIN is not set")
	}

	googleAdminService, err := newTestGoogleAdminService(context.Background(), credentials)
	if err != nil {
		t.Fatalf("Unable to retrieve Google Admin service: %v", err)
	}
	userService := NewUserService(googleAdminService, domain)

	users, err := userService.GetAllUsers(context.Background())
	if err != nil {
		t.Fatalf("Error getting all users: %v", err)
	}

	if len(users) == 0 {
		t.Fatalf("Expected to find users, but got none")
	}

}
