package google

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func NewGoogleAdminService(ctx context.Context, credentials []byte) (*admin.Service, error) {
	var info struct {
		Type string `json:"type"`
	}
	if err := json.Unmarshal(credentials, &info); err != nil {
		return nil, fmt.Errorf("unable to parse credentials: %w", err)
	}

	switch info.Type {
	case "service_account":
		// Domain-wide delegation requires impersonating an admin user.
		adminEmail := os.Getenv("GOOGLE_ADMIN_EMAIL")
		if adminEmail == "" {
			return nil, fmt.Errorf("GOOGLE_ADMIN_EMAIL is not set for service account impersonation")
		}
		jwtConfig, err := google.JWTConfigFromJSON(credentials, admin.AdminDirectoryUserReadonlyScope)
		if err != nil {
			return nil, err
		}
		jwtConfig.Subject = adminEmail
		client := jwtConfig.Client(ctx)
		return admin.NewService(ctx, option.WithHTTPClient(client))
	default:
		googleConfig, err := google.ConfigFromJSON(credentials, admin.AdminDirectoryUserReadonlyScope)
		if err != nil {
			return nil, err
		}
		client := getClient(googleConfig)
		return admin.NewService(ctx, option.WithHTTPClient(client))
	}
}
