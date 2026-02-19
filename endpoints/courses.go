package endpoints

import "fmt"

type courseEndpoints struct{}

var Courses = courseEndpoints{}

func (courseEndpoints) Create() string {
	return "/learn/api/public/v3/courses"
}

func (courseEndpoints) Update(courseId string) string {
	return Courses.GetByCourseId(courseId)
}

func (courseEndpoints) Copy(courseId string) string {
	return fmt.Sprintf("/learn/api/public/v2/courses/courseId:%s/copy", courseId)
}

func (courseEndpoints) GetByCourseId(courseId string) string {
	return fmt.Sprintf("/learn/api/public/v3/courses/courseId:%s", courseId)
}

func (courseEndpoints) GetTask(courseId string, uri string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/tasks/%s", courseId, uri)
}

func (courseEndpoints) GetById(id string) string {
	return fmt.Sprintf("/learn/api/public/v3/courses/%s", id)
}

func GetUsers(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users", courseID)
}

func (courseEndpoints) AddChildCourse(courseId string, childId string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/children/courseId:%s", courseId, childId)
}

//TODO: Right now, only accepts courseID and username. Add a way to do this with system ids
func (courseEndpoints) GetMembership(courseId string, username string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users/userName:%s", courseId, username)
}

func (courseEndpoints) CreateMembership(courseId string, username string) string {
	return Courses.GetMembership(courseId, username)
	//return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users/userName:%s", courseId, username)
}

func (courseEndpoints) DeleteMembership(courseId string, username string) string {
	return Courses.GetMembership(courseId, username)
	//return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users/userName:%s", courseId, username)
}

//TODO: Add the rest
