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

func GetMemberships(username string) string {
	return fmt.Sprintf("/learn/api/public/v1/users/userName:%s/courses", username)
}

//TODO: Add the rest
