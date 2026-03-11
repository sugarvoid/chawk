package endpoints

import "fmt"

type courseEndpoints struct{}

var Courses = courseEndpoints{}

func (courseEndpoints) Create() string {
	return "/learn/api/public/v3/courses"
}

func (courseEndpoints) Update(courseID string) string {
	return Courses.GetByCourseId(courseID)
}

func (courseEndpoints) Copy(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v2/courses/courseId:%s/copy", courseID)
}

func (courseEndpoints) GetByCourseId(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v3/courses/courseId:%s", courseID)
}

func (courseEndpoints) GetTask(courseID string, uri string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/tasks/%s", courseID, uri)
}

func (courseEndpoints) GetById(id string) string {
	return fmt.Sprintf("/learn/api/public/v3/courses/%s", id)
}

func GetUsers(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users", courseID)
}

func (courseEndpoints) AddChildCourse(courseID string, childId string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/children/courseId:%s", courseID, childId)
}

//TODO: Right now, only accepts courseID and username. Add a way to do this with system ids
func (courseEndpoints) GetMembership(courseID string, username string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users/userName:%s", courseID, username)
}

func (courseEndpoints) GetContent(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/contents", courseID)
}

func (courseEndpoints) CreateMembership(courseID string, username string) string {
	return Courses.GetMembership(courseID, username)
	//return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users/userName:%s", courseId, username)
}

func (courseEndpoints) DeleteMembership(courseID string, username string) string {
	return Courses.GetMembership(courseID, username)
	//return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/users/userName:%s", courseId, username)
}

//TODO: Add the rest
