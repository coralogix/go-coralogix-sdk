package coralogix

import (
	"context"
	"log/slog"
	"runtime"
)

type CoralogixHandler struct {
	// Next represents the next handler in the chain.
	Next slog.Handler
	// cxLogger is the Coralogix logger.
	cxLogger  *CoralogixLogger
	AddSource bool
}

type source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type logMessage struct {
	Message string         `json:"message"`
	Level   string         `json:"level"`
	Data    map[string]any `json:"data,omitempty"`
	Source  source         `json:"source,omitempty"`
}

func NewCoralogixHandler(privateKey, applicationName, subsystemName string, next slog.Handler) *CoralogixHandler {
	logger := NewCoralogixLogger(
		privateKey,
		applicationName,
		subsystemName,
	)

	return &CoralogixHandler{
		Next:     next,
		cxLogger: logger,
	}
}

// Handle handles the provided log record.
func (h *CoralogixHandler) Handle(ctx context.Context, r slog.Record) error {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	log := logMessage{
		Message: r.Message,
		Level:   r.Level.String(),
		Data:    map[string]interface{}{},
	}

	if h.AddSource {
		log.Source = source{
			Function: f.Function,
			File:     f.File,
			Line:     f.Line,
		}
	}

	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			attrToMap(log.Data, a)
			return true
		})
	}

	h.cxLogger.Log(levelSlogToCoralogix(r.Level), log, "", "", f.Function, "")

	return h.Next.Handle(ctx, r)
}

// WithAttrs returns a new Coralogix whose attributes consists of handler's attributes followed by attrs.
func (h *CoralogixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &CoralogixHandler{
		Next:      h.Next.WithAttrs(attrs),
		cxLogger:  h.cxLogger,
		AddSource: h.AddSource,
	}
}

// WithGroup returns a new Coralogix with a group, provided the group's name.
func (h *CoralogixHandler) WithGroup(name string) slog.Handler {
	return &CoralogixHandler{
		Next:      h.Next.WithGroup(name),
		cxLogger:  h.cxLogger,
		AddSource: h.AddSource,
	}
}

// Enabled reports whether the logger emits log records at the given context and level.
// Note: We handover the decision down to the next handler.
func (h *CoralogixHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.Next.Enabled(ctx, level)
}

func (h *CoralogixHandler) Stop() {
	h.cxLogger.Destroy()
}

func attrToMap(m map[string]any, a slog.Attr) {
	switch v := a.Value.Any().(type) {
	case error:
		m[a.Key] = v.Error()
	case []slog.Attr:
		m2 := map[string]any{}
		for _, a2 := range v {
			attrToMap(m2, a2)
			m[a.Key] = m2
		}
	default:
		m[a.Key] = a.Value.Any()
	}
}

func levelSlogToCoralogix(level slog.Level) uint {
	switch level {
	case slog.LevelDebug:
		return Level.DEBUG
	case slog.LevelInfo:
		return Level.INFO
	case slog.LevelWarn:
		return Level.WARNING
	case slog.LevelError:
		return Level.ERROR
	default:
		return uint(level)
	}
}
