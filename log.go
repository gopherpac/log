package log

import (
	"context"
	"io"
	"os"
	"unsafe"

	"cdr.dev/slog"
	"cdr.dev/slog/sloggers/sloghuman"
	"cdr.dev/slog/sloggers/slogjson"
)

// Level defined log level
type Level slog.Level

// Field contains additional logging fields
type Field slog.Field

// JSON creates a json logger
func JSON(w io.Writer) slog.Logger {
	return slogjson.Make(w)
}

// Terminal creates a terminal logger
func Terminal(w io.Writer) slog.Logger {
	return sloghuman.Make(w)
}

var (
	// LevelDebug is used for development and debugging messages.
	LevelDebug = Level(slog.LevelDebug)

	// LevelInfo is used for normal informational messages.
	LevelInfo = Level(slog.LevelInfo)

	// LevelWarn is used when something has possibly gone wrong.
	LevelWarn = Level(slog.LevelWarn)

	// LevelError is used when something has certainly gone wrong.
	LevelError = Level(slog.LevelError)

	// LevelCritical is used when when something has gone wrong and should
	// be immediately investigated.
	LevelCritical = Level(slog.LevelCritical)

	// LevelFatal is used when the process is about to exit due to an error.
	LevelFatal = Level(slog.LevelFatal)
)

// Logger defines a logger interface
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Critical(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)
}

type wrapper struct {
	logger  *slog.Logger
	context context.Context
}

// Config represents a logger configuration
type Config struct {
	Writer  io.Writer
	Level   Level
	Context context.Context
	Make    func(io.Writer) slog.Logger
	Fields  []Field
}

// New creates a new logger
func New(config Config) Logger {
	logger := config.
		Make(config.Writer).
		Leveled(slog.Level(config.Level)).
		With((*(*[]slog.Field)(unsafe.Pointer(&config.Fields)))...)
	return &wrapper{
		logger:  &logger,
		context: config.Context,
	}
}

func (w *wrapper) Debug(msg string, fields ...Field) {
	w.logger.Debug(w.context, msg, (*(*[]slog.Field)(unsafe.Pointer(&fields)))...)
}

func (w *wrapper) Info(msg string, fields ...Field) {
	w.logger.Info(w.context, msg, (*(*[]slog.Field)(unsafe.Pointer(&fields)))...)
}

func (w *wrapper) Warn(msg string, fields ...Field) {
	w.logger.Warn(w.context, msg, (*(*[]slog.Field)(unsafe.Pointer(&fields)))...)
}

func (w *wrapper) Error(msg string, fields ...Field) {
	w.logger.Error(w.context, msg, (*(*[]slog.Field)(unsafe.Pointer(&fields)))...)
}

func (w *wrapper) Critical(msg string, fields ...Field) {
	w.logger.Critical(w.context, msg, (*(*[]slog.Field)(unsafe.Pointer(&fields)))...)
}

func (w *wrapper) Fatal(msg string, fields ...Field) {
	w.logger.Fatal(w.context, msg, (*(*[]slog.Field)(unsafe.Pointer(&fields)))...)
}

// DefaultConfig represents default logger configuration
var DefaultConfig = Config{
	Writer:  os.Stdout,
	Make:    Terminal,
	Level:   LevelDebug,
	Context: context.Background(),
	Fields:  []Field{},
}

// Default is the default logger
var Default = New(DefaultConfig)
