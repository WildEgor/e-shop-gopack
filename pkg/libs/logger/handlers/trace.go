package slogger

import (
	"context"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"log/slog"
)

const TraceIDKey = "trace_id"

// TracingHandler slog with traces
type TracingHandler struct {
	handler slog.Handler
}

// NewTracingHandler create tracing handler
func NewTracingHandler(h slog.Handler) *TracingHandler {
	if lh, ok := h.(*TracingHandler); ok {
		h = lh.Handler()
	}
	return &TracingHandler{h}
}

// Enabled implements Handler.Enabled by reporting whether level is at least as large as h's level.
func (h *TracingHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements Handler.Handle.
func (h *TracingHandler) Handle(ctx context.Context, r slog.Record) error {
	span := trace.SpanFromContext(ctx)

	if span.IsRecording() {
		if r.Level >= slog.LevelError {
			span.SetStatus(codes.Error, r.Message)
		}

		if spanCtx := span.SpanContext(); spanCtx.HasTraceID() {
			r.AddAttrs(slog.String(TraceIDKey, spanCtx.TraceID().String()))
		}
	}

	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *TracingHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewTracingHandler(h.handler.WithAttrs(attrs))
}

// WithGroup implements Handler.WithGroup.
func (h *TracingHandler) WithGroup(name string) slog.Handler {
	return NewTracingHandler(h.handler.WithGroup(name))
}

// Handler returns the Handler wrapped by handler.
func (h *TracingHandler) Handler() slog.Handler {
	return h.handler
}

var _ slog.Handler = (*TracingHandler)(nil)
