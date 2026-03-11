package chawk

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	endpoints "github.com/sugarvoid/chawk/endpoints"
)

//TODO: remove testing print staments

type DiscussionService struct {
	client *BlackboardClient
}

type discussion struct {
	ID string `json:"id"`
}

type message struct {
	ID     string `json:"id"`
	Author string `json:"userId"`
}

// func (ds *DiscussionService) getAllDiscussions(ctx context.Context, courseID string) ([]discussion, error) {
// 	courseID, err := RequiredString(courseID, "courseID")
// 	if err != nil {
// 		return nil, err
// 	}

// 	url := endpoints.Discussions.GetAll(courseID)

// 	var allDiscussions []discussion

// 	for {
// 		resp, err := ds.client.Get(ctx, url)
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to get discussions: %w", err)
// 		}

// 		body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
// 		resp.Body.Close()
// 		if err != nil {
// 			return nil, fmt.Errorf("failed to read response body: %w", err)
// 		}

// 		if resp.StatusCode != http.StatusOK {
// 			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
// 		}

// 		var result struct {
// 			Results []discussion `json:"results"`
// 			Paging  struct {
// 				NextPage string `json:"nextPage"`
// 			} `json:"paging"`
// 		}

// 		if err := json.Unmarshal(body, &result); err != nil {
// 			return nil, fmt.Errorf("failed to parse response: %w", err)
// 		}

// 		allDiscussions = append(allDiscussions, result.Results...)

// 		if result.Paging.NextPage == "" {
// 			break
// 		}
// 		url = result.Paging.NextPage
// 	}

// 	return allDiscussions, nil
// }

func (d *DiscussionService) getDiscussions(ctx context.Context, courseID string) ([]discussion, error) {
	url := endpoints.Discussions.GetAll(courseID)

	var allDiscussions []discussion

	for {
		resp, err := d.client.Get(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to get discussions: %w", err)
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}

		var result struct {
			Results []discussion `json:"results"`
			Paging  struct {
				NextPage string `json:"nextPage"`
			} `json:"paging"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allDiscussions = append(allDiscussions, result.Results...)

		if result.Paging.NextPage == "" {
			break
		}
		url = result.Paging.NextPage
	}

	return allDiscussions, nil
}

func (d *DiscussionService) getMessages(ctx context.Context, courseID, forumID string) ([]message, error) {
	url := endpoints.Discussions.GetMessages(courseID, forumID)

	var allMessages []message

	for {
		resp, err := d.client.Get(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to get messages: %w", err)
		}

		body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		resp.Body.Close()

		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		if resp.StatusCode != http.StatusOK {
			return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
		}

		var result struct {
			Results []message `json:"results"`
			Paging  struct {
				NextPage string `json:"nextPage"`
			} `json:"paging"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allMessages = append(allMessages, result.Results...)

		if result.Paging.NextPage == "" {
			break
		}
		url = result.Paging.NextPage
	}

	return allMessages, nil
}

// ClearDiscussionStudentReplies deletes all posts from users with a given role.
// TODO: TEST THIS! MIGHT WORK!
func (d *DiscussionService) ClearStudentReplies(ctx context.Context, courseID, role string) error {
	courseID, err := RequiredString(courseID, "courseID")
	if err != nil {
		return err
	}

	if role == "" {
		role = "Student"
	}

	forums, err := d.getDiscussions(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get discussion IDs: %w", err)
	}

	for _, forum := range forums {
		messages, err := d.getMessages(ctx, courseID, forum.ID)
		if err != nil {
			return fmt.Errorf("failed to get messages for forum %s: %w", forum.ID, err)
		}

		for _, msg := range messages {
			username, err := d.client.Users.GetUserByUsername(ctx, msg.Author)
			if err != nil {
				// TODO: Alert me of fails, but keep going
				fmt.Errorf("failed to get username for user %s: %v", msg.Author, err)
				continue
			}

			courseMem, err := d.client.Courses.GetMembership(ctx, msg.Author, courseID)
			if err != nil {
				// TODO: Alert me of fails, but keep going
				fmt.Errorf("failed to get course role for user %s: %v", username, err)
				continue
			}

			if *courseMem.CourseRoleID == "student" {
				if err := d.deletePost(ctx, courseID, forum.ID, msg.ID); err != nil {
					fmt.Errorf("failed to delete post %s: %v", msg.ID, err)
				}
			}
		}
	}

	return nil
}

func (d *DiscussionService) deletePost(ctx context.Context, courseID, forumID, messageID string) error {
	url := endpoints.Discussions.DeleteMessage(courseID, forumID, messageID)

	resp, err := d.client.Delete(ctx, url)
	if err != nil {
		return fmt.Errorf("request failed: %w", err)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	resp.Body.Close()

	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}

	fmt.Printf("%s message %s deleted.", courseID, messageID)
	return nil
}
