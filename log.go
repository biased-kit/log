package log

import (
	"context"
	"fmt"
	"runtime"
	"strings"
	"time"
)

type kind int

const (
	debug = iota
	info
	warn
	erro
)

type extraError interface {
	error
	StackTrace() *runtime.Frames
	KeyvalPairs() []interface{}
}

var DefaultLogger Logger = New(StdoutRecorder{})

// Logger provides several methods to story your message,
// Each method coresponds the purpose the message should serve.
type Logger interface {
	// Error records application(transaction) errors. It Never should be used inside the public library.
	Error(ctx context.Context, err error) error
	// Warn records errors that could be handled. It is signals that app/lib doesn't work as expected, and could do it better.
	Warn(ctx context.Context, err error) error
	// Info is useful result of your app. Obvious ussage is for console util. the Servers could prints the transaction log.
	// Only for App, not for lib.
	Info(ctx context.Context, msg string, keyvals ...interface{}) error
	// Debug print usful info that help to resolve the issue. Debug = Verbose.
	Debug(ctx context.Context, msg string, keyvals ...interface{}) error
}

// Recorder write the log records. By default stdout recorder is used.
// you could define any other implementations.
type Recorder interface {
	Record(ctx context.Context, knd kind, msg string, keyvals []interface{}) (err error)
}

type StdoutRecorder struct{}

func (l StdoutRecorder) Record(ctx context.Context, knd kind, msg string, keyvals []interface{}) (err error) {

	tm := time.Now().UTC()

	switch knd {
	case debug:
		kv := keyvals[len(keyvals)-6:]
		branch, gid, file, line, _, _ := kv[0], kv[1], kv[2], kv[3], kv[4], kv[5]
		_, err = fmt.Printf("%v[DEBUG] [%s.%s] %s:%d %s; env: %v", tm, branch, gid, file, line, msg, keyvals[:len(keyvals)-6])
	case info:
		_, err = fmt.Printf("%v[INFO] %s; %v", tm, msg, keyvals)
	case warn:
		_, err = fmt.Printf("%v[WARN] %v; %v", tm, msg, keyvals)
	case erro:
		_, err = fmt.Printf("%v[ERROR] %v; %v", tm, msg, keyvals)

	default:
		err = fmt.Errorf("unsupported log type %v", knd)
	}

	return err
}

type StdoutLogger struct {
	Recorder
}

func New(recorder Recorder) Logger {
	return StdoutLogger{recorder}
}

func (l StdoutLogger) Debug(ctx context.Context, msg string, keyvals ...interface{}) (err error) {
	isLog, _ := ctx.Value(trailKey).(bool)
	if !isLog {
		return
	}

	branch := GIDValue(ctx)
	gid := GID()

	file, ffile, line := Caller(4)
	fname := Funcname(ffile)

	keyvals = append(keyvals, branch, gid, file, line, ffile, fname)

	return l.Record(ctx, debug, msg, keyvals)
}

func (l StdoutLogger) Info(ctx context.Context, msg string, keyvals ...interface{}) error {
	return l.Record(ctx, info, msg, keyvals)
}

func (l StdoutLogger) Warn(ctx context.Context, err error) error {
	msg, keyvals := err2msg(err)
	return l.Record(ctx, warn, msg, keyvals)
}

func (l StdoutLogger) Error(ctx context.Context, err error) error {
	msg, keyvals := err2msg(err)
	return l.Record(ctx, erro, msg, keyvals)
}

func err2msg(err error) (string, []interface{}) {
	msg := err.Error()
	var keyvals []interface{}
	if e, ok := err.(extraError); ok {
		keyvals = append(keyvals, "stack", e.StackTrace())
		keyvals = append(keyvals, e.KeyvalPairs())
	}
	return msg, keyvals
}

func GIDValue(ctx context.Context) string {
	branch, _ := ctx.Value(gidKey).(string)
	return branch
}

// GID returns the goroutine number
// the caller is running on.
func GID() string {
	b := make([]byte, 20)
	runtime.Stack(b, false)
	l := len("goroutine ")
	for i := l; i < len(b); i++ {
		if b[i] == ' ' {
			return string(b[l:i])
		}
	}
	return "unknown"
}

func Caller(depth int) (file string, function string, line int) {
	var rpc [1]uintptr
	if runtime.Callers(depth, rpc[:]) < 1 {
		return
	}

	f := runtime.FuncForPC(rpc[0])
	file, line = f.FileLine(rpc[0])

	return file, f.Name(), line
}

// funcname removes the path prefix component of a function's name reported by func.Name().
func Funcname(name string) string {
	i := strings.LastIndex(name, "/")
	name = name[i+1:]
	i = strings.Index(name, ".")
	return name[i+1:]
}

//*************************** Contex Helpers **************************//

// gidKey serves to store gid chain in context
// use string public type because of possible vendoring
const gidKey string = "!@gidKey@!"

// trailKey is used to store control flag in context.
// the flag controls whether the trail should be executed
const trailKey string = "!@trailKey@!"

// WithDebug create a context that activates trail Log
func WithDebug(parent context.Context) context.Context {
	ctx := context.WithValue(parent, trailKey, true)
	return ctx
}

// WithGID adds to context current goroutine id.
// you could call it before passing the context to new goroutine, in order to display "goroutines branch"
func WithGID(parent context.Context) context.Context {
	var gid string
	val := parent.Value(gidKey)
	if val != nil {
		gid, _ = val.(string)
		gid += "."
	}

	gid = gid + GID()
	ctx := context.WithValue(parent, gidKey, gid)

	return ctx
}
