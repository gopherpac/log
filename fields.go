package log

import (
	"runtime"
	"time"

	"github.com/gopherpac/log/encode"
	"go.chromium.org/luci/common/runtime/goroutine"
)

// DefaultTextFormat defines default template of plain text log entry.
const DefaultTextFormat = "%{time} [%{lvl}]: %{msg} %{ctx} %{fields}"

// EntryWriter writes a part of entry information to the intermediate buffer,
// it can either write the the full entry ot part or just additional info like
// time, location, goroutine ID, etc. The responsibility if EntryWriter is to
// take some peace of information from log entry and pass to a formatting
// function.
type EntryWriter interface {
	Write([]byte, Entry, []Field) []byte
}

// Level writes the entry level to the buffer provided.
type Level struct {
	Fmt func([]byte, string) []byte
}

// Construct sets the default log level formatter if nothing was provided.
func (l *Level) Construct() {
	if l.Fmt == nil {
		l.Fmt = func(b []byte, v string) []byte {
			return encode.String(b, "lvl: ???")
		}
	}
}

func (l *Level) Write(b []byte, e Entry, f []Field) []byte {
	return l.Fmt(b, lvl2str(e.Lvl))
}

// Message writes the entry message to the buffer provided.
type Message struct {
	Fmt func([]byte, string) []byte
}

// Construct sets the default message formatter if nothing was provided.
func (m *Message) Construct() {
	if m.Fmt == nil {
		// panic("log.Message.Fmt not set")
		m.Fmt = func(b []byte, v string) []byte {
			return encode.String(b, "msg: ???")
		}
	}
}

func (m *Message) Write(b []byte, e Entry, f []Field) []byte {
	return m.Fmt(b, e.Msg)
}

// Timestamp writes the log entry timestamp in the format provided.
type Timestamp struct {
	Now  func() time.Time
	Kind string
	Fmt  func([]byte, time.Time, string) []byte
}

// Construct sets default values for the time writer.
func (t *Timestamp) Construct() {
	if t.Now == nil {
		t.Now = func() time.Time { return time.Now().UTC() }
	}
	if t.Kind == "" {
		t.Kind = time.RFC3339
	}
	if t.Fmt == nil {
		t.Fmt = func(b []byte, v time.Time, format string) []byte {
			return encode.String(b, "time: ???")
		}
	}
}

func (t *Timestamp) Write(b []byte, e Entry, f []Field) []byte {
	return t.Fmt(b, t.Now(), t.Kind)
}

// Location writes the file:line information about the place from where the
// log method was called.
type Location struct {
	Depth int
	Fmt   func([]byte, string, int) []byte
}

// Construct sets the default Depth and Fmt if some field was not initialied.
func (l *Location) Construct() {
	if l.Depth == 0 {
		l.Depth = 2
	}
	if l.Fmt == nil {
		l.Fmt = func(b []byte, file string, line int) []byte {
			return encode.String(b, "loc: ???")
		}
	}
}

func (l *Location) Write(b []byte, e Entry, f []Field) []byte {
	// copied from https://golang.org/src/log/log.go
	_, file, line, ok := runtime.Caller(l.Depth)
	if !ok {
		file = "???"
		line = 0
	}
	return l.Fmt(b, file, line)
}

// GoroutineID writes ID of the goroutine from where the log method was called.
type GoroutineID struct {
	Fmt func([]byte, uint64) []byte
}

// Construct sets the default value for the Fmt field if nothing was provided.
func (g *GoroutineID) Construct() {
	if g.Fmt == nil {
		g.Fmt = func(b []byte, v uint64) []byte {
			return encode.String(b, "gid: ???")
		}
	}
}

func (g *GoroutineID) Write(b []byte, e Entry, f []Field) []byte {
	return g.Fmt(b, uint64(goroutine.CurID()))
}

// Context writes logger context fields.
type Context struct {
	Fmt func([]byte, []Field) []byte
}

func (c *Context) Write(b []byte, e Entry, f Field) []byte {
	fields := e.Ctx.Value(ckey)
	if fields == nil {
		return b
	}

	return c.Fmt(b, fields.([]Field))
}

// Fields writes fields of the log entry.
// type Fields struct {
// 	Fmt func([]byte, []Field) []byte
// }

// func (fs *Fields) Write(b []byte, e Entry, f []Field) []byte {
// 	return fs.Fmt(b, f)
// }

func lvl2str(lvl int) string {
	switch lvl {
	case DEBUG:
		return "DEBU"
	case INFO:
		return "INFO"
	case WARNING:
		return "WARN"
	case ERROR:
		return "ERRO"
	case CRITICAL:
		return "CRIT"
	case FATAL:
		return "FATA"
	default:
		return "????"
	}
}
