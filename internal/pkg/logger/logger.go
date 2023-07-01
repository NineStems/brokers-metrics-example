package logger

import (
	"go.uber.org/zap"
)

type Logger struct {
	Log *zap.SugaredLogger
}

func InitSugarZapLogger(log *zap.Logger) *Logger {
	return &Logger{
		Log: log.Sugar(),
	}
}

// Debug uses fmt.Sprint to construct and log a message.
func (l *Logger) Debug(args ...interface{}) {
	l.Log.Debug(args...)
}

// Debugf uses fmt.Sprintf to log a templated message.
func (l *Logger) Debugf(template string, args ...interface{}) {
	l.Log.Debugf(template, args...)
}

// Debugw logs a message with some additional context.
func (l *Logger) Debugw(msg string, keysAndValues ...interface{}) {
	l.Log.Debugw(msg, keysAndValues...)
}

// Info uses fmt.Sprint to construct and log a message.
func (l *Logger) Info(args ...interface{}) {
	l.Log.Info(args...)
}

// Infof uses fmt.Sprintf to log a templated message.
func (l *Logger) Infof(template string, args ...interface{}) {
	l.Log.Infof(template, args...)
}

// Infow logs a message with some additional context.
func (l *Logger) Infow(msg string, keysAndValues ...interface{}) {
	l.Log.Infow(msg, keysAndValues...)
}

// Warn uses fmt.Sprint to construct and log a message.
func (l *Logger) Warn(args ...interface{}) {
	l.Log.Warn(args...)
}

// Warnf uses fmt.Sprintf to log a templated message.
func (l *Logger) Warnf(template string, args ...interface{}) {
	l.Log.Warnf(template, args...)
}

// Warnw logs a message with some additional context.
func (l *Logger) Warnw(msg string, keysAndValues ...interface{}) {
	l.Log.Warnw(msg, keysAndValues...)
}

// Error uses fmt.Sprint to construct and log a message.
func (l *Logger) Error(args ...interface{}) {
	l.Log.Error(args...)
}

// Errorf uses fmt.Sprintf to log a templated message.
func (l *Logger) Errorf(template string, args ...interface{}) {
	l.Log.Errorf(template, args...)
}

// Errorw logs a message with some additional context.
func (l *Logger) Errorw(msg string, keysAndValues ...interface{}) {
	l.Log.Errorw(msg, keysAndValues...)
}

// Fatal uses fmt.Sprint to construct and log a message, then calls os.Exit.
func (l *Logger) Fatal(args ...interface{}) {
	l.Log.Fatal(args...)
}

// Fatalf uses fmt.Sprintf to log a templated message, then calls os.Exit.
func (l *Logger) Fatalf(template string, args ...interface{}) {
	l.Log.Fatalf(template, args...)
}

// Fatalw logs a message with some additional context, then calls os.Exit.
func (l *Logger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.Log.Fatalw(msg, keysAndValues...)
}

// Panic uses fmt.Sprint to construct and log a message, then panics.
func (l *Logger) Panic(args ...interface{}) {
	l.Log.Panic(args...)
}

// Panicf uses fmt.Sprintf to log a templated message, then panics.
func (l *Logger) Panicf(template string, args ...interface{}) {
	l.Log.Panicf(template, args...)
}

// Panicw logs a message with some additional context, then panics.
func (l *Logger) Panicw(msg string, keysAndValues ...interface{}) {
	l.Log.Panicw(msg, keysAndValues...)
}

// Print logs a message at level Debug on the ZapLogger.
func (l *Logger) Print(args ...interface{}) {
	l.Log.Debug(args...)
}

// Printf logs a message at level Debug on the ZapLogger.
func (l *Logger) Printf(template string, args ...interface{}) {
	l.Log.Debugf(template, args...)
}

// With return a log with an extra field.
func (l *Logger) With(key string, value interface{}) *Logger {
	return &Logger{l.Log.With(zap.Any(key, value))}
}

// WithField return a log with an extra field.
func (l *Logger) WithField(key string, value interface{}) *Logger {
	return &Logger{l.Log.With(zap.Any(key, value))}
}

// WithFields return a log with extra fields.
func (l *Logger) WithFields(fields map[string]interface{}) *Logger {
	i := 0
	var clog *Logger
	for k, v := range fields {
		if i == 0 {
			clog = l.WithField(k, v)
		} else if clog != nil {
			clog = clog.WithField(k, v)
		}
		i++
	}
	return clog
}
