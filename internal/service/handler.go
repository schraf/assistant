package service

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"

	"github.com/google/uuid"
	"github.com/schraf/assistant/internal/log"
	"github.com/schraf/assistant/pkg/models"
)

func HandleRequest(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	requestId := uuid.New()
	logger := log.NewLogger()

	logger = logger.With(
		slog.String("request_id", requestId.String()),
	)

	logger.Info("received_request")

	//--==================================================================--
	//--== AUTHENTICATE THE REQUEST
	//--==================================================================--

	apiToken := os.Getenv("API_TOKEN")
	if apiToken == "" {
		logger.WarnContext(ctx, "no_api_token")

		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}

	requestToken := r.Header.Get("X-API-Token")

	if requestToken == "" {
		logger.WarnContext(ctx, "missing X-API-Token header")

		http.Error(w, "missing X-API-Token header", http.StatusBadRequest)
		return
	}

	if requestToken != apiToken {
		logger.WarnContext(ctx, "invalid_api_token",
			slog.String("request_token", requestToken),
		)

		http.Error(w, "invalid api token", http.StatusUnauthorized)
		return
	}

	logger.Info("authenticated_request")

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
	//--== BUILD CONFIG OBJECT
	//--==================================================================--

	configPrefix := "X-Config-"
	config := make(map[string]any)

	for name, values := range r.Header {
		if strings.HasPrefix(name, configPrefix) {
			key := name[len(configPrefix):]

			if len(values) == 1 {
				config[key] = values[0]
			} else if len(values) > 1 {
				config[key] = values
			}
		}
	}

	//--==================================================================--
	//--== START THE CLOUD RUN JOB
	//--==================================================================--

	req := models.ContentRequest{
		Id:   requestId,
		Body: body,
	}

	if err := startJob(ctx, contentType, config, req); err != nil {
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
