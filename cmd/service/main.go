package main

import (
	"log/slog"
	"os"
	"strings"

	"github.com/schraf/assistant/internal/config"
	"github.com/schraf/assistant/internal/log"
	"github.com/schraf/assistant/internal/scheduler"
	"github.com/schraf/assistant/internal/scheduler/gcp"
	"github.com/schraf/assistant/internal/scheduler/local"
	"github.com/schraf/assistant/internal/service"
)

func main() {
	logger := log.NewLogger()
	slog.SetDefault(logger)

	if err := config.LoadEnv(".env"); err != nil {
		logger.Error("load_env_failed",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	var scheduler scheduler.JobScheduler

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	hostname := "0.0.0.0"
	if localOnly := strings.ToLower(os.Getenv("LOCAL_ONLY")); localOnly == "true" || localOnly == "1" {
		hostname = "127.0.0.1"
		scheduler = local.NewLocalJobScheduler()
	} else {
		scheduler = gcp.NewCloudRunJobScheduler()
	}

	if err := service.StartServer(hostname, port, scheduler); err != nil {
		os.Exit(1)
	}

	os.Exit(0)
}
