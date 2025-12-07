package mocks

import (
	"context"

	"github.com/schraf/assistant/pkg/models"
)

// MockJobScheduler is a mock implementation of models.JobScheduler.
type MockJobScheduler struct {
	ScheduleJobFunc func(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error
}

// ScheduleJob calls ScheduleJobFunc if set, otherwise returns nil.
func (m *MockJobScheduler) ScheduleJob(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error {
	if m.ScheduleJobFunc != nil {
		return m.ScheduleJobFunc(ctx, contentType, config, request)
	}

	return nil
}
