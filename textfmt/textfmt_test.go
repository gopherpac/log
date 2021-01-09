package textfmt

import (
	"reflect"
	"testing"
)

type result struct {
	fields  []string
	indents []string
}

type fixture struct {
	format   string
	expected result
}

var fixtures = []fixture{
	{
		format: "%{message}",
		expected: result{
			fields:  []string{"message"},
			indents: []string{"", ""},
		},
	},
	{
		format: "> %{time} %{gid} [%{level}]: %{message} %{ctx} %{fields} <",
		expected: result{
			fields:  []string{"time", "gid", "level", "message", "ctx", "fields"},
			indents: []string{"> ", " ", " [", "]: ", " ", " ", " <"},
		},
	},
}

func TestParse(t *testing.T) {
	for _, fx := range fixtures {
		t.Run("", func(t *testing.T) {
			fields, indents := Parse(fx.format)
			if !reflect.DeepEqual(fields, fx.expected.fields) {
				t.Errorf(
					"fields not match: \n\texpected: %v\n\tactual: %v\n",
					fx.expected.fields,
					fields,
				)
				return
			}
			if !reflect.DeepEqual(indents, fx.expected.indents) {
				t.Errorf(
					"indents not match: \n\texpected: %v\n\tactual: %v\n",
					fx.expected.indents,
					indents,
				)
				return
			}
		})
	}
}
