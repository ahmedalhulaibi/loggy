package loggy

import (
	"context"

	"go.uber.org/zap"
)

// Logger is an extension of a zap.s
// It is configured with a list of fields
// Configured fields are context keys (as string) to extract request-scoped values from context.Context
type Logger struct {
	s *zap.SugaredLogger
}

func New(zapLogger *zap.SugaredLogger) Logger {
	return Logger{
		s: zapLogger,
	}
}

// With creates a child logger, and optionally adds some context to that logger.
// The child logger inherits the context of its parent.
func (l Logger) With(ctx context.Context, args ...interface{}) (context.Context, Logger) {
	l = l.extractLogger(ctx)
	newLogger := New(l.s.With(args...))
	return context.WithValue(ctx, loggerctxkey, newLogger), newLogger
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Debug(ctx context.Context, args ...interface{}) {
	l.extractLogger(ctx).s.Debug(args...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Info(ctx context.Context, args ...interface{}) {
	l.extractLogger(ctx).s.Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Warn(ctx context.Context, args ...interface{}) {
	l.extractLogger(ctx).s.Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Error(ctx context.Context, args ...interface{}) {
	l.extractLogger(ctx).s.Error(args...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l Logger) DPanic(ctx context.Context, args ...interface{}) {
	l.extractLogger(ctx).s.DPanic(args...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l Logger) Panic(ctx context.Context, args ...interface{}) {
	l.extractLogger(ctx).s.Panic(args...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l Logger) Fatal(ctx context.Context, args ...interface{}) {
	l.extractLogger(ctx).s.Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l Logger) Debugf(ctx context.Context, template string, args ...interface{}) {
	l.extractLogger(ctx).s.Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l Logger) Infof(ctx context.Context, template string, args ...interface{}) {
	l.extractLogger(ctx).s.Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l Logger) Warnf(ctx context.Context, template string, args ...interface{}) {
	l.extractLogger(ctx).s.Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l Logger) Errorf(ctx context.Context, template string, args ...interface{}) {
	l.extractLogger(ctx).s.Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See zapcore.DPanicLevel for details.)
func (l Logger) DPanicf(ctx context.Context, template string, args ...interface{}) {
	l.extractLogger(ctx).s.DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l Logger) Panicf(ctx context.Context, template string, args ...interface{}) {
	l.extractLogger(ctx).s.Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l Logger) Fatalf(ctx context.Context, template string, args ...interface{}) {
	l.extractLogger(ctx).s.Fatalf(template, args...)
}

// Debugw logs a message with some additional context.
func (l Logger) Debugw(ctx context.Context, msg string, args ...interface{}) {
	l.extractLogger(ctx).s.Debugw(msg, args...)
}

// Infow logs a message with some additional context.
func (l Logger) Infow(ctx context.Context, msg string, args ...interface{}) {
	l.extractLogger(ctx).s.Infow(msg, args...)
}

// Warnw logs a message with some additional context.
func (l Logger) Warnw(ctx context.Context, msg string, args ...interface{}) {
	l.extractLogger(ctx).s.Warnw(msg, args...)
}

// Errorw logs a message with some additional context.
func (l Logger) Errorw(ctx context.Context, msg string, args ...interface{}) {
	l.extractLogger(ctx).s.Errorw(msg, args...)
}

// DPanicw logs a message with some additional context. In development, the logger then panics. (See zapcore.DPanicLevel for details.)
func (l Logger) DPanicw(ctx context.Context, msg string, args ...interface{}) {
	l.extractLogger(ctx).s.DPanicw(msg, args...)
}

// Panicw logs a message with some additional context, then panics.
func (l Logger) Panicw(ctx context.Context, msg string, args ...interface{}) {
	l.extractLogger(ctx).s.Panicw(msg, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit.
func (l Logger) Fatalw(ctx context.Context, msg string, args ...interface{}) {
	l.extractLogger(ctx).s.Fatalw(msg, args...)
}

type logContextKey string

const (
	loggerctxkey = logContextKey("logger")
)

func (l Logger) extractLogger(ctx context.Context) Logger {
	logger, ok := ctx.Value(loggerctxkey).(Logger)
	if !ok {
		return l
	}
	return logger
}
