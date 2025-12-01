package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/schraf/assistant/internal/config"
	"github.com/schraf/assistant/internal/gemini"
	"github.com/schraf/assistant/internal/job"
	"github.com/schraf/assistant/internal/log"
	"github.com/schraf/assistant/internal/notify"
	"github.com/schraf/assistant/internal/telegraph"
	_ "github.com/schraf/newspaper-assistant/pkg/generator"
	_ "github.com/schraf/research-assistant/pkg/generator"
)

func main() {
	ctx := context.Background()
	logger := log.NewLogger()
	slog.SetDefault(logger)

	if err := config.LoadEnv(".env"); err != nil {
		logger.ErrorContext(ctx, "load_env_failed",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Create assistant
	assistant, err := gemini.NewClient(ctx)
	if err != nil {
		logger.ErrorContext(ctx, "failed_creating_assistant",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	// Create dependencies
	publisher := telegraph.NewPublisher()
	notifier := notify.NewEmailNotifier()

	// Create processor
	processor := job.NewProcessor(assistant, publisher, notifier, logger)

	// Process the job
	if err := processor.Process(ctx); err != nil {
		logger.ErrorContext(ctx, "job_failed",
			slog.String("error", err.Error()),
		)
		os.Exit(1)
	}

	os.Exit(0)
}
