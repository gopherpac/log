package log

import (
	"context"
	"io"
	"io/ioutil"
	"reflect"
	"testing"
)

type memstream struct {
	buf []byte
}

func (m *memstream) Write(b []byte) (int, error) {
	m.buf = append(m.buf, b...)
	return len(b), nil
}

type mementry struct {
	ctx    context.Context
	lvl    int
	msg    string
	fields []Field
}

type memwriter struct {
	entries []mementry
}

func (m *memwriter) Write(w io.Writer, e Entry, fields []Field) error {
	// m.entries = append(m.entries, mementry{
	// 	msg:    e.Msg,
	// 	lvl:    e.Lvl,
	// 	fields: fields,
	// })
	return nil
}

func TestWith(t *testing.T) {
	t.Run("Null", func(t *testing.T) {
		c := context.Background()
		c = With(c, Field{Key: "k1", Val: "v1"})

		expected := []Field{
			{
				Key: "k1",
				Val: "v1",
			},
		}
		actual := c.Value(ckey)

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v but was %v", expected, actual)
			return
		}
	})

	t.Run("NotNull", func(t *testing.T) {
		c := context.Background()
		c = With(c, Field{Key: "k1", Val: "v1"})
		c = With(c, Field{Key: "k2", Val: "v2"})

		expected := []Field{
			{
				Key: "k1",
				Val: "v1",
			},
			{
				Key: "k2",
				Val: "v2",
			},
		}
		actual := c.Value(ckey)

		if !reflect.DeepEqual(expected, actual) {
			t.Errorf("expected %v but was %v", expected, actual)
			return
		}
	})
}

func TestLog(t *testing.T) {
	level := func(lvl int, expected []mementry) func(t *testing.T) {
		return func(t *testing.T) {
			s := ioutil.Discard
			w := &memwriter{}

			l := New(Config{
				Level:  lvl,
				Stream: s,
				Writer: w,
			})

			l.Debug("debug message")
			l.Info("info message")
			l.Warn("warning message")
			l.Error("error message")
			l.Critical("critical message")
			l.Fatal("fatal message")

			actual := w.entries

			if !reflect.DeepEqual(expected, actual) {
				t.Errorf("expected %v but was %v", expected, actual)
				return
			}
		}
	}

	t.Run("DEBU", level(DEBUG, []mementry{
		{
			lvl: DEBUG,
			msg: "debug message",
		},
		{
			lvl: INFO,
			msg: "info message",
		},
		{
			lvl: WARNING,
			msg: "warning message",
		},
		{
			lvl: ERROR,
			msg: "error message",
		},
		{
			lvl: CRITICAL,
			msg: "critical message",
		},
		{
			lvl: FATAL,
			msg: "fatal message",
		},
	}))

	t.Run("INFO", level(INFO, []mementry{
		{
			lvl: INFO,
			msg: "info message",
		},
		{
			lvl: WARNING,
			msg: "warning message",
		},
		{
			lvl: ERROR,
			msg: "error message",
		},
		{
			lvl: CRITICAL,
			msg: "critical message",
		},
		{
			lvl: FATAL,
			msg: "fatal message",
		},
	}))

	t.Run("WARN", level(WARNING, []mementry{
		{
			lvl: WARNING,
			msg: "warning message",
		},
		{
			lvl: ERROR,
			msg: "error message",
		},
		{
			lvl: CRITICAL,
			msg: "critical message",
		},
		{
			lvl: FATAL,
			msg: "fatal message",
		},
	}))

	t.Run("ERRO", level(ERROR, []mementry{
		{
			lvl: ERROR,
			msg: "error message",
		},
		{
			lvl: CRITICAL,
			msg: "critical message",
		},
		{
			lvl: FATAL,
			msg: "fatal message",
		},
	}))

	t.Run("CRIT", level(CRITICAL, []mementry{
		{
			lvl: CRITICAL,
			msg: "critical message",
		},
		{
			lvl: FATAL,
			msg: "fatal message",
		},
	}))

	t.Run("FATA", level(FATAL, []mementry{
		{
			lvl: FATAL,
			msg: "fatal message",
		},
	}))
}
