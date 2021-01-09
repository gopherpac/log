package log

import (
	"context"
	"fmt"
	"io"
	"os"
	"sync"
)

const (
	// DEBUG defines debug log level.
	DEBUG = 0
	// INFO defines info log level.
	INFO = 1
	// WARNING defines warning log level.
	WARNING = 2
	// ERROR defines error log level.
	ERROR = 3
	// CRITICAL defines critical log level.
	CRITICAL = 4
	// FATAL defines fatal log level.
	FATAL = 5
)

const (
	// TypeBool is index of boolean encoder in encoders tabel of a Writer.
	TypeBool = iota
	// TypeString is index of string encoder in encoders tabel of a Writer.
	TypeString
	// TypeByte is index of byte encoder in encoders tabel of a Writer.
	TypeByte
	// TypeInt is index of int encoder in encoders tabel of a Writer.
	TypeInt
	// TypeInt8 is index of int8 encoder in encoders tabel of a Writer.
	TypeInt8
	// TypeInt16 is index of int16 encoder in encoders tabel of a Writer.
	TypeInt16
	// TypeInt32 is index of int32 encoder in encoders tabel of a Writer.
	TypeInt32
	// TypeInt64 is index of int64 encoder in encoders tabel of a Writer.
	TypeInt64
	// TypeUint is index of uint encoder in encoders tabel of a Writer.
	TypeUint
	// TypeUint8 is index of uint8 encoder in encoders tabel of a Writer.
	TypeUint8
	// TypeUint16 is index of uint16 encoder in encoders tabel of a Writer.
	TypeUint16
	// TypeUint32 is index of uint32 encoder in encoders tabel of a Writer.
	TypeUint32
	// TypeUint64 is index of uint64 encoder in encoders tabel of a Writer.
	TypeUint64
	// TypeFloat32 is index of float32 encoder in encoders tabel of a Writer.
	TypeFloat32
	// TypeFloat64 is index of float64 encoder in encoders tabel of a Writer.
	TypeFloat64
	// TypeTime is index of time.Time encoder in encoders tabel of a Writer.
	TypeTime
	// TypeLevel is index of log.Entry.Lvl encoder in encoders tabel of a Writer.
	TypeLevel
	// TypeMessage is index of log.Entry.Msg encoder in encoders tabel of a Writer.
	TypeMessage
	// TypeContext is index of log.Entry.Ctx encoder in encoders tabel of a Writer.
	TypeContext
	// TypeLocation is index of log.Location encoder in encoders tabel of a Writer.
	TypeLocation
	// TypeGoroutineID is index of goroutine ID encoder in encoders tabel of a Writer.
	TypeGoroutineID
	// TypeFields is index of custom fields encoder in encoders tabel of a Writer.
	TypeFields
)

// Field defines additional custom field to be logged.
type Field struct {
	// Key is the field name, it can be represented in the output log entry.
	Key string
	// Val is the field value, it should be represented in the output log
	// entry.
	Val interface{}
	// Type is the value type. It will be used as encoder index in encoders
	// table of a Writer.
	Type int
}

// Encoder encodes a log information into buffer.
type Encoder interface {
	Encode(b []byte, k string, v interface{}) []byte
}

// Entry contains general log entry information.
type Entry struct {
	// Context passed to logger during creation via New or With functions.
	Ctx context.Context
	// Lvl represents the level of the log entry.
	Lvl int
	// Msg contains text message of the log entry. It does not mind format
	// verbs, instead custom fields should be used to provide more context
	// about the entry.
	Msg string
}

// Logger defines the logger interface.
// type Logger interface {
// 	// Debug logs entries with the debug level.
// 	Debug(msg string, fields ...Field)
// 	// Info logs entries with the info level.
// 	Info(msg string, fields ...Field)
// 	// Warning logs entries with the warning level.
// 	Warn(msg string, fields ...Field)
// 	// Error logs entries with the error level.
// 	Error(msg string, fields ...Field)
// 	// Critical logs entries with the critical level.
// 	Critical(msg string, fields ...Field)
// 	// Fatal logs entries with the fatal level.
// 	Fatal(msg string, fields ...Field)
// 	// Context returns the context passed to logger duiring creation.
// 	Context() context.Context
// 	// With clones the logger but replaces its context.
// 	With(ctx context.Context) Logger
// }

// Map represents set of fields. It's used in logging context fields.
type Map []Field

var maps = sync.Pool{
	New: func() interface{} {
		m := Map(make([]Field, 1<<5))
		return &m
	},
}

func (m *Map) Str(k string, v string) Map {
	return append(*m, Field{Key: k, Val: v, Type: TypeString})
}

func Fields() *Map {
	return maps.Get().(*Map)
}

// Constructor allows to perform delayed initialization. Its main goals is to
// provide a way to set default values in structs. This approach allows to
// provide more declarative API.
type Constructor interface {
	Construct()
}

// Writer writes log entry with the custom fields into buffer. We need to deal
// with a buffer instead of io.Writer to ensure we do not perform additional
// allocations.
type Writer interface {
	Write(w io.Writer, e Entry, fields []Field) error
}

// Config represent a logger configuration.
type Config struct {
	// Level is the logger log level. Entries with the level below than the
	// provided will not be logged.
	Level int
	// Stream is the output stream where entries will be logged. It should be
	// thread safe as the library does not care about stream concurrency. It
	// allows to utilize implementations best fitting a certain scenarios.
	Stream io.Writer
	// Writer serializes a log entry to the stream provided. The writer is
	// responsible for entry format (plain text, json, etc.).
	Writer Writer
	// Context can contain additional fields to be logged with each entry.
	Context context.Context
}

const fmterr = "cannot write log entry: %v \nentry=%v\nfields=%v"

// New creates a new Logger instance accoirding to the config provided.
// If the config.Writer supports the Constructor interface its method
// Construct will be called.
func New(config Config) *Logger {
	if ctor, ok := config.Writer.(Constructor); ok {
		ctor.Construct()
	}

	print := func(ctx context.Context, lvl int, msg string, fields []Field) {
		e := Entry{
			Ctx: ctx,
			Lvl: lvl,
			Msg: msg,
		}

		if err := config.Writer.Write(config.Stream, e, fields); err != nil {
			fmt.Fprintf(os.Stderr, fmterr, err, e, fields)
		}
	}

	discard := func(ctx context.Context, lvl int, msg string, fields []Field) {}

	writers := make([]writer, FATAL+1)
	for i := 0; i < len(writers); i++ {
		if i >= int(config.Level) {
			writers[i] = print
		} else {
			writers[i] = discard
		}
	}

	return &Logger{writers: writers, context: config.Context}
}

func (l *Logger) Debug(msg string, fields ...Field) {
	l.writers[DEBUG](l.context, DEBUG, msg, fields)
}

func (l *Logger) Info(msg string, fields ...Field) {
	l.writers[INFO](l.context, INFO, msg, fields)
}

func (l *Logger) Info2(msg string, fields Map) {
	l.writers[INFO](l.context, INFO, msg, fields)
}

func (l *Logger) Warn(msg string, fields ...Field) {
	l.writers[WARNING](l.context, WARNING, msg, fields)
}

func (l *Logger) Error(msg string, fields ...Field) {
	l.writers[ERROR](l.context, ERROR, msg, fields)
}

func (l *Logger) Critical(msg string, fields ...Field) {
	l.writers[CRITICAL](l.context, CRITICAL, msg, fields)
}

func (l *Logger) Fatal(msg string, fields ...Field) {
	l.writers[FATAL](l.context, FATAL, msg, fields)
}

func (l *Logger) Context() context.Context {
	return l.context
}

func (l *Logger) With(ctx context.Context) *Logger {
	return &Logger{writers: l.writers, context: ctx}
}

type writer func(ctx context.Context, lvl int, msg string, fields []Field)

type Logger struct {
	writers []writer
	context context.Context
}

type ck int

const ckey = ck(0)

// With returns context with the field added.
func With(c context.Context, f Field) context.Context {
	fields := c.Value(ckey)
	if fields == nil {
		fields = []Field{}
	}

	fields = append(fields.([]Field), f)
	return context.WithValue(c, ckey, fields)
}
