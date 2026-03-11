package chawk

import "time"

// TODO: These are untested and just were in the python version, so I brought ti over

// formatDate converts a Blackboard ISO 8601 timestamp to MM-DD-YYYY.
// Example: "2024-06-27T14:15:14.634Z" -> "06-27-2024"
func formatDate(dateString string) string {
	t, err := time.Parse("2006-01-02T15:04:05.000Z", dateString)
	if err != nil {
		return ""
	}
	return t.Format("01-02-2006")
}

// parseDate converts a MM-DD-YYYY string back to a Blackboard ISO 8601 timestamp.
// Example: "06-27-2024" -> "2024-06-27T00:00:00.000Z"
func parseDate(dateString string) string {
	t, err := time.Parse("01-02-2006", dateString)
	if err != nil {
		return ""
	}
	return t.UTC().Format("2006-01-02T15:04:05.000Z")
}
