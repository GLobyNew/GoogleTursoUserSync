package google

import (
	"context"
	"fmt"
	"log"

	admin "google.golang.org/api/admin/directory/v1"
)

type UserService struct {
	adminSrv *admin.Service
	domain   string
}

func NewUserService(adminSrv *admin.Service, domain string) *UserService {
	return &UserService{
		adminSrv: adminSrv,
		domain:   domain,
	}
}

func (u *UserService) GetAllUsers(ctx context.Context) ([]*admin.User, error) {
	r, err := u.adminSrv.Users.List().Domain(u.domain).OrderBy("email").Do()
	if err != nil {
		log.Fatalf("Unable to retrieve users in domain: %v", err)
	}
	return r.Users, nil
}

func (u *UserService) ListAllUsers(ctx context.Context) error {
	users, err := u.GetAllUsers(ctx)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		fmt.Print("No users found.\n")
	} else {
		fmt.Print("Users:\n")
		for _, u := range users {
			fmt.Printf("%s (%s)\n", u.PrimaryEmail, u.Name.FullName)
		}
	}
	return nil

}
