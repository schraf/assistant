package retry

import (
	"context"
	"log/slog"
	"math"
	"time"
)

type Retryer struct {
	MaxRetries       int
	InitialBackoff   time.Duration
	MaxBackoff       time.Duration
	IsRetryableError func(error) bool
	Attempt          func(context.Context) error
}

func (r Retryer) Try(ctx context.Context) error {
	attempt := 0

	for {
		err := r.Attempt(ctx)
		if err == nil {
			return nil
		}

		if !r.IsRetryableError(err) {
			return err
		}

		if attempt < r.MaxRetries {
			backoff := r.calculateBackoff(attempt)

			slog.Warn("retry_attempt_failed",
				slog.Int("attempt", attempt),
				slog.Int("max_retries", r.MaxRetries),
				slog.String("error", err.Error()),
				slog.Duration("backoff", backoff),
			)

			select {
			case <-ctx.Done():
				return ctx.Err()
			case <-time.After(backoff):
			}
		} else {
			return err
		}

		attempt++
	}
}

func (r Retryer) calculateBackoff(attempt int) time.Duration {
	backoff := float64(r.InitialBackoff) * math.Pow(2, float64(attempt))
	if backoff > float64(r.MaxBackoff) {
		backoff = float64(r.MaxBackoff)
	}

	return time.Duration(backoff)
}
