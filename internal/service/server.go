package service

import (
	"fmt"
	"log/slog"
	"net/http"

	"github.com/schraf/assistant/internal/scheduler"
)

func StartServer(hostname string, port string, scheduler scheduler.JobScheduler) error {
	address := fmt.Sprintf("%s:%s", hostname, port)

	handler := NewHandler(scheduler)
	http.HandleFunc("/content", handler.HandleRequest)

	slog.Info("starting_server",
		slog.String("host", hostname),
		slog.String("port", port),
	)

	if err := http.ListenAndServe(address, nil); err != nil {
		slog.Error("failed_starting_service",
			slog.String("error", err.Error()),
		)

		return err
	}

	slog.Info("service_shutdown")
	return nil
}
