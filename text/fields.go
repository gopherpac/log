package text

import (
	"time"

	"github.com/gopherpac/log/encode"
)

// Bool writes boolean to the buffer.
func Bool(b []byte, v interface{}) []byte {
	return encode.Bool(b, v.(bool))
}

// String writes string to the buffer.
func String(b []byte, v interface{}) []byte {
	return encode.String(b, v.(string))
}

// Int writes int to the buffer.
func Int(b []byte, v interface{}) []byte {
	return encode.Int(b, v.(int))
}

// Int8 writes int8 to the buffer.
func Int8(b []byte, v interface{}) []byte {
	return encode.Int8(b, v.(int8))
}

// Int16 writes int16 to the buffer.
func Int16(b []byte, v interface{}) []byte {
	return encode.Int16(b, v.(int16))
}

// Int32 writes int32 to the buffer.
func Int32(b []byte, v interface{}) []byte {
	return encode.Int32(b, v.(int32))
}

// Int64 writes int64 to the buffer.
func Int64(b []byte, v interface{}) []byte {
	return encode.Int64(b, v.(int64))
}

// Uint writes uint to the buffer.
func Uint(b []byte, v interface{}) []byte {
	return encode.Uint(b, v.(uint))
}

// Uint8 writes uint8 to the buffer.
func Uint8(b []byte, v interface{}) []byte {
	return encode.Uint8(b, v.(uint8))
}

// Uint16 writes uint16 to the buffer.
func Uint16(b []byte, v interface{}) []byte {
	return encode.Uint16(b, v.(uint16))
}

// Uint32 writes uint32 to the buffer.
func Uint32(b []byte, v interface{}) []byte {
	return encode.Uint32(b, v.(uint32))
}

// Uint64 writes uint64 to the buffer.
func Uint64(b []byte, v interface{}) []byte {
	return encode.Uint64(b, v.(uint64))
}

// Float32 writes float32 to the buffer.
func Float32(b []byte, v interface{}) []byte {
	return encode.Float32(b, v.(float32))
}

// Float64 writes float64 to the buffer.
func Float64(b []byte, v interface{}) []byte {
	return encode.Float64(b, v.(float64))
}

// Time returns writer of time.Time to the buffer.
func Time(format string) func(b []byte, v interface{}) []byte {
	return func(b []byte, v interface{}) []byte {
		return encode.Time(b, v.(time.Time), format)
	}
}
