package endpoints

import "fmt"

type gradeEndpoints struct{}

var Grades = gradeEndpoints{}

func (gradeEndpoints) GetColumns(courseID string) string {
	return fmt.Sprintf("/learn/api/public/v2/courses/courseId:%s/gradebook/columns", courseID)
}

func (gradeEndpoints) GetColumn(courseID, columnID string) string {
	return fmt.Sprintf("/learn/api/public/v2/courses/courseId:%s/gradebook/columns/%s", courseID, columnID)
}
