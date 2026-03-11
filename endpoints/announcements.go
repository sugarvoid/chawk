package endpoints

import "fmt"

type announcementEndpoints struct{}

var Announcements = announcementEndpoints{}

func (announcementEndpoints) Create() string {
	return "/learn/api/public/v3/courses"
}

func (announcementEndpoints) GetAllByCourseId(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/announcements", courseID)
}

func (announcementEndpoints) GetSingleById(courseID, announcementID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/announcements/%s", courseID, announcementID)
}

func (announcementEndpoints) DeleteById(courseID, announcementID string) string {
	return Announcements.GetSingleById(courseID, announcementID)
	//return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/announcements/%s", courseID, announcementID)
}
