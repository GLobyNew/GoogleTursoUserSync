package google

import (
	"context"
	"encoding/json"
	"fmt"
	"log"

	admin "google.golang.org/api/admin/directory/v1"
)

type GoogleUser struct {
	PrimaryEmail string
	tgID         int64
}

type UserService struct {
	adminSrv *admin.Service
	domain   string
}

type CustomFieldsTgID struct {
	TgID int64 `json:"tgID"`
}

func NewUserService(adminSrv *admin.Service, domain string) *UserService {
	return &UserService{
		adminSrv: adminSrv,
		domain:   domain,
	}
}

func (u *UserService) GetUserByEmail(ctx context.Context, email, customFieldMask string) (*admin.User, error) {
	user, err := u.adminSrv.Users.Get(email).Projection("custom").CustomFieldMask(customFieldMask).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve user %s: %v", email, err)
	}
	return user, nil
}

func parseTgIDFromUser(user *admin.User) (int64, error) {
	if user.CustomSchemas == nil || user.CustomSchemas["MessengerInfo"] == nil {
		return 0, fmt.Errorf("CustomSchemas or MessengerInfo is nil")
	}
	mesInfo, err := user.CustomSchemas["MessengerInfo"].MarshalJSON()
	if err != nil {
		return 0, fmt.Errorf("Error marshaling MessengerInfo: %v", err)
	}

	var cusFields CustomFieldsTgID
	err = json.Unmarshal(mesInfo, &cusFields)
	if err != nil {
		return 0, fmt.Errorf("Error unmarshaling MessengerInfo: %v", err)
	}
	return cusFields.TgID, nil
}

func (u *UserService) GetAllUsers(ctx context.Context, customFieldMask string) ([]GoogleUser, error) {
	var r *admin.Users
	var err error
	if customFieldMask == "" {
		r, err = u.adminSrv.Users.List().Domain(u.domain).OrderBy("email").Do()
	} else {
		r, err = u.adminSrv.Users.List().Domain(u.domain).Projection("custom").CustomFieldMask(customFieldMask).OrderBy("email").Do()
	}
	if err != nil {
		log.Fatalf("Unable to retrieve users in domain: %v", err)
	}

	googleUsers := make([]GoogleUser, len(r.Users))
	for i, user := range r.Users {
		tgID, err := parseTgIDFromUser(user)
		if err != nil {
			tgID = 0
		}
		googleUsers[i] = GoogleUser{
			PrimaryEmail: user.PrimaryEmail,
			tgID:         tgID,
		}
	}

	return googleUsers, nil
}

func (u *UserService) PrintAllUsers(ctx context.Context, customFieldMask string) error {
	users, err := u.GetAllUsers(ctx, customFieldMask)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		fmt.Print("No users found.\n")
	} else {
		fmt.Print("Users:\n")
		for _, u := range users {
			fmt.Printf("%v (%v)\n", u.PrimaryEmail, u.tgID)
		}
	}
	return nil

}

func (u *UserService) PrintAllUsersWithoutTgID(ctx context.Context, customFieldMask string) error {
	users, err := u.GetAllUsers(ctx, customFieldMask)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		fmt.Print("No users found.\n")
	} else {
		fmt.Print("Users without tgID:\n")
		for _, u := range users {
			if u.tgID == 0 {
				fmt.Printf("%v\n", u.PrimaryEmail)
			}
		}
	}
	return nil
}
