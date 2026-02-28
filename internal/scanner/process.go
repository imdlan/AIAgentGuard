package scanner

import (
	"bufio"
	"fmt"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)

// ProcessInfo represents information about a running process
type ProcessInfo struct {
	PID        int
	UserName   string
	Command    string
	CPUUsage   float64
	ConnectsTo []string // Remote addresses the process connects to
}

// SuspiciousProcess represents a process that may be malicious
type SuspiciousProcess struct {
	ProcessInfo
	Reason string
	Risk   model.RiskLevel
}

// ScanProcesses scans for suspicious processes and network connections
func ScanProcesses() model.RiskLevel {
	var suspiciousProcesses []SuspiciousProcess

	// Get running processes
	processes, err := getRunningProcesses()
	if err != nil {
		return model.Low
	}

	// Check for reverse shells
	reverseShells := detectReverseShells(processes)
	suspiciousProcesses = append(suspiciousProcesses, reverseShells...)

	// Check for high CPU usage (potential crypto mining)
	highCPU := detectHighCPUUsage(processes)
	suspiciousProcesses = append(suspiciousProcesses, highCPU...)

	// Check for suspicious process names
	suspiciousNames := detectSuspiciousProcessNames(processes)
	suspiciousProcesses = append(suspiciousProcesses, suspiciousNames...)

	// Determine overall risk level
	if len(suspiciousProcesses) == 0 {
		return model.Low
	}

	// Check for critical risks
	for _, proc := range suspiciousProcesses {
		if proc.Risk == model.High || proc.Risk == model.Critical {
			return model.High
		}
	}

	if len(suspiciousProcesses) > 2 {
		return model.Medium
	}

	return model.Low
}

// ScanProcessesDetailed scans for suspicious processes and returns detailed information
func ScanProcessesDetailed() (model.RiskLevel, []model.RiskDetail) {
	details := []model.RiskDetail{}

	// Get running processes
	processes, err := getRunningProcesses()
	if err != nil {
		return model.Low, details
	}

	// Check for reverse shells
	reverseShells := detectReverseShells(processes)
	// Check for high CPU usage
	highCPU := detectHighCPUUsage(processes)
	// Check for suspicious process names
	suspiciousNames := detectSuspiciousProcessNames(processes)

	// Combine all suspicious processes
	allSuspicious := append(reverseShells, highCPU...)
	allSuspicious = append(allSuspicious, suspiciousNames...)

	if len(allSuspicious) == 0 {
		return model.Low, details
	}

	// Convert to model.ProcessInfo format
	suspiciousProcs := []model.ProcessInfo{}
	for _, proc := range allSuspicious {
		// Extract process name from command
		nameParts := strings.Fields(proc.Command)
		name := "unknown"
		if len(nameParts) > 0 {
			name = nameParts[0]
		}
		
		suspiciousProcs = append(suspiciousProcs, model.ProcessInfo{
			PID:         proc.PID,
			Name:        name,
			CommandLine: proc.Command,
			User:        proc.UserName,
			RiskReason:  proc.Reason,
		})
	}

	// Determine overall risk
	risk := model.Low
	for _, proc := range allSuspicious {
		if proc.Risk == model.High || proc.Risk == model.Critical {
			risk = model.High
			break
		}
	}
	if risk == model.Low && len(allSuspicious) > 2 {
		risk = model.Medium
	}

	// Build detail
	detail := model.RiskDetail{
		Type:        risk,
		Category:    "process",
		Description: fmt.Sprintf("发现 %d 个可疑进程", len(allSuspicious)),
		Details: model.RiskSpecificInfo{
			SuspiciousProcs: suspiciousProcs,
		},
	}

	if risk == model.High || risk == model.Critical {
		detail.Remediation = model.RemediationInfo{
			Summary: "立即调查并终止可疑进程",
			Steps: []model.RemediationStep{
				{
					Step:        1,
					Action:      "查看进程详情",
					Command:     "ps -p <PID> -f",
					Explanation: "查看进程的完整命令行和启动信息",
				},
				{
					Step:        2,
					Action:      "终止可疑进程",
					Command:     "kill <PID>",
					Explanation: "如果确认是恶意进程，终止它",
				},
			},
			Commands:    []string{"ps aux | grep <process_name>"},
			Priority:    "CRITICAL",
			RiskAfter:    model.Low,
		}
	}

	details = append(details, detail)
	return risk, details
}


// getRunningProcesses retrieves a list of running processes
func getRunningProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	switch runtime.GOOS {
	case "darwin", "linux":
		return getUnixProcesses()
	case "windows":
		return getWindowsProcesses()
	default:
		return processes, fmt.Errorf("unsupported platform")
	}
}

// getUnixProcesses retrieves processes on Unix-like systems
func getUnixProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	// Use ps command to get process list
	cmd := exec.Command("ps", "aux")
	output, err := cmd.Output()
	if err != nil {
		return processes, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	lineNum := 0
	for scanner.Scan() {
		lineNum++
		// Skip header
		if lineNum == 1 {
			continue
		}

		line := scanner.Text()
		fields := strings.Fields(line)
		if len(fields) < 11 {
			continue
		}

		// Parse ps aux output
		// USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
		pid, err := strconv.Atoi(fields[1])
		if err != nil {
			continue
		}

		cpu, err := strconv.ParseFloat(fields[2], 64)
		if err != nil {
			continue
		}

		// Join the command part (may contain spaces)
		command := strings.Join(fields[10:], " ")

		processes = append(processes, ProcessInfo{
			PID:      pid,
			UserName: fields[0],
			Command:  command,
			CPUUsage: cpu,
		})
	}

	return processes, nil
}

// getWindowsProcesses retrieves processes on Windows
func getWindowsProcesses() ([]ProcessInfo, error) {
	var processes []ProcessInfo

	// Use tasklist command
	cmd := exec.Command("tasklist", "/fo", "csv", "/nh")
	output, err := cmd.Output()
	if err != nil {
		return processes, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		// Parse CSV output
		fields := strings.Split(strings.Trim(line, "\""), "\",\"")
		if len(fields) < 5 {
			continue
		}

		// Extract PID
		pid, err := strconv.Atoi(strings.Trim(fields[1], "\""))
		if err != nil {
			continue
		}

		processes = append(processes, ProcessInfo{
			PID:     pid,
			Command: fields[0],
		})
	}

	return processes, nil
}

// detectReverseShells detects processes that may be reverse shells
func detectReverseShells(processes []ProcessInfo) []SuspiciousProcess {
	var suspicious []SuspiciousProcess

	// Check for network connections
	connections, err := getNetworkConnections()
	if err != nil {
		return suspicious
	}

	for _, proc := range processes {
		// Check if process is a shell
		if isShellCommand(proc.Command) {
			// Check if shell has external network connections
			for _, conn := range connections {
				if conn.PID == proc.PID && isExternalConnection(conn.RemoteAddr) {
					suspicious = append(suspicious, SuspiciousProcess{
						ProcessInfo: proc,
						Reason:      fmt.Sprintf("reverse shell detected: connects to %s", conn.RemoteAddr),
						Risk:        model.Critical,
					})
				}
			}
		}

		// Check for suspicious command patterns
		if isReverseShellPattern(proc.Command) {
			suspicious = append(suspicious, SuspiciousProcess{
				ProcessInfo: proc,
				Reason:      "reverse shell pattern detected in command",
				Risk:        model.High,
			})
		}
	}

	return suspicious
}

// NetworkConnection represents a network connection
type NetworkConnection struct {
	PID        int
	Protocol   string
	LocalAddr  string
	RemoteAddr string
	State      string
}

// getNetworkConnections retrieves active network connections
func getNetworkConnections() ([]NetworkConnection, error) {
	var connections []NetworkConnection

	switch runtime.GOOS {
	case "darwin", "linux":
		return getUnixNetworkConnections()
	case "windows":
		return getWindowsNetworkConnections()
	default:
		return connections, fmt.Errorf("unsupported platform")
	}
}

// getUnixNetworkConnections retrieves network connections on Unix
func getUnixNetworkConnections() ([]NetworkConnection, error) {
	var connections []NetworkConnection

	// Use netstat or ss command
	cmd := exec.Command("netstat", "-anp")
	output, err := cmd.Output()
	if err != nil {
		// Try ss as fallback
		cmd = exec.Command("ss", "-tuln")
		output, err = cmd.Output()
		if err != nil {
			return connections, err
		}
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		// Parse netstat output
		// Proto Local Address           Foreign Address         State       PID/Program name
		if len(fields) < 6 {
			continue
		}

		protocol := fields[0]
		localAddr := fields[3]
		remoteAddr := fields[4]
		state := fields[5]

		var pid int
		if len(fields) > 6 {
			pidStr := strings.Split(fields[6], "/")[0]
			pid, _ = strconv.Atoi(pidStr)
		}

		connections = append(connections, NetworkConnection{
			PID:        pid,
			Protocol:   protocol,
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
			State:      state,
		})
	}

	return connections, nil
}

// getWindowsNetworkConnections retrieves network connections on Windows
func getWindowsNetworkConnections() ([]NetworkConnection, error) {
	var connections []NetworkConnection

	// Use netstat command
	cmd := exec.Command("netstat", "-ano")
	output, err := cmd.Output()
	if err != nil {
		return connections, err
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)

		if len(fields) < 5 {
			continue
		}

		protocol := fields[0]
		localAddr := fields[1]
		remoteAddr := fields[2]
		state := fields[3]

		pid, _ := strconv.Atoi(fields[4])

		connections = append(connections, NetworkConnection{
			PID:        pid,
			Protocol:   protocol,
			LocalAddr:  localAddr,
			RemoteAddr: remoteAddr,
			State:      state,
		})
	}

	return connections, nil
}

// isShellCommand checks if a command is a shell
func isShellCommand(cmd string) bool {
	shells := []string{"bash", "sh", "zsh", "fish", "dash", "tcsh", "csh", "cmd", "powershell"}
	cmdLower := strings.ToLower(filepath.Base(cmd))
	for _, shell := range shells {
		if strings.Contains(cmdLower, shell) {
			return true
		}
	}
	return false
}

// isExternalConnection checks if an address is external (not localhost/private)
func isExternalConnection(addr string) bool {
	// Skip empty addresses
	if addr == "" || addr == "*:*" {
		return false
	}

	// Check for localhost patterns
	localhostPatterns := []string{"127.0.0.1", "[::1]", "0.0.0.0", "[::]", "localhost"}
	for _, pattern := range localhostPatterns {
		if strings.Contains(addr, pattern) {
			return false
		}
	}

	// Check for private network ranges
	privateRanges := []string{"192.168.", "10.", "172.16.", "172.17.", "172.18.", "172.19.",
		"172.20.", "172.21.", "172.22.", "172.23.", "172.24.", "172.25.",
		"172.26.", "172.27.", "172.28.", "172.29.", "172.30.", "172.31."}
	for _, range_ := range privateRanges {
		if strings.HasPrefix(addr, range_) {
			return false
		}
	}

	// If not local/private, consider it external
	return true
}

// isReverseShellPattern checks for common reverse shell patterns
func isReverseShellPattern(cmd string) bool {
	patterns := []string{
		"bash -i",
		"nc -l",
		"netcat -l",
		"socat TCP",
		"powershell -enc",
		"cmd /c",
	}

	cmdLower := strings.ToLower(cmd)
	for _, pattern := range patterns {
		if strings.Contains(cmdLower, pattern) {
			return true
		}
	}

	return false
}

// detectHighCPUUsage detects processes with unusually high CPU usage
func detectHighCPUUsage(processes []ProcessInfo) []SuspiciousProcess {
	var suspicious []SuspiciousProcess

	for _, proc := range processes {
		// High CPU threshold: > 80% for non-system processes
		if proc.CPUUsage > 80.0 {
			// Check if it's a suspicious process (not known system tools)
			if !isKnownSystemProcess(proc.Command) {
				suspicious = append(suspicious, SuspiciousProcess{
					ProcessInfo: proc,
					Reason:      fmt.Sprintf("high CPU usage: %.1f%%", proc.CPUUsage),
					Risk:        model.Medium,
				})
			}
		}
	}

	return suspicious
}

// detectSuspiciousProcessNames detects processes with suspicious names
func detectSuspiciousProcessNames(processes []ProcessInfo) []SuspiciousProcess {
	var suspicious []SuspiciousProcess

	suspiciousKeywords := []string{
		"miner",
		"cryptonight",
		"xmr",
		"bitcoin",
		"backdoor",
		"trojan",
		"malware",
		"rat",
		"keylog",
		"inject",
	}

	for _, proc := range processes {
		cmdLower := strings.ToLower(proc.Command)
		for _, keyword := range suspiciousKeywords {
			if strings.Contains(cmdLower, keyword) {
				suspicious = append(suspicious, SuspiciousProcess{
					ProcessInfo: proc,
					Reason:      fmt.Sprintf("suspicious keyword: %s", keyword),
					Risk:        model.High,
				})
				break
			}
		}
	}

	return suspicious
}

// isKnownSystemProcess checks if a process is a known system process
func isKnownSystemProcess(cmd string) bool {
	systemProcesses := []string{
		"kernel_task",
		"WindowServer",
		"launchd",
		"systemd",
		"init",
		"kthreadd",
		"ssh",
		"docker",
		"chrome",
		"firefox",
		"safari",
	}

	cmdLower := strings.ToLower(filepath.Base(cmd))
	for _, proc := range systemProcesses {
		if strings.Contains(cmdLower, proc) {
			return true
		}
	}

	return false
}

// GetSuspiciousProcesses returns detailed list of suspicious processes
func GetSuspiciousProcesses() []SuspiciousProcess {
	var suspicious []SuspiciousProcess

	processes, err := getRunningProcesses()
	if err != nil {
		return suspicious
	}

	reverseShells := detectReverseShells(processes)
	suspicious = append(suspicious, reverseShells...)

	highCPU := detectHighCPUUsage(processes)
	suspicious = append(suspicious, highCPU...)

	suspiciousNames := detectSuspiciousProcessNames(processes)
	suspicious = append(suspicious, suspiciousNames...)

	return suspicious
}
