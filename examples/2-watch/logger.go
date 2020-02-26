package main

import "log"

const (
	// None level
	None uint = iota
	// Error level
	Error
	// Warn level
	Warn
	// Info level
	Info
	// Debug level
	Debug
)

// Logger is a simple logger with four different levels: Debug, Info, Warn, Error
type Logger struct {
	Level uint
}

// SetLevel changes the verbosity level
func (l *Logger) SetLevel(level string) {
	switch level {
	case "None":
		l.Level = None
	case "Error":
		l.Level = Error
	case "Warn":
		l.Level = Warn
	case "Info":
		l.Level = Info
	case "Debug":
		l.Level = Debug
	}
}

// Debug logs in debug level (1)
func (l *Logger) Debug(msg string, args ...interface{}) {
	if l.Level >= Debug {
		log.Printf(msg, args...)
	}
}

// Info logs in info level (2)
func (l *Logger) Info(msg string, args ...interface{}) {
	if l.Level >= Info {
		log.Printf(msg, args...)
	}
}

// Warn logs in warn level (3)
func (l *Logger) Warn(msg string, args ...interface{}) {
	if l.Level >= Warn {
		log.Printf(msg, args...)
	}
}

// Error logs in error level (4)
func (l *Logger) Error(msg string, args ...interface{}) {
	if l.Level >= Error {
		log.Printf(msg, args...)
	}
}
