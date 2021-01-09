package log

// import (
// 	"bytes"
// 	"fmt"
// 	"strconv"
// 	"testing"
// 	"time"
// 	_ "unsafe"

// 	"go.chromium.org/luci/common/runtime/goroutine"
// )

// func str(b []byte) string {
// 	n := bytes.IndexByte(b, 0)
// 	if n == -1 {
// 		return string(b)
// 	}
// 	return string(b[:n])
// }

// func getbuf() []byte {
// 	b := make([]byte, 1<<10)
// 	b = b[:0]
// 	return b
// }

// func expectName(expected string, field TextMarshaler) func(t *testing.T) {
// 	return func(t *testing.T) {
// 		if field.Name() != expected {
// 			t.Errorf("TextTime.Name() should return '%s'", expected)
// 			return
// 		}
// 	}
// }

// func TestTextTime(t *testing.T) {
// 	t.Run("Name", expectName("time", &TextTime{}))

// 	t.Run("Default", func(t *testing.T) {
// 		b := getbuf()

// 		f := TextTime{}
// 		f.Construct()

// 		b = f.Marshal(b, Entry{
// 			Fields: []Field{
// 				{
// 					Key: "time",
// 					Val: time.Now().UTC(),
// 				},
// 			},
// 		})

// 		if _, err := time.Parse(time.RFC3339, str(b)); err != nil {
// 			t.Errorf("unmarshal failed: %v", err)
// 			return
// 		}
// 	})

// 	t.Run("Format", func(t *testing.T) {
// 		format := time.RFC850

// 		b := getbuf()

// 		f := TextTime{Format: format}
// 		f.Construct()

// 		b = f.Marshal(b, Entry{
// 			Fields: []Field{
// 				{
// 					Key: "time",
// 					Val: time.Now().UTC(),
// 				},
// 			},
// 		})

// 		if _, err := time.Parse(format, str(b)); err != nil {
// 			t.Errorf("unmarshal failed: %v", err)
// 			return
// 		}
// 	})
// }

// func TestTextLevel(t *testing.T) {
// 	run := func(level int, expected string) func(t *testing.T) {
// 		return func(t *testing.T) {
// 			b := getbuf()

// 			f := TextLevel{}
// 			b = f.Marshal(b, Entry{Lvl: level})

// 			actual := str(b)
// 			if expected != actual {
// 				t.Errorf("expected level name: %s but was: %s", expected, actual)
// 				return
// 			}
// 		}
// 	}

// 	t.Run("Name", expectName("lvl", &TextLevel{}))

// 	t.Run("DEBU", run(DEBUG, "DEBU"))
// 	t.Run("INFO", run(INFO, "INFO"))
// 	t.Run("WARN", run(WARNING, "WARN"))
// 	t.Run("ERRO", run(ERROR, "ERRO"))
// 	t.Run("CRIT", run(CRITICAL, "CRIT"))
// 	t.Run("FATA", run(FATAL, "FATA"))
// }

// func TestTextMessage(t *testing.T) {
// 	t.Run("Name", expectName("msg", &TextMessage{}))

// 	t.Run("Marshal", func(t *testing.T) {
// 		expected := "the test message!"
// 		b := getbuf()

// 		f := TextMessage{}
// 		b = f.Marshal(b, Entry{Msg: expected})

// 		actual := str(b)
// 		if expected != actual {
// 			t.Errorf("expected entry message '%s' but was '%s'", expected, actual)
// 			return
// 		}
// 	})
// }

// func TestTextGID(t *testing.T) {
// 	t.Run("Name", expectName("gid", &TextGID{}))

// 	t.Run("Marshal", func(t *testing.T) {
// 		expected := 1
// 		b := getbuf()

// 		f := TextGID{Index: 0}
// 		e := Entry{Fields: []Field{{Key: "", Val: goroutine.ID(1)}}}

// 		b = f.Marshal(b, e)

// 		actual := str(b)
// 		if strconv.Itoa(expected) != actual {
// 			t.Errorf("expected value '%d' but was '%s'", expected, actual)
// 			return
// 		}
// 	})
// }

// func TestTextFields(t *testing.T) {
// 	t.Run("Name", expectName("fields", &TextFields{}))

// 	t.Run("Marshal", func(t *testing.T) {
// 		fields := []Field{
// 			{
// 				Key: "request_id",
// 				Val: "123",
// 			},
// 			{
// 				Key: "user",
// 				Val: "dmr",
// 			},
// 			{
// 				Key: "rank",
// 				Val: 80,
// 			},
// 		}

// 		b := getbuf()

// 		f := TextFields{Index: 0}
// 		f.Construct()

// 		b = f.Marshal(b, Entry{Fields: fields})

// 		expected := "request_id=123 user=dmr rank=80"
// 		actual := str(b)

// 		if expected != actual {
// 			t.Errorf("expected '%s' but was '%s'", expected, actual)
// 			return
// 		}
// 	})
// }

// func TestTextFormat(t *testing.T) {
// 	now := time.Now().UTC()

// 	f := TextFormat{
// 		Format: "> %{time}	[%{gid}]	%{lvl}:	%{msg}	{%{fields}}\n",
// 		Fields: []TextMarshaler{
// 			&TextTime{Index: 1},
// 			&TextGID{Index: 0},
// 			&TextLevel{},
// 			&TextMessage{},
// 			&TextFields{Index: 2},
// 		},
// 	}
// 	f.Construct()

// 	b := getbuf()
// 	e := Entry{
// 		Lvl: INFO,
// 		Msg: "the test message!",
// 		Fields: []Field{
// 			{
// 				Key: "gid",
// 				Val: goroutine.ID(1),
// 			},
// 			{
// 				Key: "time",
// 				Val: now,
// 			},
// 			{
// 				Key: "request_id",
// 				Val: "123",
// 			},
// 			{
// 				Key: "user",
// 				Val: "dmr",
// 			},
// 			{
// 				Key: "rank",
// 				Val: 80,
// 			},
// 		},
// 	}

// 	b = f.Marshal(b, e)

// 	expected := fmt.Sprintf(
// 		"> %s	[1]	INFO:	%s	{request_id=123 user=dmr rank=80}\n",
// 		now.Format(time.RFC3339),
// 		e.Msg,
// 	)

// 	actual := string(b)
// 	if expected != actual {
// 		t.Errorf("\nexpected '%s'\nbut was '%s'", expected, actual)
// 		return
// 	}
// }
