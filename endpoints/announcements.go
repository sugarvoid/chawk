package endpoints

import "fmt"

type announcementEndpoints struct{}

var Announcements = announcementEndpoints{}

func (announcementEndpoints) Create() string {
	return "/learn/api/public/v3/courses"
}

func (announcementEndpoints) GetAllByCourseId(courseId string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/announcements", courseId)
}

func (announcementEndpoints) GetSingleById(courseID, announcementID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/announcements/%s", courseID, announcementID)
}
