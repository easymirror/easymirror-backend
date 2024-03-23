package log

import (
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/rs/zerolog"
)

// NewMiddlewareLogger creates and returns a logging middleware for Echo
func NewMiddlewareLogger() echo.MiddlewareFunc {
	logger := zerolog.New(os.Stdout)
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogURI:       true,
		LogStatus:    true,
		LogError:     true,
		HandleError:  true, // forwards error to the global error handler, so it can decide appropriate status code
		LogLatency:   true,
		LogProtocol:  true,
		LogRemoteIP:  true,
		LogHost:      true,
		LogMethod:    true,
		LogUserAgent: true,

		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			// Log all values
			if v.Error == nil {
				logger.Info().
					Str("URI", v.URI).
					Int("status", v.Status).
					Str("id", v.RequestID).
					Str("user-agent", v.UserAgent).
					Str("method", v.Method).
					Str("latency_human", v.Latency.String()).
					Str("host", v.Host).
					Str("remote_ip", v.RemoteIP).
					Str("protocol", v.Protocol).
					Time("time_utc", time.Now().UTC()).
					Msg("request")
			} else {
				logger.Error().
					Err(v.Error).
					Str("URI", v.URI).
					Int("status", v.Status).
					Str("id", v.RequestID).
					Str("user-agent", v.UserAgent).
					Str("method", v.Method).
					Str("latency_human", v.Latency.String()).
					Str("host", v.Host).
					Str("remote_ip", v.RemoteIP).
					Str("protocol", v.Protocol).
					Time("time_utc", time.Now().UTC()).
					Msg("request error")
			}
			return nil
		},
	})
}
