package audit

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

// LogLevel represents the severity level of an audit event
type LogLevel string

const (
	// DebugLevel for detailed debugging information
	DebugLevel LogLevel = "DEBUG"
	// InfoLevel for general informational messages
	InfoLevel LogLevel = "INFO"
	// WarnLevel for warning messages
	WarnLevel LogLevel = "WARN"
	// ErrorLevel for error events
	ErrorLevel LogLevel = "ERROR"
	// CriticalLevel for critical security events
	CriticalLevel LogLevel = "CRITICAL"
)

// EventType represents the type of audit event
type EventType string

const (
	// PolicyViolation event when a security policy is violated
	PolicyViolation EventType = "policy_violation"
	// CommandBlocked event when a command is blocked
	CommandBlocked EventType = "command_blocked"
	// FileAccessDenied event when file access is denied
	FileAccessDenied EventType = "file_access_denied"
	// NetworkAccessDenied event when network access is denied
	NetworkAccessDenied EventType = "network_access_denied"
	// ScanCompleted event when a security scan completes
	ScanCompleted EventType = "scan_completed"
	// RiskDetected event when a security risk is detected
	RiskDetected EventType = "risk_detected"
	// SystemEvent for general system events
	SystemEvent EventType = "system"
)

// AuditEvent represents a single audit log entry
type AuditEvent struct {
	Timestamp  time.Time              `json:"timestamp"`
	Level      LogLevel               `json:"level"`
	EventType  EventType              `json:"event_type"`
	Message    string                 `json:"message"`
	Details    map[string]interface{} `json:"details,omitempty"`
	Hostname   string                 `json:"hostname"`
	ProcessID  int                    `json:"pid"`
	UserName   string                 `json:"user,omitempty"`
	SourceIP   string                 `json:"source_ip,omitempty"`
	Command    string                 `json:"command,omitempty"`
	Path       string                 `json:"path,omitempty"`
	PolicyRule string                 `json:"policy_rule,omitempty"`
	RiskLevel  string                 `json:"risk_level,omitempty"`
}

// AuditLogger handles audit logging with multiple output options
type AuditLogger struct {
	mu            sync.Mutex
	fileLogger    *log.Logger
	fileHandle    *os.File
	consoleLogger *log.Logger
	enabled       bool
	jsonFormat    bool
	minLevel      LogLevel
}

var (
	// Global logger instance
	defaultLogger *AuditLogger
	once          sync.Once
)

// Init initializes the global audit logger
func Init(logFilePath string, jsonFormat bool, minLevel LogLevel) error {
	var initErr error
	once.Do(func() {
		defaultLogger = &AuditLogger{
			consoleLogger: log.New(os.Stdout, "", 0),
			enabled:       true,
			jsonFormat:    jsonFormat,
			minLevel:      minLevel,
		}

		if logFilePath != "" {
			// Create directory if it doesn't exist
			logDir := filepath.Dir(logFilePath)
			if err := os.MkdirAll(logDir, 0755); err != nil {
				initErr = fmt.Errorf("failed to create log directory: %w", err)
				return
			}

			// Open log file
			file, err := os.OpenFile(logFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
			if err != nil {
				initErr = fmt.Errorf("failed to open log file: %w", err)
				return
			}

			defaultLogger.fileHandle = file
			defaultLogger.fileLogger = log.New(file, "", 0)
		}
	})

	return initErr
}

// Close closes the audit logger and releases resources
func Close() error {
	if defaultLogger != nil && defaultLogger.fileHandle != nil {
		return defaultLogger.fileHandle.Close()
	}
	return nil
}

// shouldLog determines if an event should be logged based on its level
func (al *AuditLogger) shouldLog(level LogLevel) bool {
	levels := map[LogLevel]int{
		DebugLevel:    0,
		InfoLevel:     1,
		WarnLevel:     2,
		ErrorLevel:    3,
		CriticalLevel: 4,
	}

	currentLevel, ok := levels[level]
	if !ok {
		return true
	}
	minLevelVal, ok := levels[al.minLevel]
	if !ok {
		return true
	}

	return currentLevel >= minLevelVal
}

// LogEvent logs an audit event
func LogEvent(event AuditEvent) error {
	if defaultLogger == nil || !defaultLogger.enabled {
		return nil
	}

	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()

	if !defaultLogger.shouldLog(event.Level) {
		return nil
	}

	// Set timestamp if not provided
	if event.Timestamp.IsZero() {
		event.Timestamp = time.Now().UTC()
	}

	// Set hostname and process ID if not provided
	if event.Hostname == "" {
		if hostname, err := os.Hostname(); err == nil {
			event.Hostname = hostname
		} else {
			event.Hostname = "unknown"
		}
	}

	if event.ProcessID == 0 {
		event.ProcessID = os.Getpid()
	}

	// Set username if not provided
	if event.UserName == "" {
		if user := os.Getenv("USER"); user != "" {
			event.UserName = user
		} else if user := os.Getenv("USERNAME"); user != "" {
			event.UserName = user
		}
	}

	// Format the log entry
	var logEntry string
	if defaultLogger.jsonFormat {
		jsonBytes, err := json.Marshal(event)
		if err != nil {
			return fmt.Errorf("failed to marshal audit event: %w", err)
		}
		logEntry = string(jsonBytes)
	} else {
		logEntry = formatPlainText(event)
	}

	// Write to file logger if configured
	if defaultLogger.fileLogger != nil {
		defaultLogger.fileLogger.Println(logEntry)
	}

	// Write to console logger
	defaultLogger.consoleLogger.Println(logEntry)

	return nil
}

// formatPlainText formats an audit event as plain text
func formatPlainText(event AuditEvent) string {
	timestamp := event.Timestamp.Format("2006-01-02T15:04:05Z")
	return fmt.Sprintf("[%s] %s %s: %s", timestamp, event.Level, event.EventType, event.Message)
}

// Convenience functions for common event types

// LogPolicyViolation logs a security policy violation
func LogPolicyViolation(policyType, resource, reason string) error {
	event := AuditEvent{
		Level:     ErrorLevel,
		EventType: PolicyViolation,
		Message:   fmt.Sprintf("Policy violation: %s access denied for %s", policyType, resource),
		Details: map[string]interface{}{
			"policy_type": policyType,
			"resource":    resource,
			"reason":      reason,
		},
		Path:       resource,
		PolicyRule: policyType,
	}

	return LogEvent(event)
}

// LogCommandBlocked logs when a command is blocked
func LogCommandBlocked(command, reason string) error {
	event := AuditEvent{
		Level:     WarnLevel,
		EventType: CommandBlocked,
		Message:   fmt.Sprintf("Command blocked: %s", command),
		Details: map[string]interface{}{
			"command": command,
			"reason":  reason,
		},
		Command: command,
	}

	return LogEvent(event)
}

// LogRiskDetected logs when a security risk is detected
func LogRiskDetected(category, description, riskLevel string) error {
	event := AuditEvent{
		Level:     ErrorLevel,
		EventType: RiskDetected,
		Message:   fmt.Sprintf("Risk detected: %s - %s", category, description),
		Details: map[string]interface{}{
			"category":    category,
			"description": description,
		},
		RiskLevel: riskLevel,
	}

	return LogEvent(event)
}

// LogScanCompleted logs when a security scan completes
func LogScanCompleted(scanType string, results map[string]string) error {
	event := AuditEvent{
		Level:     InfoLevel,
		EventType: ScanCompleted,
		Message:   fmt.Sprintf("Scan completed: %s", scanType),
		Details: map[string]interface{}{
			"scan_type": scanType,
			"results":   results,
		},
	}

	return LogEvent(event)
}

// LogSystemEvent logs a general system event
func LogSystemEvent(level LogLevel, message string) error {
	event := AuditEvent{
		Level:     level,
		EventType: SystemEvent,
		Message:   message,
	}

	return LogEvent(event)
}

// GetLogger returns the default audit logger instance
func GetLogger() *AuditLogger {
	return defaultLogger
}

// SetEnabled enables or disables audit logging
func SetEnabled(enabled bool) {
	if defaultLogger != nil {
		defaultLogger.mu.Lock()
		defer defaultLogger.mu.Unlock()
		defaultLogger.enabled = enabled
	}
}

// IsEnabled returns whether audit logging is enabled
func IsEnabled() bool {
	if defaultLogger != nil {
		defaultLogger.mu.Lock()
		defer defaultLogger.mu.Unlock()
		return defaultLogger.enabled
	}
	return false
}

// RotateLog rotates the log file by renaming and creating a new one
func RotateLog() error {
	if defaultLogger == nil || defaultLogger.fileHandle == nil {
		return fmt.Errorf("no log file configured")
	}

	defaultLogger.mu.Lock()
	defer defaultLogger.mu.Unlock()

	oldPath := defaultLogger.fileHandle.Name()
	newPath := oldPath + "." + time.Now().Format("20060102-150405")

	// Close current file
	if err := defaultLogger.fileHandle.Close(); err != nil {
		return fmt.Errorf("failed to close current log file: %w", err)
	}

	// Rename current file
	if err := os.Rename(oldPath, newPath); err != nil {
		return fmt.Errorf("failed to rotate log file: %w", err)
	}

	// Open new log file
	file, err := os.OpenFile(oldPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0640)
	if err != nil {
		return fmt.Errorf("failed to open new log file: %w", err)
	}

	defaultLogger.fileHandle = file
	defaultLogger.fileLogger = log.New(file, "", 0)

	return nil
}

// GetLogFilePath returns the current log file path
func GetLogFilePath() string {
	if defaultLogger != nil && defaultLogger.fileHandle != nil {
		return defaultLogger.fileHandle.Name()
	}
	return ""
}

// InitDefault initializes audit logging with default settings
// Logs to ~/.agent-guard/audit.log in JSON format
func InitDefault() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		homeDir = os.TempDir()
	}

	logDir := filepath.Join(homeDir, ".agent-guard")
	logPath := filepath.Join(logDir, "audit.log")

	return Init(logPath, true, InfoLevel)
}

func init() {
	// Auto-initialize with default settings on package import
	// This can be overridden by explicit Init() call
	if runtime.GOOS != "windows" {
		_ = InitDefault()
	}
}
