package log

import (
	"io"
	"sync"

	"github.com/gopherpac/log/encode"
	"github.com/gopherpac/log/textfmt"
)

// TextEncoders is the default text encoders table.
var TextEncoders = []Encoder{}

// Named object can tell self name.
type Named interface {
	Name() string
}

// Name returns the name of time field used in format stirng.
func (t *Timestamp) Name() string {
	return "time"
}

// Name returns the name of level field used in format string.
func (t *Level) Name() string {
	return "lvl"
}

// Name returns the name of message field used in format string.
func (t *Message) Name() string {
	return "msg"
}

// Name returns the name of goroutine ID field used in format string.
func (t *GoroutineID) Name() string {
	return "gid"
}

// Name returns the name of fields log field used in format string.
// func (t *Fields) Name() string {
// 	return "fields"
// }

// Text formats a log entry in plain text.
type Text struct {
	Format string
	Fields []EntryWriter

	indents []string
	fields  []int

	buffers *sync.Pool
}

// Construct initializes the Text formatter.
func (t *Text) Construct() {
	fields, indents := textfmt.Parse(t.Format)

	t.indents = indents
	for _, name := range fields {
		k := -1
		for i, f := range t.Fields {
			named, ok := f.(Named)
			if !ok {
				continue
			}
			if named.Name() == name {
				k = i
				break
			}
		}
		if k != -1 {
			t.fields = append(t.fields, k)
		}
	}

	for _, f := range t.Fields {
		if c, ok := f.(Constructor); ok {
			c.Construct()
		}
	}

	t.buffers = &sync.Pool{New: func() interface{} {
		b := make([]byte, 1<<10)
		return &b
	}}
}

func (t *Text) Write(w io.Writer, e Entry, fields []Field) error {
	b := t.buffers.Get().(*[]byte)
	*b = (*b)[:0]
	*b = encode.String(*b, t.indents[0])

	for i := range t.fields {
		*b = t.Fields[t.fields[i]].Write(*b, e, fields)
		*b = encode.String(*b, t.indents[i+1])
	}

	_, err := w.Write(*b)
	t.buffers.Put(b)

	return err
}

// TextFields creates fields writer in text format. Fields will be written like
// %{f1.Key}%{delim}%{f1.Val}%{sep}%{f2.Key}%{delim}%{f2.Val}...
func TextFields(sep, delim string) func([]byte, []Field) []byte {
	return func(b []byte, fields []Field) []byte {
		if len(fields) < 1 {
			return b
		}

		last := len(fields) - 1
		for _, f := range fields[:last] {
			b = encode.String(b, f.Key)
			b = encode.String(b, delim)
			// b = f.Fmt(b, f.Val)
			b = encode.String(b, sep)
		}

		l := fields[last]
		b = encode.String(b, l.Key)
		b = encode.String(b, delim)
		// b = l.Fmt(b, l.Val)

		return b
	}
}
