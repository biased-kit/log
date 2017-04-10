package log

import (
	"context"
	"testing"

	"github.com/biased-kit/errors"
	opentracing "github.com/opentracing/opentracing-go"
	"github.com/opentracing/opentracing-go/log"
)

func TestDefaultDebug(t *testing.T) {
	spn, ctx := ctxWithSpan()
	Debug(ctx, "test msg", "key1", "value1")
	if len(spn.vals) != 0 {
		t.Fail()
	}
}

func TestDebug(t *testing.T) {
	spn, ctx := ctxWithSpan()
	ctx = WithDebug(ctx)

	Debug(ctx, "test msg", "key1", "value1")
	if spn.vals["msg"] != "test msg" {
		t.Fail()
	}
}

func TestError(t *testing.T) {
	e := errors.New("test_err")
	spn, ctx := ctxWithSpan()
	Error(ctx, e)

	_, ok := spn.vals["stack"]
	if !ok {
		t.Fail()
	}
}

func TestNoError(t *testing.T) {

	spn, ctx := ctxWithSpan()
	Error(ctx, nil)

	if len(spn.vals) > 1 || spn.vals["lvl"] != "error" {
		t.Fatalf("unxepected keyvals: %v", spn.vals)
	}
}

func TestConextWithoutSpan(t *testing.T) {
	// should not panic
	ctx := context.Background()
	Debug(ctx, "test msg")

	e := errors.New("test err")
	Error(ctx, e)
}

func ctxWithSpan() (*span, context.Context) {
	ctx := context.Background()
	spn := &span{}
	ctx = opentracing.ContextWithSpan(ctx, spn)
	return spn, ctx
}

type span struct {
	keyvals []interface{}
	vals    map[string]interface{}
}

func (s *span) LogKV(alternatingKeyValues ...interface{}) {
	s.keyvals = alternatingKeyValues
	s.vals = kv2map(s.keyvals)
}

func kv2map(kv []interface{}) map[string]interface{} {
	m := map[string]interface{}{}
	for i := 0; i < len(kv); i += 2 {
		m[kv[i].(string)] = kv[i+1]
	}
	return m
}

func (*span) Finish() {}

func (*span) FinishWithOptions(opts opentracing.FinishOptions) {}
func (*span) Context() opentracing.SpanContext {
	return nil
}

func (*span) SetOperationName(operationName string) opentracing.Span {
	return nil
}
func (*span) SetTag(key string, value interface{}) opentracing.Span {
	return nil
}

func (*span) LogFields(fields ...log.Field) {}

func (*span) SetBaggageItem(restrictedKey, value string) opentracing.Span {
	return nil
}
func (*span) BaggageItem(restrictedKey string) string {
	return ""
}
func (*span) Tracer() opentracing.Tracer {
	return nil
}

func (*span) LogEvent(event string)                                 {}
func (*span) LogEventWithPayload(event string, payload interface{}) {}
func (*span) Log(data opentracing.LogData)                          {}
