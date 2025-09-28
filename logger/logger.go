package logger

import (
	"io"

	"github.com/sirupsen/logrus"
)

// Logger interface defines the logging methods
type Logger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
}

// LogrusLogger implements the Logger interface using logrus
type LogrusLogger struct {
	*logrus.Logger
}

// New creates a new logger instance
func New(level string) Logger {
	logger := logrus.New()

	// Set log level
	switch level {
	case "debug":
		logger.SetLevel(logrus.DebugLevel)
	case "info":
		logger.SetLevel(logrus.InfoLevel)
	case "warn":
		logger.SetLevel(logrus.WarnLevel)
	case "error":
		logger.SetLevel(logrus.ErrorLevel)
	default:
		logger.SetLevel(logrus.InfoLevel)
	}

	// Set JSON formatter for structured logging
	logger.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
	})

	return &LogrusLogger{Logger: logger}
}

// SetOutput sets the output destination for the logger
func (l *LogrusLogger) SetOutput(output io.Writer) {
	l.Logger.SetOutput(output)
}
