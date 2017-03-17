package log

import (
	"context"
	"fmt"
	"testing"
)

func TestWithGID(t *testing.T) {
	gid := GID()
	ctx := context.Background()

	ctx = WithGID(ctx)
	got := ctx.Value(gidKey).(string)
	if got != gid {
		t.Fatal(got)
	}

	ctx = WithGID(ctx)
	got = ctx.Value(gidKey).(string)
	if got != gid+"."+gid {
		t.Fatal(got)
	}
}

func TestDebug(t *testing.T) {
	buf := new(debuf)
	lgr := New(buf)
	ctx := context.Background()

	lgr.Debug(ctx, "msg", "test")
	if buf.msg != "" {
		t.Fatal(*buf)
	}

	ctx = WithDebug(ctx)

	lgr.Debug(ctx, "msg", "test")
	if buf.fname != "TestDebug" {
		t.Fatal(buf.fname)
	}
}

type debuf struct {
	fname string
	msg   string
}

func (s *debuf) Record(ctx context.Context, knd kind, msg string, keyvals []interface{}) (err error) {
	kv := keyvals[len(keyvals)-6:]

	branch, gid, file, line, _, fname := kv[0], kv[1], kv[2], kv[3], kv[4], kv[5]
	s.msg = fmt.Sprintf(" [%s.%s] %s:%d ", branch, gid, file, line)
	s.fname = fname.(string)

	return nil
}
