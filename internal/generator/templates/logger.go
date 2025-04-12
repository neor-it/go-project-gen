// internal/generator/templates/logger.go - Templates for logger files
package templates

// LoggerTemplate returns the content of the logger.go file
func LoggerTemplate() string {
	return `// internal/logger/logger.go - Logger implementation
package logger

import (
	"os"
	"strings"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger interface defines the methods that the logger should implement
type Logger interface {
	Debug(msg string, keysAndValues ...interface{})
	Info(msg string, keysAndValues ...interface{})
	Warn(msg string, keysAndValues ...interface{})
	Error(msg string, keysAndValues ...interface{})
	Fatal(msg string, keysAndValues ...interface{})
	SetLevel(level string)
}

// ZapLogger implements the Logger interface using Zap
type ZapLogger struct {
	logger *zap.SugaredLogger
	level  zapcore.Level
	core   zapcore.Core
	atom   zap.AtomicLevel
}

// NewLogger creates a new logger
func NewLogger() Logger {
	// Default to info level
	defaultLevel := getLogLevelFromEnv()
	atom := zap.NewAtomicLevelAt(defaultLevel)

	// Create encoder configuration
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalColorLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// Create core
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		zapcore.AddSync(os.Stdout),
		atom,
	)

	// Create logger
	zapLogger := zap.New(core)
	defer zapLogger.Sync()

	// Return ZapLogger
	return &ZapLogger{
		logger: zapLogger.Sugar(),
		level:  defaultLevel,
		core:   core,
		atom:   atom,
	}
}

// Debug logs a debug message
func (l *ZapLogger) Debug(msg string, keysAndValues ...interface{}) {
	l.logger.Debugw(msg, keysAndValues...)
}

// Info logs an info message
func (l *ZapLogger) Info(msg string, keysAndValues ...interface{}) {
	l.logger.Infow(msg, keysAndValues...)
}

// Warn logs a warning message
func (l *ZapLogger) Warn(msg string, keysAndValues ...interface{}) {
	l.logger.Warnw(msg, keysAndValues...)
}

// Error logs an error message
func (l *ZapLogger) Error(msg string, keysAndValues ...interface{}) {
	l.logger.Errorw(msg, keysAndValues...)
}

// Fatal logs a fatal message and exits
func (l *ZapLogger) Fatal(msg string, keysAndValues ...interface{}) {
	l.logger.Fatalw(msg, keysAndValues...)
}

// SetLevel sets the logger level
func (l *ZapLogger) SetLevel(level string) {
	newLevel := parseLogLevel(level)
	l.atom.SetLevel(newLevel)
	l.level = newLevel
}

// getLogLevelFromEnv gets the log level from environment variable or returns the default
func getLogLevelFromEnv() zapcore.Level {
	levelStr := os.Getenv("LOGGING_LEVEL")
	if levelStr == "" {
		return zapcore.InfoLevel
	}
	return parseLogLevel(levelStr)
}

// parseLogLevel parses a string log level to a zapcore.Level
func parseLogLevel(level string) zapcore.Level {
	switch strings.ToLower(level) {
	case "debug":
		return zapcore.DebugLevel
	case "info":
		return zapcore.InfoLevel
	case "warn", "warning":
		return zapcore.WarnLevel
	case "error":
		return zapcore.ErrorLevel
	case "fatal":
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}
`
}
