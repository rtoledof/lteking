package slogctx

import (
	"context"
	"log/slog"
)

type contextKey int

const (
	attrContextKey = contextKey(iota + 1)
)

type entry struct {
	key string
	val any
}

func With(ctx context.Context, key string, val any) context.Context {
	v := ctx.Value(attrContextKey)
	if v == nil {
		return context.WithValue(ctx, attrContextKey, []entry{{key, val}})
	}
	attrs := v.([]entry)

	return context.WithValue(ctx, attrContextKey, append(attrs, entry{key, val}))

}

type ctxHandler struct {
	handler slog.Handler
}

func NewContextAttrsHandler(h slog.Handler) *ctxHandler {
	return &ctxHandler{h}
}

func (h *ctxHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h *ctxHandler) Handle(ctx context.Context, r slog.Record) error {
	if ctx == nil {
		return h.handler.Handle(ctx, r)
	}
	v := ctx.Value(attrContextKey)
	if v != nil {
		attrs := v.([]entry)
		for _, v := range attrs {
			r.AddAttrs(slog.Any(v.key, v.val))
		}
	}
	return h.handler.Handle(ctx, r)
}

func (h *ctxHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return NewContextAttrsHandler(h.handler.WithAttrs(attrs))
}

func (h *ctxHandler) WithGroup(name string) slog.Handler {
	return NewContextAttrsHandler(h.handler.WithGroup(name))
}

func (h *ctxHandler) Handler() slog.Handler {
	return h.handler
}
