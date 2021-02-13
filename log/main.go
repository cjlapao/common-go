package log

import (
	"fmt"
	"os"

	strcolor "github.com/cjlapao/common-go/strcolor"
)

// Log Interface
type Log interface {
	Log(format string, level Level, words ...string)
	LogHighlight(format string, level Level, highlightColor strcolor.ColorCode, words ...string)
	Info(format string, words ...string)
	Success(format string, words ...string)
	TaskSuccess(format string, isComplete bool, words ...string)
	Warn(format string, words ...string)
	TaskWarn(format string, words ...string)
	Command(format string, words ...string)
	Disabled(format string, words ...string)
	Notice(format string, words ...string)
	Debug(format string, words ...string)
	Trace(format string, words ...string)
	Error(format string, words ...string)
	LogError(message error)
	TaskError(format string, isComplete bool, words ...string)
	Fatal(format string, words ...string)
	FatalError(e error, format string, words ...string)
}

// Logger Default structure
type Logger struct {
	Loggers       []Log
	LogLevel      Level
	HighlighColor strcolor.ColorCode
}

var globalLogger *Logger

// Level Entity
type Level int

// LogLevel Enum Definition
const (
	Error Level = iota
	Warning
	Info
	Debug
	Trace
)

// Get Creates a new Logger instance
func Get() *Logger {
	if globalLogger == nil {
		result := Logger{
			LogLevel:      Info,
			HighlighColor: strcolor.BrightYellow,
		}
		result.Loggers = []Log{}
		result.AddCmdLogger()

		debug := os.Getenv("DT_DEBUG")
		if debug == "true" {
			result.LogLevel = Debug
		}

		trace := os.Getenv("DT_TRACE")
		if trace == "trace" {
			result.LogLevel = Trace
		}

		globalLogger = &result
		return &result
	}

	return globalLogger
}

// AddCmdLogger Add a command line logger to the system
func (l *Logger) AddCmdLogger() {
	found := false
	for _, logger := range l.Loggers {
		xType := fmt.Sprintf("%T", logger)
		if xType == "CmdLogger" {
			found = true
			break
		}
	}

	if !found {
		l.Loggers = append(l.Loggers, new(CmdLogger))
	}
}

// Log Log information message
func (l *Logger) Log(format string, level Level, words ...string) {
	for _, logger := range l.Loggers {
		logger.Log(format, level, words...)
	}
}

// LogHighlight Log information message
func (l *Logger) LogHighlight(format string, level Level, words ...string) {
	for _, logger := range l.Loggers {
		logger.LogHighlight(format, level, l.HighlighColor, words...)
	}
}

// Info log information message
func (l *Logger) Info(format string, words ...string) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Info(format, words...)
		}
	}
}

// Success log message
func (l *Logger) Success(format string, words ...string) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Success(format, words...)
		}
	}
}

// TaskSuccess log message
func (l *Logger) TaskSuccess(format string, isComplete bool, words ...string) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.TaskSuccess(format, isComplete, words...)
		}
	}
}

// Warn log message
func (l *Logger) Warn(format string, words ...string) {
	if l.LogLevel >= Warning {
		for _, logger := range l.Loggers {
			logger.Warn(format, words...)
		}
	}
}

// TaskWarn log message
func (l *Logger) TaskWarn(format string, words ...string) {
	if l.LogLevel >= Warning {
		for _, logger := range l.Loggers {
			logger.TaskWarn(format, words...)
		}
	}
}

// Command log message
func (l *Logger) Command(format string, words ...string) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Command(format, words...)
		}
	}
}

// Disabled log message
func (l *Logger) Disabled(format string, words ...string) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Disabled(format, words...)
		}
	}
}

// Notice log message
func (l *Logger) Notice(format string, words ...string) {
	if l.LogLevel >= Info {
		for _, logger := range l.Loggers {
			logger.Notice(format, words...)
		}
	}
}

// Debug log message
func (l *Logger) Debug(format string, words ...string) {
	if l.LogLevel >= Debug {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Trace log message
func (l *Logger) Trace(format string, words ...string) {
	if l.LogLevel >= Trace {
		for _, logger := range l.Loggers {
			logger.Debug(format, words...)
		}
	}
}

// Error log message
func (l *Logger) Error(format string, words ...string) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Error(format, words...)
		}
	}
}

// LogError log message
func (l *Logger) LogError(message error) {
	if l.LogLevel >= Error {
		if message != nil {
			for _, logger := range l.Loggers {
				logger.LogError(message)
			}
		}
	}
}

// TaskError log message
func (l *Logger) TaskError(format string, isComplete bool, words ...string) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.TaskError(format, isComplete, words...)
		}
	}
}

// Fatal log message
func (l *Logger) Fatal(format string, words ...string) {
	if l.LogLevel >= Error {
		for _, logger := range l.Loggers {
			logger.Fatal(format, words...)
		}
	}
}

// FatalError log message
func (l *Logger) FatalError(e error, format string, words ...string) {
	for _, logger := range l.Loggers {
		logger.Error(format, words...)
	}

	if e != nil {
		panic(e)
	}
}
