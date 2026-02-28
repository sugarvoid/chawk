package chawk

import (
	"context"
	"errors"
)

type GradebookService struct {
	client *BlackboardClient
}

type GradeColumn struct {
}

func (gs *GradebookService) CreateColumn(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("CreateGradeColumn not implemented")
}

func (gs *GradebookService) DeleteColumn(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("DeleteColumn not implemented")
}

func (gs *GradebookService) GetColumnValue(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("GetColumnValue not implemented")
}

func (gs *GradebookService) UpdateColumnValue(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("UpdateColumnValue not implemented")
}

func (gs *GradebookService) UpdateColumnPro(ctx context.Context, courseID string, announcementID string) error {
	// TODO: implement
	return errors.New("UpdateColumnPro not implemented")
}
