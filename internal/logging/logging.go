package logging

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"
	"runtime"
	"strings"
)

type keyType int

const key = keyType(0)

type LoggerImpl struct {
	next slog.Handler
}

func NewLoggerImpl(next slog.Handler) *LoggerImpl {
	return &LoggerImpl{next: next}
}

func (h *LoggerImpl) Enabled(ctx context.Context, rec slog.Level) bool {
	return h.next.Enabled(ctx, rec)
}

func (h *LoggerImpl) Handle(ctx context.Context, rec slog.Record) error {
	if c, ok := ctx.Value(key).(logCtx); ok {
		v := reflect.ValueOf(c)
		t := v.Type()

		for i := 0; i < v.NumField(); i++ {
			field := v.Field(i)
			if !field.IsZero() {
				fieldName := strings.ToLower(t.Field(i).Name)
				rec.Add(fieldName, field.Interface())
			}
		}
	}

	if rec.PC != 0 {
		fs := runtime.CallersFrames([]uintptr{rec.PC})
		f, _ := fs.Next()
		rec.Add("source", fmt.Sprintf("%s:%d", f.File, f.Line))
	}

	return h.next.Handle(ctx, rec)
}

func (h *LoggerImpl) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &LoggerImpl{next: h.next.WithAttrs(attrs)}
}

func (h *LoggerImpl) WithGroup(name string) slog.Handler {
	return &LoggerImpl{next: h.next.WithGroup(name)}
}
