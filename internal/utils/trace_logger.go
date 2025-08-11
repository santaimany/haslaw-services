package utils

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"haslaw-be-services/internal/middleware"
)

// Logger levels
const (
	InfoLevel  = "INFO"
	WarnLevel  = "WARN"
	ErrorLevel = "ERROR"
	DebugLevel = "DEBUG"
)

// TraceLogger provides logging with trace ID context
type TraceLogger struct {
	traceID string
}

// NewTraceLogger creates a new logger with trace ID
func NewTraceLogger(c *gin.Context) *TraceLogger {
	traceID := middleware.GetTraceID(c)
	return &TraceLogger{traceID: traceID}
}

// NewTraceLoggerWithID creates a new logger with specific trace ID
func NewTraceLoggerWithID(traceID string) *TraceLogger {
	return &TraceLogger{traceID: traceID}
}

// formatMessage formats log message with trace ID and level
func (tl *TraceLogger) formatMessage(level, message string) string {
	timestamp := time.Now().Format("2006-01-02 15:04:05.000")
	return fmt.Sprintf("[%s] [TRACE: %s] [%s] %s", timestamp, tl.traceID, level, message)
}

// Info logs info level message
func (tl *TraceLogger) Info(message string) {
	log.Println(tl.formatMessage(InfoLevel, message))
}

// Infof logs formatted info level message
func (tl *TraceLogger) Infof(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	tl.Info(message)
}

// Warn logs warning level message
func (tl *TraceLogger) Warn(message string) {
	log.Println(tl.formatMessage(WarnLevel, message))
}

// Warnf logs formatted warning level message
func (tl *TraceLogger) Warnf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	tl.Warn(message)
}

// Error logs error level message
func (tl *TraceLogger) Error(message string) {
	log.Println(tl.formatMessage(ErrorLevel, message))
	// Also write to stderr for error visibility
	fmt.Fprintf(os.Stderr, "%s\n", tl.formatMessage(ErrorLevel, message))
}

// Errorf logs formatted error level message
func (tl *TraceLogger) Errorf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	tl.Error(message)
}

// Debug logs debug level message
func (tl *TraceLogger) Debug(message string) {
	if gin.Mode() == gin.DebugMode {
		log.Println(tl.formatMessage(DebugLevel, message))
	}
}

// Debugf logs formatted debug level message
func (tl *TraceLogger) Debugf(format string, args ...interface{}) {
	message := fmt.Sprintf(format, args...)
	tl.Debug(message)
}

// WithError logs error with additional context
func (tl *TraceLogger) WithError(err error, context string) {
	if err != nil {
		tl.Errorf("%s: %v", context, err)
	}
}

// LogDatabaseQuery logs database queries with trace ID
func (tl *TraceLogger) LogDatabaseQuery(query string, duration time.Duration) {
	tl.Infof("DB Query executed in %v: %s", duration, query)
}

// LogAPICall logs external API calls
func (tl *TraceLogger) LogAPICall(method, url string, statusCode int, duration time.Duration) {
	tl.Infof("API Call: %s %s - Status: %d, Duration: %v", method, url, statusCode, duration)
}

// LogUserAction logs user actions for audit trail
func (tl *TraceLogger) LogUserAction(userID, action, resource string) {
	tl.Infof("User Action: UserID=%s, Action=%s, Resource=%s", userID, action, resource)
}
