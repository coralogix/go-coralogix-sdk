package coralogix

import (
	"context"
	"log/slog"
	"runtime"
)

type CoralogixHandler struct {
	cxLogger *CoralogixLogger

	opts   slog.HandlerOptions
	data   map[string]interface{}
	groups []string
}

type source struct {
	Function string `json:"function"`
	File     string `json:"file"`
	Line     int    `json:"line"`
}

type logMessage struct {
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
	Source  *source        `json:"source,omitempty"`
}

func NewCoralogixHandler(privateKey, applicationName, subsystemName string, opts *slog.HandlerOptions) *CoralogixHandler {
	logger := NewCoralogixLogger(
		privateKey,
		applicationName,
		subsystemName,
	)

	return &CoralogixHandler{
		cxLogger: logger,
		opts:     *opts,
	}
}

func (h *CoralogixHandler) cloneData() map[string]interface{} {
	clone := map[string]interface{}{}
	for k, v := range h.data {
		clone[k] = v
	}

	return clone
}

func (h *CoralogixHandler) cloneGroups() []string {
	clone := make([]string, len(h.groups))
	for i, group := range h.groups {
		clone[i] = group
	}

	return clone
}

// Handle handles the provided log record.
func (h *CoralogixHandler) Handle(ctx context.Context, r slog.Record) error {
	fs := runtime.CallersFrames([]uintptr{r.PC})
	f, _ := fs.Next()

	log := logMessage{
		Message: r.Message,
		Data:    h.cloneData(),
	}

	grouped := groupToMap(log.Data, h.groups)
	if r.NumAttrs() > 0 {
		r.Attrs(func(a slog.Attr) bool {
			attrToMap(grouped, a)
			return true
		})
	}

	if h.opts.AddSource {
		log.Source = &source{
			Function: f.Function,
			File:     f.File,
			Line:     f.Line,
		}
	}

	category := ""
	if v, ok := log.Data["Category"]; ok {
		category = v.(string)
		delete(log.Data, "Category")
	}

	className := ""
	if v, ok := log.Data["ClassName"]; ok {
		className = v.(string)
		delete(log.Data, "ClassName")
	}

	threadId := ""
	if v, ok := log.Data["ThreadId"]; ok {
		threadId = v.(string)
		delete(log.Data, "ThreadId")
	}

	h.cxLogger.Log(levelSlogToCoralogix(r.Level), log, category, className, f.Function, threadId)
	return nil
}

// WithAttrs returns a new Coralogix whose attributes consists of handler's attributes followed by attrs.
func (h *CoralogixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	data := h.cloneData()
	grouped := groupToMap(data, h.groups)
	for _, attr := range attrs {
		attrToMap(grouped, attr)
	}

	return &CoralogixHandler{
		cxLogger: h.cxLogger,
		opts:     h.opts,

		data:   data,
		groups: h.cloneGroups(),
	}
}

// WithGroup returns a new Coralogix with a group, provided the group's name.
func (h *CoralogixHandler) WithGroup(name string) slog.Handler {
	return &CoralogixHandler{
		cxLogger: h.cxLogger,
		opts:     h.opts,

		data:   h.cloneData(),
		groups: append(h.cloneGroups(), name),
	}
}

// Enabled reports whether the logger emits log records at the given context and level.
// Note: We handover the decision down to the next handler.
func (h *CoralogixHandler) Enabled(ctx context.Context, level slog.Level) bool {
	minLevel := slog.LevelInfo
	if h.opts.Level != nil {
		minLevel = h.opts.Level.Level()
	}
	return level >= minLevel
}

func (h *CoralogixHandler) Stop() {
	h.cxLogger.Destroy()
}

func attrToMap(m map[string]any, a slog.Attr) {
	switch v := a.Value.Any().(type) {
	case []slog.Attr:
		m2 := map[string]any{}
		for _, a2 := range v {
			attrToMap(m2, a2)
			m[a.Key] = m2
		}
	default:
		m[a.Key] = v
	}
}

func groupToMap(m map[string]any, groups []string) map[string]any {
	for _, group := range groups {
		if _, ok := m[group]; !ok {
			m[group] = map[string]any{}
		}
		m = m[group].(map[string]any)
	}
	return m
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
