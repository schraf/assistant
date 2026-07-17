package local

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log/slog"

	"github.com/schraf/assistant/pkg/models"
)

type LocalJobScheduler struct{}

func NewLocalJobScheduler() *LocalJobScheduler {
	return &LocalJobScheduler{}
}

func (s *LocalJobScheduler) ScheduleJob(ctx context.Context, contentType string, config map[string]any, request models.ContentRequest) error {
	//--==================================================================--
	//--== ENCODE THE REQUEST BODY
	//--==================================================================--

	requestBodyJson, err := json.Marshal(request.Body)
	if err != nil {
		return fmt.Errorf("failed to marshal content request: %w", err)
	}

	encodedRequestBody := base64.StdEncoding.EncodeToString(requestBodyJson)

	//--==================================================================--
	//--== ENCODE THE CONFIG
	//--==================================================================--

	configJson, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	encodedConfig := base64.StdEncoding.EncodeToString(configJson)

	//--==================================================================--
	//--== RUN LOCAL JOB
	//--==================================================================--

	// TODO: run async job

	slog.Info("scheduled_job",
		slog.String("request", encodedRequestBody),
		slog.String("config", encodedConfig),
	)

	return nil
}
