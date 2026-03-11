package chawk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"

	endpoints "github.com/sugarvoid/chawk/endpoints"
)

type AnnouncementService struct {
	client *BlackboardClient
}

type Duration struct {
	Type  *string `json:"type"`
	Start *string `json:"start"`
	End   *string `json:"end"`
}

type AnnouncementAvailability struct {
	Duration Duration `json:"duration"`
}

type Announcement struct {
	ID            string                   `json:"id"`
	Title         string                   `json:"title"`
	Body          string                   `json:"body"`
	Draft         bool                     `json:"draft"`
	Availability  AnnouncementAvailability `json:"availability"`
	CreatorUserID string                   `json:"creatorUserId"`
	Created       string                   `json:"created"`
	Modified      string                   `json:"modified"`
	Position      int                      `json:"position"`
	Creator       string                   `json:"creator"`
}

// GetAllAnnouncements returns all the announcements from a course
func (c *AnnouncementService) GetAllAnnouncements(ctx context.Context, courseID string) ([]Announcement, error) {
	courseID, err := RequiredString(courseID, "courseID")
	if err != nil {
		return nil, err
	}

	url := endpoints.Announcements.GetAllByCourseId(courseID)

	var allAnnouncements []Announcement

	for {
		resp, err := c.client.Get(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to get announcements: %w", err)
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
			Results []Announcement `json:"results"`
			Paging  struct {
				NextPage string `json:"nextPage"`
			} `json:"paging"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allAnnouncements = append(allAnnouncements, result.Results...)

		if result.Paging.NextPage == "" {
			break
		}
		url = result.Paging.NextPage
	}

	return allAnnouncements, nil
}

// GetAnnouncement get a single announcement by its ID
func (c *AnnouncementService) GetAnnouncement(ctx context.Context, courseID string, announcementID string) (Announcement, error) {
	// TODO: implement Blackboard announcement get
	courseID = strings.TrimSpace(courseID)
	exists, _ := c.client.Courses.DoesCourseExist(ctx, courseID)

	if !exists {
		return Announcement{}, ErrCourseNotFound
	}

	//url := endpoints.GetAllAnnouncements(c.BaseURL, courseID)

	return Announcement{}, errors.New("GetAnnouncement not implemented")
}

func (c *AnnouncementService) UpdateAnnouncement(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("UpdateAnnouncement not implemented")
}

func (c *AnnouncementService) DeleteAnnouncement(ctx context.Context, courseID, announcementID string) error {
	url := endpoints.Announcements.DeleteById(courseID, announcementID)

	resp, err := c.client.Delete(ctx, url)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusNoContent:
		return nil
	case http.StatusForbidden:
		return ErrInsufficientPrivileges
	case http.StatusNotFound:
		return fmt.Errorf("announcement %s not found in course %s", announcementID, courseID)
	default:
		body, _ := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}

func (c *AnnouncementService) DeleteAllAnnouncements(ctx context.Context, courseID string) error {
	announcements, err := c.GetAllAnnouncements(ctx, courseID)
	if err != nil {
		return fmt.Errorf("failed to get announcements: %w", err)
	}

	for _, a := range announcements {
		if err := c.DeleteAnnouncement(ctx, courseID, a.ID); err != nil {
			return fmt.Errorf("failed to delete announcement %s: %w", a.ID, err)
		}
		//fmt.Printf("Deleted announcement %q\n", a.Title)
	}

	return nil
}
