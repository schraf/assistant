package scheduler

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
)

type JobScheduler interface {
	ScheduleJob(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error
}
