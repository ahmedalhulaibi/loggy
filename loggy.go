package loggy

import (
	"context"

	"go.uber.org/zap"
)

type LoggyCtxKey = string

// Logger is an extension of a zap.SugaredLogger
// It is configured with a list of fields
// Configured fields are context keys (as string) to extract request-scoped values from context.Context
type Logger struct {
	*zap.SugaredLogger
	fields []LoggyCtxKey
}

func New(zapLogger *zap.SugaredLogger) Logger {
	return Logger{
		SugaredLogger: zapLogger,
	}
}

// WithFields returns a new logger which will inspect for additional fields
// Configured fields are context keys (as string) to extract request-scoped values from context.Context
func (l Logger) WithFields(fields ...LoggyCtxKey) Logger {
	l.fields = append(l.fields, fields...)
	return l
}

// Debug logs a message at DebugLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Debug(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Debug(args...)
}

// Info logs a message at InfoLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Info(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Info(args...)
}

// Warn uses fmt.Sprint to construct and log a message.
// Warn logs a message at WarnLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Warn(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Warn(args...)
}

// Error uses fmt.Sprint to construct and log a message.
// Error logs a message at ErrorLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
func (l Logger) Error(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Error(args...)
}

// DPanic logs a message at DPanicLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
//
// If the logger is in development mode, it then panics (DPanic means
// "development panic"). This is useful for catching errors that are
// recoverable, but shouldn't ever happen.
func (l Logger) DPanic(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).DPanic(args...)
}

// Panic logs a message at PanicLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
//
// The logger then panics, even if logging at PanicLevel is disabled.
func (l Logger) Panic(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).DPanic(args...)
}

// Fatal logs a message at FatalLevel. The message includes any fields passed
// at the log site, as well as any fields extracted from the context.
//
// The logger then calls os.Exit(1), even if logging at FatalLevel is
// disabled.
func (l Logger) Fatal(ctx context.Context, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Fatal(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l Logger) Debugf(ctx context.Context, template string, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Debugf(template, args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l Logger) Infof(ctx context.Context, template string, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Infof(template, args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l Logger) Warnf(ctx context.Context, template string, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Warnf(template, args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l Logger) Errorf(ctx context.Context, template string, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Errorf(template, args...)
}

// DPanicf uses fmt.Sprintf to log a templated message. In development, the logger then panics. (See zapcore.DPanicLevel for details.)
func (l Logger) DPanicf(ctx context.Context, template string, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).DPanicf(template, args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l Logger) Panicf(ctx context.Context, template string, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Panicf(template, args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l Logger) Fatalf(ctx context.Context, template string, args ...interface{}) {
	l.SugaredLogger.With(l.extractArgsFromCtx(ctx)...).Fatalf(template, args...)
}

// Debugw logs a message with some additional context.
func (l Logger) Debugw(ctx context.Context, msg string, args ...KeyVal) {
	finalArgs := append(l.extractArgsFromCtx(ctx), keyvals(args).toGenericSlice()...)
	l.SugaredLogger.Debugw(msg, finalArgs...)
}

// Infow logs a message with some additional context.
func (l Logger) Infow(ctx context.Context, msg string, args ...KeyVal) {
	finalArgs := append(l.extractArgsFromCtx(ctx), keyvals(args).toGenericSlice()...)
	l.SugaredLogger.Infow(msg, finalArgs...)
}

// Warnw logs a message with some additional context.
func (l Logger) Warnw(ctx context.Context, msg string, args ...KeyVal) {
	finalArgs := append(l.extractArgsFromCtx(ctx), keyvals(args).toGenericSlice()...)
	l.SugaredLogger.Warnw(msg, finalArgs...)
}

// Errorw logs a message with some additional context.
func (l Logger) Errorw(ctx context.Context, msg string, args ...KeyVal) {
	finalArgs := append(l.extractArgsFromCtx(ctx), keyvals(args).toGenericSlice()...)
	l.SugaredLogger.Errorw(msg, finalArgs...)
}

// DPanicw logs a message with some additional context. In development, the logger then panics. (See zapcore.DPanicLevel for details.)
func (l Logger) DPanicw(ctx context.Context, msg string, args ...KeyVal) {
	finalArgs := append(l.extractArgsFromCtx(ctx), keyvals(args).toGenericSlice()...)
	l.SugaredLogger.DPanicw(msg, finalArgs...)
}

// Panicw logs a message with some additional context, then panics.
func (l Logger) Panicw(ctx context.Context, msg string, args ...KeyVal) {
	finalArgs := append(l.extractArgsFromCtx(ctx), keyvals(args).toGenericSlice()...)
	l.SugaredLogger.Panicw(msg, finalArgs...)
}

// Fatalw logs a message with some additional context, then calls os.Exit.
func (l Logger) Fatalw(ctx context.Context, msg string, args ...KeyVal) {
	finalArgs := append(l.extractArgsFromCtx(ctx), keyvals(args).toGenericSlice()...)
	l.SugaredLogger.Fatalw(msg, finalArgs...)
}

// KeyVal represents a key-value pair for a single log attribute
type KeyVal struct {
	Key string
	Val interface{}
}

type keyvals []KeyVal

func (kvs keyvals) toGenericSlice() []interface{} {
	argsI := make([]interface{}, 0, len(kvs)*2)

	for _, kv := range kvs {
		argsI = append(argsI, kv.Key, kv.Val)
	}

	return argsI
}

func (l Logger) extractArgsFromCtx(ctx context.Context) []interface{} {
	ctxArgs := make([]interface{}, 0, len(l.fields)*2)

	for _, field := range l.fields {
		val := ctx.Value(field)
		if val != nil {
			ctxArgs = append(ctxArgs, field, val)
		}
	}

	return ctxArgs
}
