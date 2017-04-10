package log

import (
	"context"
	"fmt"
	"runtime"

	"github.com/biased-kit/errors"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/ext"
)

type debKeyType struct{}

var debKey debKeyType = debKeyType{}

// WithDebug create a context that activates trail Log
func WithDebug(parent context.Context) context.Context {
	ctx := context.WithValue(parent, debKey, struct{}{})
	return ctx
}

// Debug stores message details in current span, that extracted from context.
// The context should contain debug flag, see WithDebug()
// Will panic if ctx is nil
func Debug(ctx context.Context, msg string, keyvals ...interface{}) {
	if isLog := ctx.Value(debKey); isLog == nil {
		return
	}

	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}
	span.SetTag("debug", true)

	kv := make([]interface{}, 0, len(keyvals)+6)
	kv = append(kv, "lvl", "debug", "msg", msg)

	pc, _, _, ok := runtime.Caller(1)
	if ok {
		c := newCall(pc)
		kv = append(kv, "caller", fmt.Sprintf("%+v", c))
	}
	kv = append(kv, keyvals...)
	span.LogKV(kv...)
}

// Error stores error details in current span, that extracted from context.
// Will panic if ctx is nil
func Error(ctx context.Context, e errors.E) {
	span := opentracing.SpanFromContext(ctx)
	if span == nil {
		return
	}
	ext.Error.Set(span, true)
	if e == nil {
		span.LogKV("lvl", "error")
		return
	}

	keyvals := e.KeyValues()
	kv := make([]interface{}, 0, len(keyvals)+6)
	kv = append(kv, "lvl", "error")
	kv = append(kv, "msg", e.Error())

	stk := stack(e.StackTrace())
	kv = append(kv, "stack", fmt.Sprint(stk))
	kv = append(kv, keyvals...)
	span.LogKV(kv...)
}
