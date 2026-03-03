package endpoints

import "fmt"

type userEndpoints struct{}

var Users = userEndpoints{}

func (userEndpoints) Create() string {
	return "/learn/api/public/v1/users"
}

func (userEndpoints) Delete() string {
	return "/learn/api/public/v1/users"
}

func (userEndpoints) GetByUsername(username string) string {
	return fmt.Sprintf("/learn/api/public/v1/users/userName:%s", username)
}

func (userEndpoints) GetId(username string) string {
	return fmt.Sprintf("/learn/api/public/v1/users/userName:%s", username)
}

func (userEndpoints) GetMemberships(username string) string {
	//return fmt.Sprintf("/learn/api/public/v1/users/userName:%s/courses", username)
	return fmt.Sprintf("/learn/api/public/v1/users/userName:%s/courses?expand=course&fields=courseId,courseRoleId,created,course.externalId,course.name", username)
}
