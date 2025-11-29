package service

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"github.com/google/uuid"
	"github.com/schraf/assistant/internal/utils"
	"github.com/schraf/assistant/pkg/models"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
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
	//--== DETERMINE CONTENT TYPE
	//--==================================================================--

	contentType := r.Header.Get("X-Content-Type")

	if contentType == "" {
		logger.WarnContext(ctx, "missing X-Content-Type header")

		http.Error(w, "missing X-Content-Type header", http.StatusBadRequest)
		return
	}

	logger.InfoContext(ctx, "request_content",
		slog.String("content_type", contentType),
	)

	//--==================================================================--
	//--== START THE CLOUD RUN JOB
	//--==================================================================--

	req := models.ContentRequest{
		Id:   requestId,
		Body: body,
	}

	if err := startJob(ctx, contentType, req); err != nil {
		logger.WarnContext(ctx, "failed_executing_job",
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
