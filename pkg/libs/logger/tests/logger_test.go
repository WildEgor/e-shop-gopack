package logger_tests

import (
	"errors"
	slogger "github.com/WildEgor/g-core/pkg/libs/logger/handlers"
	"github.com/WildEgor/g-core/pkg/libs/logger/models"
	"log/slog"
	"testing"
)

func TestShowLogs(t *testing.T) {
	logger := slogger.NewLogger(
		slogger.WithOrganization("test"),
		slogger.WithAppName("app"),
		slogger.WithLevel("debug"),
		slogger.WithFormat("json"),
	)

	slog.SetDefault(logger)

	slog.Debug("debug or error", models.LogEntryAttr(&models.LogEntry{
		Err: errors.New("err"),
		Props: map[string]interface{}{
			"id": 1,
		},
	}))

	slog.Info("info")
}
