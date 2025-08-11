package utils

import (
	"fmt"
	"log"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

type LogLevel string

const (
	LogLevelInfo  LogLevel = "INFO"
	LogLevelWarn  LogLevel = "WARN"
	LogLevelError LogLevel = "ERROR"
	LogLevelDebug LogLevel = "DEBUG"
	LogLevelFatal LogLevel = "FATAL"
)

type LogEntry struct {
	Timestamp string   `json:"timestamp"`
	TraceID   string   `json:"trace_id"`
	Level     LogLevel `json:"level"`
	Message   string   `json:"message"`
	Method    string   `json:"method,omitempty"`
	Path      string   `json:"path,omitempty"`
	Status    int      `json:"status,omitempty"`
	Error     string   `json:"error,omitempty"`
	File      string   `json:"file,omitempty"`
	Line      int      `json:"line,omitempty"`
	UserAgent string   `json:"user_agent,omitempty"`
	IP        string   `json:"ip,omitempty"`
}

type Logger struct {
	traceID string
	context *gin.Context
}

func NewLogger(c *gin.Context) *Logger {
	traceID := GetTraceID(c)
	return &Logger{
		traceID: traceID,
		context: c,
	}
}

func getCallerInfo() (string, int) {
	_, file, line, ok := runtime.Caller(3)
	if !ok {
		return "unknown", 0
	}

	for i := len(file) - 1; i > 0; i-- {
		if file[i] == '/' || file[i] == '\\' {
			file = file[i+1:]
			break
		}
	}
	return file, line
}

func (l *Logger) formatLog(level LogLevel, message string, err error) LogEntry {
	file, line := getCallerInfo()

	entry := LogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05.000"),
		TraceID:   l.traceID,
		Level:     level,
		Message:   message,
		File:      file,
		Line:      line,
	}

	if l.context != nil {
		entry.Method = l.context.Request.Method
		entry.Path = l.context.Request.URL.Path
		entry.Status = l.context.Writer.Status()
		entry.UserAgent = l.context.Request.UserAgent()
		entry.IP = l.context.ClientIP()
	}

	if err != nil {
		entry.Error = err.Error()
	}

	return entry
}

func (l *Logger) logEntry(entry LogEntry) {
	if entry.Error != "" {
		log.Printf("[%s - TRACE: %s] %s %s:%d - %s | Error: %s | %s %s (Status: %d) | IP: %s",
			entry.Level, entry.TraceID, entry.File, "", entry.Line, entry.Message,
			entry.Error, entry.Method, entry.Path, entry.Status, entry.IP)
	} else {
		log.Printf("[%s - TRACE: %s] %s %s:%d - %s | %s %s (Status: %d) | IP: %s",
			entry.Level, entry.TraceID, entry.File, "", entry.Line, entry.Message,
			entry.Method, entry.Path, entry.Status, entry.IP)
	}
}

func (l *Logger) Info(message string) {
	entry := l.formatLog(LogLevelInfo, message, nil)
	l.logEntry(entry)
}

func (l *Logger) Warn(message string) {
	entry := l.formatLog(LogLevelWarn, message, nil)
	l.logEntry(entry)
}

func (l *Logger) Error(message string, err error) {
	entry := l.formatLog(LogLevelError, message, err)
	l.logEntry(entry)
}

func (l *Logger) Debug(message string) {
	entry := l.formatLog(LogLevelDebug, message, nil)
	l.logEntry(entry)
}

func (l *Logger) Fatal(message string, err error) {
	entry := l.formatLog(LogLevelFatal, message, err)
	l.logEntry(entry)
	log.Fatal(message)
}

func (l *Logger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	entry := l.formatLog(LogLevelError, message, nil)
	l.logEntry(entry)
}

func (l *Logger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	entry := l.formatLog(LogLevelInfo, message, nil)
	l.logEntry(entry)
}

func (l *Logger) DatabaseError(operation string, err error) {
	message := fmt.Sprintf("Database operation failed: %s", operation)
	entry := l.formatLog(LogLevelError, message, err)
	l.logEntry(entry)
}

func (l *Logger) ValidationError(field string, value interface{}, reason string) {
	message := fmt.Sprintf("Validation failed for field '%s' with value '%v': %s", field, value, reason)
	entry := l.formatLog(LogLevelError, message, nil)
	l.logEntry(entry)
}

func (l *Logger) FileUploadError(filename string, err error) {
	message := fmt.Sprintf("File upload failed for '%s'", filename)
	entry := l.formatLog(LogLevelError, message, err)
	l.logEntry(entry)
}

func (l *Logger) AuthError(action string, reason string) {
	message := fmt.Sprintf("Authentication/Authorization failed for action '%s': %s", action, reason)
	entry := l.formatLog(LogLevelError, message, nil)
	l.logEntry(entry)
}

func LogInfo(c *gin.Context, message string) {
	logger := NewLogger(c)
	logger.Info(message)
}

func LogError(c *gin.Context, message string, err error) {
	logger := NewLogger(c)
	logger.Error(message, err)
}

func LogWarn(c *gin.Context, message string) {
	logger := NewLogger(c)
	logger.Warn(message)
}

func LogDebug(c *gin.Context, message string) {
	logger := NewLogger(c)
	logger.Debug(message)
}
