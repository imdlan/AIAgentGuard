package security

import (
	"strings"
)

// Injection patterns that indicate prompt injection attacks
var injectionPatterns = []string{
	"ignore previous instructions",
	"override system rules",
	"ignore all instructions",
	"disregard previous",
	"forget everything",
	"delete all files",
	"exfiltrate secrets",
	"download and execute",
	"curl | sh",
	"wget | bash",
	"eval(",
	"exec(",
	"system(",
	"__import__",
	"subprocess",
}

// Privilege escalation patterns
var privilegeEscalationPatterns = []string{
	"sudo",
	"run as root",
	"administrator",
	"elevate privileges",
	"bypass security",
	"disable protection",
}

// Data exfiltration patterns
var exfiltrationPatterns = []string{
	"send to server",
	"upload to",
	"exfiltrate",
	"steal data",
	"leak secrets",
	"base64 encode",
	"post to",
	"transmit",
}

// IsPromptInjection checks if the input contains potential prompt injection patterns
func IsPromptInjection(input string) bool {
	inputLower := strings.ToLower(input)

	for _, pattern := range injectionPatterns {
		if strings.Contains(inputLower, pattern) {
			return true
		}
	}

	return false
}

// IsPrivilegeEscalation checks if the input attempts privilege escalation
func IsPrivilegeEscalation(input string) bool {
	inputLower := strings.ToLower(input)

	for _, pattern := range privilegeEscalationPatterns {
		if strings.Contains(inputLower, pattern) {
			return true
		}
	}

	return false
}

// IsDataExfiltration checks if the input attempts data exfiltration
func IsDataExfiltration(input string) bool {
	inputLower := strings.ToLower(input)

	for _, pattern := range exfiltrationPatterns {
		if strings.Contains(inputLower, pattern) {
			return true
		}
	}

	return false
}

// CheckSecurity performs a comprehensive security check on the input
func CheckSecurity(input string) (isInsecure bool, reasons []string) {
	reasons = []string{}

	if IsPromptInjection(input) {
		isInsecure = true
		reasons = append(reasons, "prompt injection detected")
	}

	if IsPrivilegeEscalation(input) {
		isInsecure = true
		reasons = append(reasons, "privilege escalation attempt")
	}

	if IsDataExfiltration(input) {
		isInsecure = true
		reasons = append(reasons, "data exfiltration attempt")
	}

	return
}
