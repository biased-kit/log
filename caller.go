package log

import (
	"bytes"
	"fmt"
	"io"
	"runtime"
	"strconv"
	"strings"
)

type stack []uintptr

func (s stack) String() string {
	var buf bytes.Buffer
	for _, pc := range s {
		if pc == 0 {
			break
		}
		c := newCall(pc)
		buf.WriteString(fmt.Sprintf("%+v %n\n", c, c))

	}
	return buf.String()
}

// Call records a single function invocation from a goroutine stack.
type call struct {
	fn *runtime.Func
	pc uintptr
}

// NewCall calls
func newCall(pc uintptr) call {
	fn := runtime.FuncForPC(pc)
	return call{fn: fn, pc: pc}
}

// Format implements fmt.Formatter with support for the following verbs.
//
//    %s    source file
//    %d    line number
//    %n    function name
//    %v    equivalent to %s:%d
//
// It accepts the '+' and '#' flags for most of the verbs as follows.
//
//    %+s   path of source file relative to the compile time GOPATH
//    %#s   full path of source file
//    %+n   import path qualified function name
//    %+v   equivalent to %+s:%d
//    %#v   equivalent to %#s:%d
func (c call) Format(s fmt.State, verb rune) {
	if c.fn == nil {
		fmt.Fprintf(s, "%%!%c(NOFUNC)", verb)
		return
	}

	switch verb {
	case 's', 'v':
		file, line := c.fn.FileLine(c.pc)
		switch {
		case s.Flag('#'):
			// done
		case s.Flag('+'):
			file = file[pkgIndex(file, c.fn.Name()):]
		default:
			const sep = "/"
			if i := strings.LastIndex(file, sep); i != -1 {
				file = file[i+len(sep):]
			}
		}
		io.WriteString(s, file)
		if verb == 'v' {
			buf := [7]byte{':'}
			s.Write(strconv.AppendInt(buf[:1], int64(line), 10))
		}

	case 'd':
		_, line := c.fn.FileLine(c.pc)
		buf := [6]byte{}
		s.Write(strconv.AppendInt(buf[:0], int64(line), 10))

	case 'n':
		name := c.fn.Name()
		if !s.Flag('+') {
			const pathSep = "/"
			if i := strings.LastIndex(name, pathSep); i != -1 {
				name = name[i+len(pathSep):]
			}
			const pkgSep = "."
			if i := strings.Index(name, pkgSep); i != -1 {
				name = name[i+len(pkgSep):]
			}
		}
		io.WriteString(s, name)
	}
}

// pkgIndex returns the index that results in file[index:] being the path of
// file relative to the compile time GOPATH, and file[:index] being the
// $GOPATH/src/ portion of file. funcName must be the name of a function in
// file as returned by runtime.Func.Name.
func pkgIndex(file, funcName string) int {
	// As of Go 1.6.2 there is no direct way to know the compile time GOPATH
	// at runtime, but we can infer the number of path segments in the GOPATH.
	// We note that runtime.Func.Name() returns the function name qualified by
	// the import path, which does not include the GOPATH. Thus we can trim
	// segments from the beginning of the file path until the number of path
	// separators remaining is one more than the number of path separators in
	// the function name. For example, given:
	//
	//    GOPATH     /home/user
	//    file       /home/user/src/pkg/sub/file.go
	//    fn.Name()  pkg/sub.Type.Method
	//
	// We want to produce:
	//
	//    file[:idx] == /home/user/src/
	//    file[idx:] == pkg/sub/file.go
	//
	// From this we can easily see that fn.Name() has one less path separator
	// than our desired result for file[idx:]. We count separators from the
	// end of the file path until it finds two more than in the function name
	// and then move one character forward to preserve the initial path
	// segment without a leading separator.
	const sep = "/"
	i := len(file)
	for n := strings.Count(funcName, sep) + 2; n > 0; n-- {
		i = strings.LastIndex(file[:i], sep)
		if i == -1 {
			i = -len(sep)
			break
		}
	}
	// get back to 0 or trim the leading separator
	return i + len(sep)
}
