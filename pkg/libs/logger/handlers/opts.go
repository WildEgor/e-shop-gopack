package slogger

import (
	"context"
	"fmt"
	"github.com/WildEgor/e-shop-gopack/pkg/libs/logger/models"
	"io"
	"log/slog"
	"os"
	"strings"
)

// Option is an application option.
type Option func(o *options)

// Options
type Options struct {
	DisableSource bool
	FullSource    bool
	DisableTime   bool
	DisableColor  bool // for cli

	// for format
	Name         string
	Organization string
	Context      string
}

// options is an application options.
type options struct {
	Options

	Level   string    // debug, info, warn, error
	Format  string    // json, text, pretty
	Output  string    // stdout, stderr, discard, or a file path
	Writer  io.Writer // set this to override Output
	Tracing bool      // enable tracing feature
}

// WithOrganization set name of organization
func WithOrganization(name string) Option {
	return func(o *options) {
		o.Organization = name
	}
}

// WithAppName application name
func WithAppName(name string) Option {
	return func(o *options) {
		o.Name = name
	}
}

// WithDisableSource disable source log
func WithDisableSource() Option {
	return func(o *options) { o.DisableSource = true }
}

// WithFullSource disable crop source
func WithFullSource() Option {
	return func(o *options) { o.FullSource = true }
}

// WithDisableTime disable time in log
func WithDisableTime() Option {
	return func(o *options) { o.DisableTime = true }
}

// WithLevel set log level. Default "info"
func WithLevel(level string) Option {
	return func(o *options) {
		if level == "" {
			level = "info"
		}
		o.Level = level
	}
}

// WithFormat set output format. Default "json"
func WithFormat(format string) Option {
	return func(o *options) {
		if format == "" {
			format = "json"
		}
		o.Format = format
	}
}

// WithOutput set output. Default "stdout"
func WithOutput(output string) Option {
	return func(o *options) {
		if output == "" {
			output = "stdout"
		}
		o.Output = output
	}
}

// WithWriter set writer
func WithWriter(w io.Writer) Option {
	return func(o *options) { o.Writer = w }
}

// WithTracing enable tracing
func WithTracing() Option {
	return func(o *options) { o.Tracing = true }
}

// ------------------------------------------------------------------------

// NewLogger create a new *slog.Logger with tracing handler wrapper
func NewLogger(opts ...Option) *slog.Logger {
	options := &options{
		Options: Options{
			DisableSource: false,
			FullSource:    false,
			DisableTime:   false,
		},
		Level:  "debug",
		Format: "json",
		Output: "stdout",
	}
	for _, o := range opts {
		o(options)
	}

	h := NewHandler(options)
	if options.Tracing {
		h = NewTracingHandler(h)
	}

	return slog.New(h)
}

// NewHandler create new slog handler. Default stdout with json handler and info level
func NewHandler(options *options) slog.Handler {
	var w io.Writer

	if options.Writer != nil {
		w = options.Writer
	} else {
		switch options.Output {
		case "stdout":
			w = os.Stdout
		case "stderr", "":
			w = os.Stderr
		case "discard":
			w = io.Discard
		default:
			var err error

			w, err = os.OpenFile(options.Output, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)

			if err != nil {
				slog.Error("failed to open log file, fallback to stderr", err)
				w = os.Stderr
			}
		}
	}

	var convLvl slog.Level
	switch options.Level {
	case "debug":
		convLvl = slog.LevelDebug
	case "info":
		convLvl = slog.LevelInfo
	case "warn":
		convLvl = slog.LevelWarn
	case "error":
		convLvl = slog.LevelError
	default:
		convLvl = slog.LevelInfo
	}

	lvl := &slog.LevelVar{}
	lvl.Set(convLvl)

	opts := NewHandlerOptions(lvl, &options.Options)
	var th slog.Handler
	switch options.Format {
	case "text":
		th = slog.NewTextHandler(w, &opts)
	case "pretty":
		th = NewCLIHandler(w, &CLIHandlerOptions{
			DisableColor:   options.DisableColor,
			HandlerOptions: opts,
		})
	case "json":
		fallthrough
	default:
		th = slog.NewJSONHandler(w, &opts)
	}

	return th
}

// NewHandlerOptions create new options with level
func NewHandlerOptions(level slog.Leveler, opt *Options) slog.HandlerOptions {
	ho := slog.HandlerOptions{
		AddSource: !opt.DisableSource,
		Level:     level,
	}

	if !opt.DisableTime && (opt.FullSource || opt.DisableSource) {
		return ho
	}

	ho.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
		if opt.DisableTime {
			if a.Key == slog.TimeKey {
				// Remove time from the output.
				return slog.Attr{}
			}
		}

		switch a.Key {
		case slog.SourceKey:
			if !opt.DisableSource && !opt.FullSource {
				return handleSourceKey(a)
			}
		case slog.MessageKey:
			return handleMsgKey(a)
		case models.SlogDataKey:
			return handleDataKey(a, opt)
		}
		return a
	}
	return ho
}

// handleDataKey process "data" key
func handleDataKey(a slog.Attr, opt *Options) slog.Attr {
	v := a.Value.Any()

	if log, ok := v.(*models.LogEntry); ok {
		l := make([]string, 0)

		if len(log.App) == 0 {
			log.App = opt.Name

			l = append(l, opt.Name)
		}

		if len(log.Organization) == 0 {
			log.Organization = opt.Organization

			l = append(l, opt.Organization)
		}

		if len(log.Label) == 0 {
			log.Label = strings.Join(l, ".")
		}

		ev, ok := log.Err.(error)
		if ok {
			log.Err = ev.Error()
		}

		v = log
	}

	return slog.Attr{
		Key:   a.Key,
		Value: slog.AnyValue(v),
	}
}

// handleMsgKey convert msg key
func handleMsgKey(a slog.Attr) slog.Attr {
	v := a.Value.String()

	return slog.Attr{
		Key:   models.SlogMessageKey,
		Value: slog.StringValue(v),
	}
}

// handleSourceKey make source key short
func handleSourceKey(a slog.Attr) slog.Attr {
	file := a.Value.String()
	if src, ok := a.Value.Any().(*slog.Source); ok {
		short := src.File

		idx := strings.LastIndexByte(src.File, '/')
		if idx > 0 {
			idx = strings.LastIndexByte(src.File[:idx], '/')
			if idx > 0 {
				short = src.File[idx+1:]
			}
		}

		file = fmt.Sprintf("%s:%d", short, src.Line)
	}

	return slog.Attr{
		Key:   a.Key,
		Value: slog.StringValue(file),
	}
}

// ------------------------------------------------------------------------

type ctxLogger struct{}

// ContextWithLogger adds logger to context.
func ContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, ctxLogger{}, l)
}

// ExtractLoggerFromContext returns logger from context.
func ExtractLoggerFromContext(ctx context.Context) *slog.Logger {
	if l, ok := ctx.Value(ctxLogger{}).(*slog.Logger); ok {
		return l
	}

	return NewLogger()
}
