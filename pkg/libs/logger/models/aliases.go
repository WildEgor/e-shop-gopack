package models

import (
	"log/slog"
	"time"
)

var (
	// SlogDataKey grouped data
	SlogDataKey = "data"
	// SlogMessageKey text field
	SlogMessageKey = "message"
	// SlogErrorKey error field
	SlogErrorKey = "error"
)

// Make type aliases to slog
type (
	Logger         = slog.Logger
	Attr           = slog.Attr
	Level          = slog.Level
	Handler        = slog.Handler
	Value          = slog.Value
	HandlerOptions = slog.HandlerOptions
	LogValuer      = slog.LogValuer
)

func Float32Attr(key string, val float32) Attr {
	return slog.Float64(key, float64(val))
}

func UInt32Attr(key string, val uint32) Attr {
	return slog.Int(key, int(val))
}

func Int32Attr(key string, val int32) Attr {
	return slog.Int(key, int(val))
}

func TimeAttr(key string, time time.Time) Attr {
	return slog.String(key, time.String())
}

func ErrAttr(err error) Attr {
	return slog.String(SlogErrorKey, err.Error())
}

func LogEntryAttr(log *LogEntry) Attr {
	return slog.Any(SlogDataKey, log)
}
