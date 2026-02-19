package chawk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"

	endpoints "chawk/endpoints"
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

func (c *AnnouncementService) GetAllAnnouncements(ctx context.Context, courseID string) ([]Announcement, error) {

	courseID, err := RequiredString(courseID, "courseID")

	if err != nil {
		return nil, err
	}

	url := endpoints.Announcements.GetAllByCourseId(courseID)

	resp, err := c.client.Get(ctx, url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed to get user: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	if err != nil {
		return nil, err
	}

	var response struct {
		Results []Announcement `json:"results"`
	}

	if err := json.Unmarshal(body, &response); err != nil {
		return nil, err
	}

	return response.Results, nil
}

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

func (c *AnnouncementService) deleteAnnouncement(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("deleteAnnouncement not implemented")
}

func (c *AnnouncementService) DeleteAnnouncements(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("DeleteAnnouncements not implemented")
}
