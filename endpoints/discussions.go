package endpoints

import "fmt"

type discussionEndpoints struct{}

var Discussions = discussionEndpoints{}

func (discussionEndpoints) GetAll(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/discussions/", courseID)
}

func (discussionEndpoints) GetMessages(courseID, forumID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/discussions/%s/messages/", courseID, forumID)
}

func (discussionEndpoints) DeleteMessage(courseID, forumID, messageID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/discussions/%s/messages/%s", courseID, forumID, messageID)
}
