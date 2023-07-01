package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	ljack "gopkg.in/natefinch/lumberjack.v2"
)

func Console(fileName string, level string) *zap.Logger {
	core := newCore(fileName, level)
	return zap.New(core)
}

func newCore(fileName string, level string) zapcore.Core {
	var low zapcore.Level

	switch level {
	case "debug":
		low = zap.DebugLevel
	case "warn":
		low = zap.WarnLevel
	case "info":
		low = zap.InfoLevel
	case "error":
		low = zap.ErrorLevel
	default:
		low = zap.InfoLevel
	}

	priority := zap.LevelEnablerFunc(func(lvl zapcore.Level) bool {
		return low <= lvl
	})

	consoleDebugging := zapcore.Lock(os.Stdout)
	fileWriter := zapcore.AddSync(&ljack.Logger{
		Filename:   fileName,
		MaxSize:    50, // megabytes
		MaxBackups: 3,
		MaxAge:     28, // days
	})

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	jsonEncoder := zapcore.NewJSONEncoder(encoderConfig)
	consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)

	file := zapcore.NewCore(jsonEncoder, fileWriter, priority)
	stdout := zapcore.NewCore(consoleEncoder, consoleDebugging, priority)

	return zapcore.NewTee(file, stdout)
}
