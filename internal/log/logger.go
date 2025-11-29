package log

import (
	"context"
	"log/slog"
	"os"
)

// cloudLoggingHandler wraps a JSON handler and converts slog levels to
// Google Cloud Logging severity levels for proper log level detection.
type cloudLoggingHandler struct {
	slog.Handler
}

// Handle converts the log record to use "severity" instead of "level"
// and maps slog levels to Cloud Logging severity values.
func (h *cloudLoggingHandler) Handle(ctx context.Context, r slog.Record) error {
	// Create a new record with the severity field
	attrs := make([]slog.Attr, 0, r.NumAttrs()+1)

	// Map slog level to Cloud Logging severity
	severity := levelToSeverity(r.Level)
	attrs = append(attrs, slog.String("severity", severity))

	// Copy all other attributes, but skip the default "level" attribute
	r.Attrs(func(a slog.Attr) bool {
		if a.Key != "level" {
			attrs = append(attrs, a)
		}
		return true
	})

	// Create a new record with the modified attributes
	newRecord := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
	for _, attr := range attrs {
		newRecord.AddAttrs(attr)
	}

	return h.Handler.Handle(ctx, newRecord)
}

// levelToSeverity maps slog.Level to Google Cloud Logging severity strings.
func levelToSeverity(level slog.Level) string {
	switch {
	case level >= slog.LevelError:
		return "ERROR"
	case level >= slog.LevelWarn:
		return "WARNING"
	case level >= slog.LevelInfo:
		return "INFO"
	case level >= slog.LevelDebug:
		return "DEBUG"
	default:
		return "DEFAULT"
	}
}

// NewLogger creates and returns a slog logger configured to output JSON logs
// with Google Cloud Logging compatible severity levels.
// Logs are written to stdout only.
// The log level defaults to Info, but will be set to Debug if the "DEBUG"
// environment variable is set.
func NewLogger() *slog.Logger {
	// Determine log level based on DEBUG environment variable
	level := slog.LevelInfo
	if os.Getenv("DEBUG") != "" {
		level = slog.LevelDebug
	}

	// Create JSON handler writing to stdout
	jsonHandler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: level,
		// ReplaceAttr removes the default "level" field since we'll add "severity" instead
		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// Remove the default "level" attribute - we'll add "severity" in the handler
			if a.Key == "level" {
				return slog.Attr{}
			}
			return a
		},
	})

	// Wrap with Cloud Logging compatible handler
	cloudHandler := &cloudLoggingHandler{Handler: jsonHandler}

	// Return the logger
	return slog.New(cloudHandler)
}
