package chawk

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"

	endpoints "github.com/sugarvoid/chawk/endpoints"
)

type GradebookService struct {
	client *BlackboardClient
}

type ColumnScore struct {
	Possible float64 `json:"possible"`
}

type ColumnAvailability struct {
	Available string `json:"available"`
}

type GradebookGrading struct {
	Type            string `json:"type,omitempty"`
	Due             string `json:"due,omitempty"`
	AttemptsAllowed int    `json:"attemptsAllowed,omitempty"`
	ScoringModel    string `json:"scoringModel,omitempty"`
}

type GradebookColumn struct {
	ID                       string             `json:"id,omitempty"`
	ExternalID               string             `json:"externalId,omitempty"`
	Name                     string             `json:"name,omitempty"`
	DisplayName              string             `json:"displayName,omitempty"`
	Description              string             `json:"description,omitempty"`
	Created                  string             `json:"created,omitempty"`
	Modified                 string             `json:"modified,omitempty"`
	Score                    ColumnScore        `json:"score,omitzero"`
	Availability             ColumnAvailability `json:"availability,omitzero"`
	Grading                  GradebookGrading   `json:"grading,omitzero"`
	IncludeInCalculations    bool               `json:"includeInCalculations,omitempty"`
	ShowStatisticsToStudents bool               `json:"showStatisticsToStudents,omitempty"`
}

func (gs *GradebookService) GetColumn(ctx context.Context, courseID string, columnID string) error {
	// TODO: implement
	return errors.New("GetColumn not implemented")
}

func (gs *GradebookService) GetColumns(ctx context.Context, courseID string) ([]GradebookColumn, error) {
	courseID, err := RequiredString(courseID, "courseID")
	if err != nil {
		return nil, err
	}

	url := endpoints.Gradebook.GetColumns(courseID)

	var allColumns []GradebookColumn

	for {
		resp, err := gs.client.Get(ctx, url)
		if err != nil {
			return nil, fmt.Errorf("failed to get gradebook columns: %w", err)
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
			Results []GradebookColumn `json:"results"`
			Paging  struct {
				NextPage string `json:"nextPage"`
			} `json:"paging"`
		}

		if err := json.Unmarshal(body, &result); err != nil {
			return nil, fmt.Errorf("failed to parse response: %w", err)
		}

		allColumns = append(allColumns, result.Results...)

		if result.Paging.NextPage == "" {
			break
		}
		url = result.Paging.NextPage
	}

	return allColumns, nil
}

func (gs *GradebookService) CreateColumnPro(ctx context.Context, courseID string, column GradebookColumn) error {
	courseID, err := RequiredString(courseID, "courseID")
	if err != nil {
		return err
	}

	url := endpoints.Gradebook.CreateColumn(courseID)

	resp, err := gs.client.Post(ctx, url, column)
	if err != nil {
		return fmt.Errorf("failed to create column: %w", err)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf("invalid request data: %s", string(body))
	case http.StatusForbidden:
		return ErrInsufficientPrivileges
	case http.StatusConflict:
		return fmt.Errorf("column %q already exists in course %s", column.Name, courseID)
	default:
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}

func (gs *GradebookService) CreateColumn(ctx context.Context, courseID, name, description string, score float64) error {
	courseID, err := RequiredString(courseID, "courseID")
	if err != nil {
		return err
	}

	data := GradebookColumn{
		Name:         name,
		Description:  description,
		Score:        ColumnScore{Possible: score},
		Availability: ColumnAvailability{Available: "Yes"},
	}

	url := endpoints.Gradebook.CreateColumn(courseID)

	resp, err := gs.client.Post(ctx, url, data)
	if err != nil {
		return fmt.Errorf("failed to create column: %w", err)
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, MAX_RESPONSE_SIZE))
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("failed to read response body: %w", err)
	}

	switch resp.StatusCode {
	case http.StatusCreated:
		return nil
	case http.StatusBadRequest:
		return fmt.Errorf("invalid request data: %s", string(body))
	case http.StatusForbidden:
		return ErrInsufficientPrivileges
	case http.StatusConflict:
		return fmt.Errorf("column %q already exists in course %s", name, courseID)
	default:
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
}

func (gs *GradebookService) DeleteColumn(ctx context.Context, courseID string, columnID string) error {
	// TODO: implement
	return errors.New("DeleteColumn not implemented")
}

func (gs *GradebookService) DeleteColumns(ctx context.Context, courseID string) error {
	// TODO: implement
	return errors.New("DeleteColumns not implemented")
}

// func (gs *GradebookService) GetColumnValue(ctx context.Context, courseID string, columnID string) error {
// 	// TODO: implement
// 	return errors.New("GetColumnValue not implemented")
// }

func (gs *GradebookService) UpdateColumnValue(ctx context.Context, courseID string, columnID string) error {
	// TODO: implement
	return errors.New("UpdateColumnValue not implemented")
}

func (gs *GradebookService) UpdateColumnPro(ctx context.Context, courseID string, columnID string) error {
	// TODO: implement
	return errors.New("UpdateColumnPro not implemented")
}
