package cannon

import (
	"context"
	"log/slog"
	"runtime"
	"time"
)

type contextKey uint8

const (
	loggerContextKey contextKey = iota
)

func NewContextWithLogger(ctx context.Context, l *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerContextKey, l)
}

func LoggerFromContext(ctx context.Context) *slog.Logger {
	v, ok := ctx.Value(loggerContextKey).(*slog.Logger)
	if !ok {
		return nil
	}
	return v
}

// CanonicalHandler wraps a Handler with an Enabled method
// that returns false for levels below a minimum.
type CanonicalHandler struct {
	handler slog.Handler
	attrs   *[]slog.Attr
}

// NewCanonicalHandler returns a canonical handler.
// All methods except Enabled delegate to handler.
func NewCanonicalHandler(h slog.Handler) *CanonicalHandler {
	// Optimization: avoid chains of LevelHandlers.
	if lh, ok := h.(*CanonicalHandler); ok {
		h = lh.Handler()
	}
	a := []slog.Attr{}
	return &CanonicalHandler{h, &a}
}

// Enabled implements Handler.Enabled by reporting whether
// level is at least as large as h's level.
func (h *CanonicalHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

// Handle implements Handler.Handle.
func (h *CanonicalHandler) Handle(ctx context.Context, r slog.Record) error {
	r.Attrs(func(a slog.Attr) bool {
		*h.attrs = append(*h.attrs, a)
		return true
	})
	return h.handler.Handle(ctx, r)
}

// WithAttrs implements Handler.WithAttrs.
func (h *CanonicalHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	v := NewCanonicalHandler(h.handler.WithAttrs(attrs))
	v.attrs = h.attrs
	*h.attrs = append(*h.attrs, attrs...)

	return v

}

// WithGroup implements Handler.WithGroup.
func (h *CanonicalHandler) WithGroup(name string) slog.Handler {
	v := NewCanonicalHandler(h.handler.WithGroup(name))
	v.attrs = h.attrs
	return v

}

// Handler returns the Handler wrapped by h.
func (h *CanonicalHandler) Handler() slog.Handler {
	return h.handler
}

func (h *CanonicalHandler) resetAttrs() {
	*h.attrs = (*h.attrs)[:0]
}

type Logger struct {
	l *slog.Logger
}

func (l *Logger) Logger() *slog.Logger {
	return l.l
}

// Emit logs a canonical log line with all attributes and resets the logger.
// This will remove all attributes from l.
// You can provide additional attributes to be logged.
// Emit should not be called by children loggers.
func (l *Logger) Emit(args ...any) error {
	if !l.l.Enabled(context.Background(), slog.LevelInfo) {
		return nil
	}
	h := l.l.Handler().(*CanonicalHandler)
	var pcs [1]uintptr
	runtime.Callers(2, pcs[:]) // skip [Callers, Infof]
	rec := slog.NewRecord(time.Now(), slog.LevelInfo, "canonical-log-line", pcs[0])
	rec.AddAttrs(*h.attrs...)
	rec.Add(args...)
	defer h.resetAttrs()
	return h.handler.Handle(context.Background(), rec)

}

func NewLogger(l *slog.Logger) (*Logger, func()) {
	h := NewCanonicalHandler(l.Handler())
	logger := Logger{
		l: slog.New(h),
	}
	return &logger, func() {
		h.resetAttrs()
	}
}
