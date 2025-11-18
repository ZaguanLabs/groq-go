package groq

import (
	"fmt"
	"os"
)

// Logger is the interface for logging
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// LeveledLogger is a simple logger implementation
type LeveledLogger struct {
	Level Level
}

type Level int

const (
	LevelDebug Level = iota
	LevelInfo
	LevelWarn
	LevelError
	LevelNone
)

func (l *LeveledLogger) Debug(msg string, args ...interface{}) {
	if l.Level <= LevelDebug {
		fmt.Fprintf(os.Stderr, "[DEBUG] "+msg+"\n", args...)
	}
}

func (l *LeveledLogger) Info(msg string, args ...interface{}) {
	if l.Level <= LevelInfo {
		fmt.Fprintf(os.Stderr, "[INFO] "+msg+"\n", args...)
	}
}

func (l *LeveledLogger) Warn(msg string, args ...interface{}) {
	if l.Level <= LevelWarn {
		fmt.Fprintf(os.Stderr, "[WARN] "+msg+"\n", args...)
	}
}

func (l *LeveledLogger) Error(msg string, args ...interface{}) {
	if l.Level <= LevelError {
		fmt.Fprintf(os.Stderr, "[ERROR] "+msg+"\n", args...)
	}
}

var defaultLogger Logger = &LeveledLogger{Level: LevelInfo}

func init() {
	if os.Getenv("GROQ_LOG") == "debug" {
		defaultLogger = &LeveledLogger{Level: LevelDebug}
	} else if os.Getenv("GROQ_LOG") == "info" {
		defaultLogger = &LeveledLogger{Level: LevelInfo}
	} else {
		// Default to info or warn? The plan says env toggle.
		// Usually SDKs are silent by default unless configured.
		defaultLogger = &LeveledLogger{Level: LevelWarn}
	}
}
