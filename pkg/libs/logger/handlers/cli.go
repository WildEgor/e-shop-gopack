package slogger

import (
	"context"
	"fmt"
	"io"
	"log/slog"
	"sync"

	slogUtils "github.com/WildEgor/g-core/pkg/libs/logger/utils"
	"github.com/fatih/color"
)

var Colors = map[slog.Level]*color.Color{
	slog.LevelDebug: color.New(color.FgWhite),
	slog.LevelInfo:  color.New(color.FgBlue),
	slog.LevelWarn:  color.New(color.FgYellow),
	slog.LevelError: color.New(color.FgRed),
}

var Strings = map[slog.Level]string{
	slog.LevelDebug: "•",
	slog.LevelInfo:  "•",
	slog.LevelWarn:  "•",
	slog.LevelError: "⨯",
}

var bold = color.New(color.Bold)

type CLIHandlerOptions struct {
	DisableColor bool
	slog.HandlerOptions
}

type CLIHandler struct {
	opts *CLIHandlerOptions
	slog.Handler
	attrsPrefix []slog.Attr
	groupPrefix string
	mu          sync.Mutex
	w           io.Writer
}

func NewCLIHandler(
	out io.Writer,
	opts *CLIHandlerOptions,
) *CLIHandler {
	h := &CLIHandler{
		opts: opts,
		w:    out,
	}

	return h
}

func (h *CLIHandler) Enabled(_ context.Context, l slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return l >= minLevel
}

func (h *CLIHandler) Handle(_ context.Context, r slog.Record) error {
	buf := slogUtils.NewBuffer()
	defer buf.Free()

	theColor := Colors[r.Level]

	if h.opts.DisableColor {
		theColor.DisableColor()
	} else {
		theColor.EnableColor()
	}

	levelEmoji := Strings[r.Level]
	padding := 4
	coloredLevel := theColor.Sprintf("%s", bold.Sprintf("%*s", padding, levelEmoji))
	_, err := buf.WriteString(coloredLevel)
	if err != nil {
		return err
	}

	_, err = buf.WriteString(" ")
	if err != nil {
		return err
	}
	_, err = buf.WriteString(fmt.Sprintf("%-25s", r.Message))
	if err != nil {
		return err
	}

	_, err = buf.WriteString("\t\t")
	if err != nil {
		return err
	}

	// write handler attributes
	if len(h.attrsPrefix) > 0 {
		for _, attr := range h.attrsPrefix {
			h.appendAttr(buf, attr, theColor, h.groupPrefix)
		}
	}

	// write attributes
	if r.NumAttrs() > 0 {
		r.Attrs(func(attr slog.Attr) bool {
			h.appendAttr(buf, attr, theColor, h.groupPrefix)
			return true
		})
	}

	err = buf.WriteByte('\n')
	if err != nil {
		return err
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	_, err = h.w.Write(buf.Bytes())
	if err != nil {
		return err
	}

	return nil
}

func (h *CLIHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}

	cloned := h.clone()
	cloned.attrsPrefix = append(cloned.attrsPrefix, attrs...)
	return cloned
}

func (h *CLIHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	cloned := h.clone()
	cloned.groupPrefix += name + "."
	return cloned
}

func (h *CLIHandler) appendAttr(buf *slogUtils.Buffer, attr slog.Attr, theColor *color.Color, groupsPrefix string) {
	_, err := buf.Write([]byte(" "))
	if err != nil {
		return
	}
	if groupsPrefix != "" {
		_, err := buf.WriteString(theColor.Sprint(groupsPrefix))
		if err != nil {
			return
		}
	}
	_, err = buf.WriteString(theColor.Sprint(attr.Key))
	if err != nil {
		return
	}
	_, err = buf.Write([]byte("="))
	if err != nil {
		return
	}

	if attr.Value.Kind() != slog.KindGroup {
		_, err = buf.Write([]byte(attr.Value.String()))
		if err != nil {
			return
		}
	} else {
		_, err = buf.Write([]byte("{"))
		if err != nil {
			return
		}
		for _, attr := range attr.Value.Group() {
			h.appendAttr(buf, attr, theColor, groupsPrefix)
		}
		_, err = buf.Write([]byte(" }"))
		if err != nil {
			return
		}
	}
}

func (h *CLIHandler) clone() *CLIHandler {
	attrsPrefix := make([]slog.Attr, len(h.attrsPrefix))
	copy(attrsPrefix, h.attrsPrefix)

	return &CLIHandler{
		w:           h.w,
		opts:        h.opts,
		attrsPrefix: attrsPrefix,
		groupPrefix: h.groupPrefix,
	}
}

var _ slog.Handler = (*CLIHandler)(nil)
