package chawk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	endpoints "github.com/sugarvoid/chawk/endpoints"
)

type CourseService struct {
	client *BlackboardClient
}

var ErrCourseNotFound = errors.New("course doesn't exist")
var ErrCourseExist = errors.New("course already exists")
var ErrUserAlreadyEnrolled = errors.New("user already enrolled in course, use UpdateMembership() if needed")
var ErrInvalidRole = errors.New("invalid course role")

type CourseAvailability struct {
	Available string         `json:"available"`
	Duration  CourseDuration `json:"duration,omitzero"`
}

type CourseDuration struct {
	Type      string    `json:"type,omitempty"`
	Start     time.Time `json:"start,omitzero"`
	End       time.Time `json:"end,omitzero"`
	DaysOfUse int       `json:"daysOfUse,omitempty"`
}

type Enrollment struct {
	Type       string    `json:"type,omitempty"`
	Start      time.Time `json:"start,omitzero"`
	End        time.Time `json:"end,omitzero"`
	AccessCode string    `json:"accessCode,omitempty"`
}

type CourseLocale struct {
	ID    string `json:"id,omitempty"`
	Force bool   `json:"force,omitempty"`
}

type CopyHistory struct {
	UUID string `json:"uuid,omitempty"`
}

type Course struct {
	// Required for creation
	CourseID string `json:"courseId"`
	Name     string `json:"name"`
	TermID   string `json:"termId"`

	// Optional for creation
	Description    string             `json:"description,omitempty"`
	Availability   CourseAvailability `json:"availability"`
	Organization   bool               `json:"organization,omitempty"`
	UltraStatus    string             `json:"ultraStatus,omitempty"`
	AllowGuests    bool               `json:"allowGuests,omitempty"`
	AllowObservers bool               `json:"allowObservers"`
	ClosedComplete bool               `json:"closedComplete"`
	Enrollment     Enrollment         `json:"enrollment,omitzero"`
	Locale         CourseLocale       `json:"locale,omitzero"`

	// Read-only (set by server, ignored on create)
	ID           string     `json:"id,omitempty"`
	UUID         string     `json:"uuid,omitempty"`
	Created      *time.Time `json:"created,omitempty"`
	Modified     *time.Time `json:"modified,omitempty"`
	ExternalID   string     `json:"externalId,omitempty"`
	DataSourceID string     `json:"dataSourceId,omitempty"`

	// Stuff Populated on GET'ing
	HasChildren       bool          `json:"hasChildren,omitempty"`
	ParentID          string        `json:"parentId,omitempty"`
	ExternalAccessURL string        `json:"externalAccessUrl,omitempty"`
	GuestAccessURL    string        `json:"guestAccessUrl,omitempty"`
	CopyHistory       []CopyHistory `json:"copyHistory"`
	IsChild           bool          `json:"-"`
}

type CourseUpdateRequest struct {
	Name         *string             `json:"name,omitempty"`
	TermID       *string             `json:"termId,omitempty"`
	Availability *CourseAvailability `json:"availability,omitempty"`
	DataSourceID *string             `json:"dataSourceId,omitempty"`
	Description  *string             `json:"description,omitempty"`
}

type CourseCreateRequest struct {
	Name         *string             `json:"name,omitempty"`
	ExternalID   *string             `json:"termId,omitempty"`
	TermID       *string             `json:"externalId,omitempty"`
	Availability *CourseAvailability `json:"availability,omitempty"`
	DataSourceID *string             `json:"dataSourceId,omitempty"`
	Description  *string             `json:"description,omitempty"`
}

type EnrollmentRequest struct {
	ChildCourseID *string                 `json:"childCourseId,omitempty"`
	DataSourceID  *string                 `json:"dataSourceId,omitempty"`
	Availability  *MembershipAvailability `json:"availability,omitempty"`
	CourseRoleID  *string                 `json:"courseRoleId,omitempty"`
	DisplayOrder  *int                    `json:"displayOrder,omitempty"`
}

type MembershipAvailability struct {
	Available *string `json:"available,omitempty"`
}

//type AvailabilityStatus string

const (
	AvailabilityYes      string = "Yes"
	AvailabilityNo       string = "No"
	AvailabilityDisabled string = "Disabled"

	RoleStudent    string = "Student"
	RoleInstructor string = "Instructor"
	RoleTA         string = "TeachingAssistant"

	//TODO: Find out what the string value is
	RoleCourseBuilder string = ""
	RoleFacilitator   string = ""
	RoleSpectator     string = ""
	RoleGraderstring  string = ""
	RoleGuest         string = ""
)

// type UserCourseEnrollment struct {
// 	Role         string       `json:"courseRoleId"`
// 	Availability Availability `json:"availability"`
// }

type CourseUser struct {
	ID           string
	UserName     string
	FirstName    string
	LastName     string
	Available    string
	CourseRoleID string
}

type courseUsersResponse struct {
	Results []struct {
		ID           string `json:"id"`
		CourseRoleID string `json:"courseRoleId"`
		User         struct {
			UserName string `json:"userName"`

			Name struct {
				Given  string `json:"given"`
				Family string `json:"family"`
			} `json:"name"`
		} `json:"user"`
		Availability struct {
			Available string `json:"available"`
		} `json:"availability"`
	} `json:"results"`
	Paging struct {
		NextPage string `json:"nextPage"`
	} `json:"paging"`
}

// Create function that only needs the bare minimum. For simple creates.
func (cs *CourseService) Create(ctx context.Context, courseID string, title string, termID string) error {
	courseID = strings.TrimSpace(courseID)
	title = strings.TrimSpace(title)
	termID = strings.TrimSpace(termID)

	if courseID == "" || title == "" || termID == "" {
		return errors.New("missing parameters: courseID, title, termID")
	}

	data := Course{
		CourseID:     courseID,
		Name:         title,
		TermID:       termID,
		Organization: false,
		Availability: CourseAvailability{
			Available: AvailabilityNo,
			Duration: CourseDuration{
				Type: "Continuous",
			},
		},
		Enrollment: Enrollment{
			Type: "InstructorLed",
		},
	}

	//fmt.Printf("%v\n", data)

	url := endpoints.Courses.Create()
	resp, err := cs.client.Post(ctx, url, data)
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
		return ErrCourseExist
		//client.Logger.Error(fmt.Sprintf("User with ID %s already exists", username))
	case 400:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("The request did not specify valid data. %s", string(body))
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}

// A more advanced version of Create()
func (cs *CourseService) CreatePro(ctx context.Context, req *CourseCreateRequest) error {
	courseID := strings.TrimSpace(*req.ExternalID)
	title := strings.TrimSpace(*req.Name)
	termID := strings.TrimSpace(*req.TermID)
	description := strings.TrimSpace(*req.Description)

	//TODO: Test if TermID is blank.
	if courseID == "" || title == "" || termID == "" {
		return errors.New("missing parameters: courseID, title, termID")
	}

	data := Course{
		CourseID:     courseID,
		Name:         title,
		TermID:       termID,
		Organization: false,
		Description:  description,
		Availability: CourseAvailability{
			Available: AvailabilityNo,
			Duration: CourseDuration{
				Type: "Continuous",
			},
		},
		Enrollment: Enrollment{
			Type: "InstructorLed",
		},
	}

	//fmt.Printf("%v\n", data)

	url := endpoints.Courses.Create()
	resp, err := cs.client.Post(ctx, url, data)
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
		return ErrCourseExist
		//client.Logger.Error(fmt.Sprintf("User with ID %s already exists", username))
	case 400:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("The request did not specify valid data. %s", string(body))
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}

// TODO: Rename parameters
func (cs *CourseService) CopyCourseByCourseID(ctx context.Context, sourceID, destinationID string) (string, error) {
	sourceID = strings.TrimSpace(sourceID)
	destinationID = strings.TrimSpace(destinationID)

	if sourceID == "" || destinationID == "" {
		return "", errors.New("sourceID and destinationID are required")
	}

	// Build request
	data := map[string]any{
		"targetCourse": map[string]string{
			"courseId": destinationID,
		},
		"copy": map[string]any{
			"adaptiveReleaseRules": true,
			"announcements":        true,
			"assessments":          true,
			"blogs":                true,
			"calendar":             true,
			"contacts":             true,
			"contentAlignments":    true,
			"contentAreas":         true,
			"discussions":          "ForumsAndStarterPosts",
			"glossary":             true,
			"gradebook":            true,
			"groupSettings":        true,
			"journals":             true,
			"retentionRules":       true,
			"rubrics":              true,
			"settings": map[string]bool{
				"availability":       false,
				"bannerImage":        true,
				"duration":           true,
				"enrollmentOptions":  true,
				"guestAccess":        true,
				"languagePack":       true,
				"navigationSettings": true,
				"observerAccess":     true,
			},
			"tasks": true,
			"wikis": true,
		},
	}

	url := endpoints.Courses.Copy(sourceID)
	resp, err := cs.client.Post(ctx, url, data)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case 202:
		// "Course was successfully created from
		taskURI := resp.Header.Get("Location")
		if taskURI == "" {
			return "", fmt.Errorf("202 Accepted received but Location header was missing")
		}
		return taskURI, nil
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return "", fmt.Errorf("failed to copy course (status %d): %s", resp.StatusCode, string(body))
	}
}

func (cs *CourseService) DoesCourseExist(ctx context.Context, courseID string) (bool, error) {
	courseID = strings.TrimSpace(courseID)
	if courseID == "" {
		return false, errors.New("courseID is required")
	}
	url := endpoints.Courses.GetByCourseId(courseID)
	resp, err := cs.client.Get(ctx, url)
	if err != nil {
		// TODO: This could be better?
		return false, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		return true, nil
	//TODO: Add other codes
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return false, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}

func (cs *CourseService) GetCourseByCourseId(ctx context.Context, courseId string) (*Course, error) {
	courseId = strings.TrimSpace(courseId)

	url := endpoints.Courses.GetByCourseId(courseId)
	resp, err := cs.client.Get(ctx, url)
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

	var course Course
	if err := json.Unmarshal(body, &course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (cs *CourseService) GetCourseById(ctx context.Context, id string) (*Course, error) {
	id = strings.TrimSpace(id)

	url := endpoints.Courses.GetById(id)
	resp, err := cs.client.Get(ctx, url)
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

	var course Course
	if err := json.Unmarshal(body, &course); err != nil {
		return nil, err
	}

	return &course, nil
}

func (cs *CourseService) AddChildCourse(ctx context.Context, courseID string, childID string) error {
	url := endpoints.Courses.AddChildCourse(courseID, childID)

	resp, err := cs.client.Put(ctx, url, nil)
	if err != nil {
		return fmt.Errorf("request failed adding child course %s to %s: %w", childID, courseID, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	return fmt.Errorf("failed to add child course %s to %s (HTTP %d - %s)", childID, courseID, resp.StatusCode, string(body))
}

// Tested 11/4/25
func (cs *CourseService) DeleteCourse(ctx context.Context, courseID string) error {
	url := endpoints.Courses.GetByCourseId(courseID)
	resp, err := cs.client.Delete(ctx, url)
	if err != nil {
		return fmt.Errorf("request failed deleting course %s: %v", courseID, err)
	}
	defer resp.Body.Close()

	//TODO: I think I can remove status accepted, need to check docs
	if resp.StatusCode == http.StatusAccepted || resp.StatusCode == http.StatusOK {
		return nil
	}

	body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	return fmt.Errorf("failed to delete course %s (HTTP code %d) %s", courseID, resp.StatusCode, string(body))
}

// func (cs *CourseService) EnrollUserByUsername(ctx context.Context, courseId string, username string) error {
// 	return errors.New("EnrollUserByUsername not implemented")
// }

func (cs *CourseService) CreateMembership(ctx context.Context, username, courseId string, update EnrollmentRequest) error {
	return cs.upsertMembership(ctx, "PUT", username, courseId, update)
}

func (cs *CourseService) UpdateMembership(ctx context.Context, username, courseId string, update EnrollmentRequest) error {
	return cs.upsertMembership(ctx, "PATCH", username, courseId, update)
}

// TODO: Create and update are the same, other than put or patch. Could this be better?
func (cs *CourseService) upsertMembership(ctx context.Context, method string, username string, courseId string, update EnrollmentRequest) error {
	//TODO: Add course not found, user already enrolled, blah blah blah...
	username = strings.TrimSpace(username)
	courseId = strings.TrimSpace(courseId)
	if username == "" || courseId == "" {
		return ErrEmptyStringParameter
	}

	url := endpoints.Courses.CreateMembership(courseId, username)
	var resp *http.Response
	var err error

	switch method {
	case "PUT":
		resp, err = cs.client.Put(ctx, url, update)
	case "PATCH":
		resp, err = cs.client.Patch(ctx, url, update)
	default:
		return fmt.Errorf("unsupported method: %s", method)
	}

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK, http.StatusNoContent, http.StatusCreated:
		return nil

	case http.StatusNotFound:
		return ErrUserNotFound

	case 409:
		return ErrUserAlreadyEnrolled

	case http.StatusBadRequest, http.StatusInternalServerError:
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

func (cs *CourseService) UpdateMembershipAvailability(ctx context.Context, username string, courseId string, availability string) error {
	//TODO: Add course not found, user already enrolled, blah blah blah...
	username = strings.TrimSpace(username)
	courseId = strings.TrimSpace(courseId)
	availability = strings.TrimSpace(availability)
	if username == "" || courseId == "" || availability == "" {
		return ErrInvalidUsername
	}

	updateReq := EnrollmentRequest{
		Availability: &MembershipAvailability{
			Available: ToPtr(availability),
		},
	}

	return cs.UpdateMembership(ctx, username, courseId, updateReq)

}

func (cs *CourseService) Update(ctx context.Context, courseID string, req *CourseUpdateRequest) (*Course, error) {
	courseID = strings.TrimSpace(courseID)
	if courseID == "" {
		return nil, errors.New("courseID is required")
	}

	// Validate at least one field is being updated
	if req.Name == nil && req.TermID == nil && req.Availability == nil && req.DataSourceID == nil {
		return nil, errors.New("at least one field must be provided for update")
	}

	// Trim strings
	if req.Name != nil {
		trimmed := strings.TrimSpace(*req.Name)
		req.Name = ToPtr(trimmed)
	}
	if req.TermID != nil {
		trimmed := strings.TrimSpace(*req.TermID)
		req.TermID = &trimmed
	}
	if req.DataSourceID != nil {
		trimmed := strings.TrimSpace(*req.DataSourceID)
		req.DataSourceID = &trimmed
	}
	if req.Description != nil {
		trimmed := strings.TrimSpace(*req.Description)
		req.Description = &trimmed
	}

	url := endpoints.Courses.Update(courseID)
	resp, err := cs.client.Patch(ctx, url, req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		var updated Course
		if err := json.NewDecoder(resp.Body).Decode(&updated); err != nil {
			return nil, fmt.Errorf("failed to decode updated course: %w", err)
		}
		return &updated, nil

	case 404:
		return nil, fmt.Errorf("course %s not found", courseID)
	case 403:
		return nil, ErrInsufficientPrivileges
	case 400:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return nil, fmt.Errorf("bad request: %s", string(body))
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

}

// EnrollUserIntoCourse is a wrapper function that calls CreateMembership.
func (cs *CourseService) EnrollUserIntoCourse(ctx context.Context, courseId string, username string, role string, availability string) error {
	updateReq := EnrollmentRequest{
		CourseRoleID: ToPtr(role),
		Availability: &MembershipAvailability{
			Available: ToPtr(availability),
		},
	}
	return cs.CreateMembership(ctx, username, courseId, updateReq)
}

// RemoveUser will remove a user from a course.
func (cs *CourseService) RemoveUser(ctx context.Context, courseId string, username string) error {
	username = strings.TrimSpace(username)
	courseId = strings.TrimSpace(courseId)
	if username == "" || courseId == "" {
		return ErrEmptyStringParameter
	}

	url := endpoints.Courses.DeleteMembership(courseId, username)
	var resp *http.Response
	var err error

	resp, err = cs.client.Delete(ctx, url)

	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNoContent {
		return nil

	} else {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf(
			"deleting membership failed: status=%d body=%s",
			resp.StatusCode,
			string(body),
		)
	}
}

//TODO: Re-evaluate if the grade stuff needs to be its own thing

func (cs *CourseService) GetUsers(ctx context.Context, courseID string) ([]CourseUser, error) {
	//TODO: Move to endpoint file
	path := fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users?expand=user&fields=id,courseRoleId,user.userName,availability.available,user.name.given,user.name.family", courseID)

	var allUsers []CourseUser

	for {
		resp, err := cs.client.Get(ctx, path)
		if err != nil {
			return nil, fmt.Errorf("failed to get course users: %w", err)
		}
		defer resp.Body.Close()

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}

		var result courseUsersResponse
		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		for _, r := range result.Results {
			allUsers = append(allUsers, CourseUser{
				ID:           r.ID,
				UserName:     r.User.UserName,
				FirstName:    r.User.Name.Given,
				LastName:     r.User.Name.Family,
				Available:    r.Availability.Available,
				CourseRoleID: r.CourseRoleID,
			})
		}

		if result.Paging.NextPage == "" {
			break
		}
		path = result.Paging.NextPage
	}

	return allUsers, nil
}
