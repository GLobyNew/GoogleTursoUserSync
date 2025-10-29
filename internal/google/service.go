package google

import (
	"context"

	"golang.org/x/oauth2/google"
	admin "google.golang.org/api/admin/directory/v1"
	"google.golang.org/api/option"
)

func NewGoogleAdminService(ctx context.Context, credentials []byte) (*admin.Service, error) {

	googleConfig, err := google.ConfigFromJSON(credentials, admin.AdminDirectoryUserReadonlyScope)
	if err != nil {
		return nil, err
	}

	client := getClient(googleConfig)

	srv, err := admin.NewService(ctx, option.WithHTTPClient(client))
	if err != nil {
		return nil, err
	}

	return srv, nil

}
