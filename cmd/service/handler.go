package main

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/schraf/assistant/internal/utils"
	"github.com/schraf/assistant/pkg/models"
)

func handleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := uuid.New()
	logger := utils.NewLogger()

	logger = logger.With(
		slog.String("request_id", requestId.String()),
	)

	logger.Info("received_research_request")

	//--==================================================================--
	//--== DECODE THE REQUEST BODY
	//--==================================================================--

	var body map[string]any

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		logger.WarnContext(ctx, "invalid_json_payload",
			slog.String("error", err.Error()),
		)

		http.Error(w, "invalid JSON payload", http.StatusBadRequest)
		return
	}

	//--==================================================================--
	//--== DETERMINE JOB TYPE
	//--==================================================================--

	jobName := r.Header.Get("X-Job-Name")

	if jobName == "" {
		logger.WarnContext(ctx, "missing X-Job-Name header")

		http.Error(w, "missing X-Job-Name header", http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "request_job",
		slog.String("job_name", jobName),
	)

	//--==================================================================--
	//--== START THE CLOUD RUN JOB
	//--==================================================================--

	req := models.ContentRequest{
		Id:   requestId,
		Body: body,
	}

	if err := startJob(ctx, jobName, req); err != nil {
		logger.WarnContext(ctx, "failed_executing_job",
			slog.String("job_name", jobName),
			slog.String("error", err.Error()),
		)

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	//--==================================================================--
	//--== SEND A RESPONSE BACK
	//--==================================================================--

	response := map[string]interface{}{
		"success":    true,
		"request_id": requestId,
		"job_name":   jobName,
		"message":    "assistant request queued",
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(response); err != nil {
		logger.WarnContext(ctx, "failed_encoding_response",
			slog.String("error", err.Error()),
		)
	}
}
