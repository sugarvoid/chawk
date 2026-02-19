package endpoints

import "fmt"

type discussionEndpoints struct{}

var Discussions = discussionEndpoints{}

func (discussionEndpoints) Discussions(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/discussions/", courseID)
}

func (discussionEndpoints) DiscussionMessages(courseID, forumID string) string {
	return fmt.Sprintf("/learn/api/public/v1/courses/courseId:%s/discussions/%s/messages/", courseID, forumID)
}
