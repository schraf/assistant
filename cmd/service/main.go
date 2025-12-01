package main

import (
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/schraf/assistant/internal/config"
	"github.com/schraf/assistant/internal/log"
	"github.com/schraf/assistant/internal/service"
)

func main() {
	logger := log.NewLogger()

	if err := config.LoadEnv(".env"); err != nil {
		logger.Error("load_env_failed",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	port := "8080"
	if envPort := os.Getenv("PORT"); envPort != "" {
		port = envPort
	}

	hostname := "0.0.0.0"
	if localOnly := strings.ToLower(os.Getenv("LOCAL_ONLY")); localOnly == "true" || localOnly == "1" {
		hostname = "127.0.0.1"
	}

	address := fmt.Sprintf("%s:%s", hostname, port)

	// Create job scheduler
	scheduler := service.NewCloudRunJobScheduler()

	// Create handler with scheduler
	handler := service.NewHandler(scheduler)
	http.HandleFunc("/content", handler.HandleRequest)

	logger.Info("starting_service",
		slog.String("host", hostname),
		slog.String("port", port),
	)

	if err := http.ListenAndServe(address, nil); err != nil {
		logger.Error("failed_starting_service",
			slog.String("error", err.Error()),
		)

		os.Exit(1)
	}

	logger.Info("service_shutdown")
	os.Exit(0)
}
