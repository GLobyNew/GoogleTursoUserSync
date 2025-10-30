package google

import (
	"context"
	"encoding/json"
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

func TestGetTgIDFromUser(t *testing.T) {
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

	testEmail := os.Getenv("GOOGLE_TEST_USER_EMAIL")
	if testEmail == "" {
		t.Fatalf("GOOGLE_TEST_USER_EMAIL is not set")
	}

	user, err := userService.GetUserByEmail(context.Background(), testEmail, "MessengerInfo")
	if err != nil {
		t.Fatalf("Error getting user by email: %v", err)
	}

	if user.CustomSchemas == nil || user.CustomSchemas["MessengerInfo"] == nil {
		t.Fatalf("CustomSchemas or MessengerInfo is nil")
	}
	mesInfo, err := user.CustomSchemas["MessengerInfo"].MarshalJSON()
	if err != nil {
		t.Fatalf("Error marshaling MessengerInfo: %v", err)
	}

	var cusFields CustomFieldsTgID
	err = json.Unmarshal(mesInfo, &cusFields)
	if err != nil {
		t.Fatalf("Error unmarshaling MessengerInfo: %v", err)
	}

	t.Logf("tgID: %v", cusFields.TgID)

}

func TestGetUserByEmail(t *testing.T) {
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

	testEmail := os.Getenv("GOOGLE_TEST_USER_EMAIL")
	if testEmail == "" {
		t.Fatalf("GOOGLE_TEST_USER_EMAIL is not set")
	}

	user, err := userService.GetUserByEmail(context.Background(), testEmail, "MessengerInfo")
	if err != nil {
		t.Fatalf("Error getting user by email: %v", err)
	}

	if user.PrimaryEmail != testEmail {
		t.Fatalf("Expected email %s, got %s", testEmail, user.PrimaryEmail)
	}

	t.Logf("Retrieved user: %s (%s)\n", user.PrimaryEmail, user.CustomSchemas["MessengerInfo"])
	t.Logf("Full user data: %+v\n", user)

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

	users, err := userService.GetAllUsers(context.Background(), "")
	if err != nil {
		t.Fatalf("Error getting all users: %v", err)
	}

	if len(users) == 0 {
		t.Fatalf("Expected to find users, but got none")
	}

}

func TestPrintAllUsersWithoutTgID(t *testing.T) {
	godotenv.Load()
	credentials, err := os.ReadFile("test-credentials.json")
	if err != nil {
		t.Fatalf("Unable to read client secret file: %v", err)
	}
	domain := os.Getenv("GOOGLE_TEST_DOMAIN")
	if domain == "" {
		t.Fatalf("GOOGLE_TEST_DOMAIN is not set")
	}
	t.Logf("\n")

	googleAdminService, err := newTestGoogleAdminService(context.Background(), credentials)
	if err != nil {
		t.Fatalf("Unable to retrieve Google Admin service: %v", err)
	}
	userService := NewUserService(googleAdminService, domain)

	err = userService.PrintAllUsersWithoutTgID(context.Background(), "MessengerInfo")
	if err != nil {
		t.Fatalf("Error printing all users without tgID: %v", err)
	}

}
