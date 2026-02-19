package chawk

import (
	endpoints "chawk/endpoints"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

type UserService struct {
	client *BlackboardClient
}

func (s *UserService) url(id string) string {
	if id == "" {
		return "/learn/api/public/v1/users"
	}
	return fmt.Sprintf("/learn/api/public/v1/users/%s", id)
}

var ErrUserNotFound = errors.New("user account doesn't exist")
var ErrUserExist = errors.New("user already exists")
var ErrMissingName = errors.New("missing first and last name")
var ErrInvalidUsername = errors.New("invalid username")

type User struct {
	ID                 string             `json:"id,omitempty"`
	UUID               string             `json:"uuid,omitempty"`
	ExternalID         string             `json:"externalId,omitempty"`
	DataSourceID       string             `json:"dataSourceId,omitempty"`
	UserName           string             `json:"userName"`
	StudentID          string             `json:"studentId,omitempty"`
	Password           string             `json:"password,omitempty"`
	EducationLevel     string             `json:"educationLevel,omitempty"`
	Gender             string             `json:"gender,omitempty"`
	Pronouns           string             `json:"pronouns,omitempty"`
	BirthDate          *time.Time         `json:"birthDate,omitempty"`
	InstitutionRoleIDs []string           `json:"institutionRoleIds"`
	SystemRoleIDs      []string           `json:"systemRoleIds,omitempty"`
	Availability       UserAvailability   `json:"availability,omitempty"`
	Name               Name               `json:"name,omitempty"`
	Job                Job                `json:"job,omitempty"`
	Contact            Contact            `json:"contact,omitempty"`
	Address            Address            `json:"address,omitempty"`
	Locale             Locale             `json:"locale,omitempty"`
	Avatar             Avatar             `json:"avatar,omitempty"`
	Pronunciation      string             `json:"pronunciation,omitempty"`
	PronunciationAudio PronunciationAudio `json:"pronunciationAudio,omitempty"`
}

type Name struct {
	Given                string `json:"given,omitempty"`
	Family               string `json:"family,omitempty"`
	Middle               string `json:"middle,omitempty"`
	Other                string `json:"other,omitempty"`
	Suffix               string `json:"suffix,omitempty"`
	Title                string `json:"title,omitempty"`
	PreferredDisplayName string `json:"preferredDisplayName,omitempty"`
}

type Job struct {
	Title      string `json:"title,omitempty"`
	Department string `json:"department,omitempty"`
	Company    string `json:"company,omitempty"`
}

type UserAvailability struct {
	Available string `json:"available"`
}

type Contact struct {
	HomePhone        string `json:"homePhone,omitempty"`
	MobilePhone      string `json:"mobilePhone,omitempty"`
	BusinessPhone    string `json:"businessPhone,omitempty"`
	BusinessFax      string `json:"businessFax,omitempty"`
	Email            string `json:"email,omitempty"`
	InstitutionEmail string `json:"institutionEmail,omitempty"`
	WebPage          string `json:"webPage,omitempty"`
}

type Address struct {
	Street1 string `json:"street1,omitempty"`
	Street2 string `json:"street2,omitempty"`
	City    string `json:"city,omitempty"`
	State   string `json:"state,omitempty"`
	ZipCode string `json:"zipCode,omitempty,omitempty"`
	Country string `json:"country,omitempty,omitempty"`
}

type Locale struct {
	ID             string  `json:"id,omitempty"`
	Calendar       *string `json:"calendar,omitempty"`
	FirstDayOfWeek *string `json:"firstDayOfWeek,omitempty"`
}

type Avatar struct {
	Source     string `json:"source,omitempty"`
	UploadID   string `json:"uploadId,omitempty"`
	ResourceID string `json:"resourceId,omitempty"`
}

type PronunciationAudio struct {
	UploadID string `json:"uploadId,omitempty"`
}

type UserUpdate struct {
	Contact            *ContactUpdate    `json:"contact,omitempty"`
	Name               *NameUpdate       `json:"name,omitempty"`
	Password           *string           `json:"password,omitempty"`
	InstitutionRoleIDs []string          `json:"institutionRoleIds,omitempty"`
	Availability       *UserAvailability `json:"availability,omitempty"`
}

type ContactUpdate struct {
	Email            *string `json:"email,omitempty"`
	InstitutionEmail *string `json:"institutionEmail,omitempty"`
}

type NameUpdate struct {
	Given  *string `json:"given,omitempty"`
	Family *string `json:"family,omitempty"`
}

func (us *UserService) CreateUser(ctx context.Context, username string, fName string, lName string, email string, password string) error {
	username = strings.TrimSpace(username)
	fName = strings.TrimSpace(fName)
	lName = strings.TrimSpace(lName)
	//email = strings.TrimSpace(email)

	email = OptionalString(email)

	password = strings.TrimSpace(password)

	if username == "" || fName == "" || lName == "" {
		return errors.New("missing parameters: username, fName, lName")
	}

	data := User{
		UserName: username,
		Password: password,
		Availability: UserAvailability{
			Available: "Yes",
		},
		Name: Name{
			Given:                fName,
			Family:               lName,
			PreferredDisplayName: "GivenName",
		},
		Contact: Contact{
			Email: email,
		},
	}

	fmt.Printf("%v\n", data)

	url := endpoints.Users.Create()
	resp, err := us.client.Post(ctx, url, data)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 201:
		return nil
		//client.Logger.Info(fmt.Sprintf("User %s was created successfully", username))
	case 403:
		return ErrInsufficientPrivileges
		//client.Logger.Error("Insufficient privileges to create a new user")
	case 409:
		return ErrUserExist
		//client.Logger.Error(fmt.Sprintf("User with ID %s already exists", username))
	case 400:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		if strings.Contains(string(body), "A database error occurred") {
			// Treat this as a "Potential Conflict/Already Exists"
			// When testing, a 409 was not being retunred
			return ErrUserExist
		}
		//client.Logger.Error(fmt.Sprintf("Error creating user: %s", string(body)))
		return fmt.Errorf("An error occurred while creating the new user: %s", string(body))
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		//client.Logger.Error(fmt.Sprintf("Unexpected status: %d, response: %s", resp.StatusCode, string(body)))
	}

}

func (s *UserService) DoesUserExist(ctx context.Context, username string) (bool, error) {
	url := endpoints.Users.GetByUsername(username)

	resp, err := s.client.Get(ctx, url)
	if err != nil {
		return false, fmt.Errorf("check user failed: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK: // 200
		return true, nil

	case http.StatusNotFound: // 404
		// This is the ONLY time we can confidently say "User does not exist"
		return false, nil

	default:
		// Any other code (401, 403, 500) is an error, not a "false"
		return false, fmt.Errorf("unexpected API status: %d", resp.StatusCode)
	}
}

// GetUserObject fetches user data and maps to User struct
func (us *UserService) GetUserByUsername(ctx context.Context, username string) (*User, error) {
	username = strings.TrimSpace(username)
	//exists := us.DoesUserExist(username)
	//if !exists {
	//	return nil, fmt.Errorf("someing?")
	//}
	url := endpoints.Users.GetByUsername(username)
	resp, err := us.client.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get user: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	if err != nil {
		return nil, err
	}

	var u User
	if err := json.Unmarshal(body, &u); err != nil {
		return nil, err
	}
	return &u, nil
}

func (us *UserService) Update(ctx context.Context, username string, update UserUpdate) error {
	username = strings.TrimSpace(username)
	if username == "" {
		return ErrInvalidUsername
	}

	url := endpoints.Users.GetByUsername(username)

	resp, err := us.client.Patch(ctx, url, update)
	if err != nil {
		return err // transport / marshal / network
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent:
		return nil

	case http.StatusNotFound:
		return ErrUserNotFound

	case http.StatusBadRequest:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("invalid user update: %s", string(body))

	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf(
			"update user failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}
}

func (us *UserService) UpdatePassword(ctx context.Context, username, newPassword string) error {
	password := strings.TrimSpace(newPassword)
	if password == "" {
		return errors.New("new password cannot be empty")
	}

	return us.Update(ctx, username, UserUpdate{
		Password: ToPtr(password),
	})
}

func (us *UserService) UpdateEmail(ctx context.Context, username, newEmail string) error {
	email := strings.TrimSpace(newEmail)

	return us.Update(ctx, username, UserUpdate{
		Contact: &ContactUpdate{
			Email: ToPtr(email),
		},
	})
}

// func (uc *BlackboardClient) UpdateInstitutionEmail(username, newEmail string) error {
// 	return uc.updateUser(username, map[string]interface{}{"contact": map[string]string{"institutionEmail": strings.TrimSpace(newEmail)}}, "institution email changed")
// }

func (us *UserService) UpdateInstitutionEmail(ctx context.Context, username, newEmail string) error {
	email := strings.TrimSpace(newEmail)

	return us.Update(ctx, username, UserUpdate{
		Contact: &ContactUpdate{
			InstitutionEmail: ToPtr(email),
		},
	})
}

// func (uc *BlackboardClient) UpdateName(username, fName, lName string) error {
// 	if fName == "" && lName == "" {
// 		return errors.New("first and last name not provided")
// 	}
// 	data := map[string]interface{}{
// 		"name": map[string]string{
// 			"given":  fName,
// 			"family": lName,
// 		},
// 	}
// 	return uc.updateUser(username, data, "name changed")
// }

func (us *UserService) UpdateName(ctx context.Context, username, fName, lName string) error {
	fName = strings.TrimSpace(fName)
	lName = strings.TrimSpace(lName)

	if fName == "" && lName == "" {
		return ErrMissingName
	}

	update := &NameUpdate{}
	if fName != "" {
		update.Given = &fName
	}
	if lName != "" {
		update.Family = &lName
	}

	return us.Update(ctx, username, UserUpdate{Name: update})
}

func (us *UserService) AddInstitutionRoles(ctx context.Context, username string, roles []string) error {
	if len(roles) == 0 {
		return errors.New("no roles provided")
	}

	return us.Update(ctx, username, UserUpdate{
		InstitutionRoleIDs: roles,
	})
}

func (us *UserService) UpdateUserAvailability(ctx context.Context, username string, availability string) error {
	//availability = strings.TrimSpace(availability)

	switch availability {
	case "Yes", "No", "Disabled":
		// All good
	default:
		return errors.New("availability must be 'Disabled', 'Yes', or 'No'")
	}

	// if !uc.DoesUserExist(username) {
	// 	return ErrUserNotFound
	// }

	return us.Update(ctx, username, UserUpdate{
		Availability: &UserAvailability{Available: availability},
	})
}
