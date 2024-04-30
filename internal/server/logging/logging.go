package logging

import (
	"context"
	"log/slog"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

// InitLogging creates a logger, makes it the [slog] default, and [echo.Echo]'s logger.
// Uses [echo.Echo.Debug] to determine logging level and source code inclusion.
func InitLogging(e *echo.Echo) {
	// Create logger
	var logger *slog.Logger
	if e.Debug {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: true,
			Level:     slog.LevelDebug,
		}))
	} else {
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
			AddSource: false,
			Level:     slog.LevelInfo,
		}))
	}

	// Make it slog's default
	slog.SetDefault(logger)

	// Setup Echo to use logger
	e.Use(middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogProtocol: true,
		LogRemoteIP: true,
		LogHost:     true,
		LogMethod:   true,
		LogURI:      true,
		LogStatus:   true,
		LogError:    true,
		HandleError: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			if v.Error == nil {
				logger.LogAttrs(context.Background(), slog.LevelInfo, "REQUEST",
					slog.String("protocol", v.Protocol),
					slog.String("ip", v.RemoteIP),
					slog.String("host", v.Host),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
				)
			} else {
				logger.LogAttrs(context.Background(), slog.LevelError, "REQUEST_ERROR",
					slog.String("protocol", v.Protocol),
					slog.String("ip", v.RemoteIP),
					slog.String("host", v.Host),
					slog.String("method", v.Method),
					slog.String("uri", v.URI),
					slog.Int("status", v.Status),
					slog.String("err", v.Error.Error()),
				)
			}
			return nil
		},
	}))
}
