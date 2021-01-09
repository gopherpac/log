package encode

import (
	"fmt"
	"math"
	"strconv"
	"time"
)

func Bool(b []byte, v bool) []byte {
	return strconv.AppendBool(b, v)
}

func Bytes(b []byte, v []byte) []byte {
	return append(b, v...)
}

func String(b []byte, v string) []byte {
	return append(b, v...)
}

func Quote(b []byte, v string) []byte {
	return strconv.AppendQuote(b, v)
}

func Int(b []byte, v int) []byte {
	return strconv.AppendInt(b, int64(v), 10)
}

func Int8(b []byte, v int8) []byte {
	return strconv.AppendInt(b, int64(v), 10)
}

func Int16(b []byte, v int16) []byte {
	return strconv.AppendInt(b, int64(v), 10)
}

func Int32(b []byte, v int32) []byte {
	return strconv.AppendInt(b, int64(v), 10)
}

func Int64(b []byte, v int64) []byte {
	return strconv.AppendInt(b, int64(v), 10)
}

func Uint(b []byte, v uint) []byte {
	return strconv.AppendUint(b, uint64(v), 10)
}

func Uint8(b []byte, v uint8) []byte {
	return strconv.AppendUint(b, uint64(v), 10)
}

func Uint16(b []byte, v uint16) []byte {
	return strconv.AppendUint(b, uint64(v), 10)
}

func Uint32(b []byte, v uint32) []byte {
	return strconv.AppendUint(b, uint64(v), 10)
}

func Uint64(b []byte, v uint64) []byte {
	return strconv.AppendUint(b, uint64(v), 10)
}

func float(b []byte, v float64, bs int) []byte {
	v64 := float64(v)
	switch {
	case math.IsNaN(v64):
		return append(b, `"NaN"`...)
	case math.IsInf(v64, 1):
		return append(b, `"+Inf"`...)
	case math.IsInf(v64, -1):
		return append(b, `"-Inf"`...)
	default:
		return strconv.AppendFloat(b, float64(v), 'f', -1, bs)
	}
}

func Float32(b []byte, v float32) []byte {
	return float(b, float64(v), 32)
}

func Float64(b []byte, v float64) []byte {
	return float(b, v, 64)
}

func Time(b []byte, v time.Time, format string) []byte {
	return v.AppendFormat(b, format)
}

type writer []byte

func (w *writer) Write(b []byte) (int, error) {
	dst := (*[]byte)(w)
	*dst = append(*dst, b...)
	return len(b), nil
}

func Value(b []byte, v interface{}) []byte {
	switch t := v.(type) {
	case bool:
		return Bool(b, t)
	case int:
		return Int(b, t)
	case int8:
		return Int8(b, t)
	case int16:
		return Int16(b, t)
	case int32:
		return Int32(b, t)
	case int64:
		return Int64(b, t)
	case uint:
		return Uint(b, t)
	case uint8:
		return Uint8(b, t)
	case uint16:
		return Uint16(b, t)
	case uint32:
		return Uint32(b, t)
	case uint64:
		return Uint64(b, t)
	case float32:
		return Float32(b, t)
	case float64:
		return Float64(b, t)
	case string:
		return String(b, t)
	default:
		w := writer(b)
		fmt.Fprintf(&w, "%v", v)
		return b
	}
}
