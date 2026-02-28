package scanner

import (
	"bufio"
	"fmt"
	"net"
	"os/exec"
	"strconv"
	"strings"
	"time"

	"github.com/imdlan/AIAgentGuard/pkg/model"
)
// ScanNetwork checks if network access is available
func ScanNetwork() model.RiskLevel {
	// Try to connect to a well-known external DNS server (8.8.8.8:53)
	timeout := 2 * time.Second
	conn, err := net.DialTimeout("tcp", "8.8.8.8:53", timeout)
	if err != nil {
		// Can't reach external network - check for local network access
		if checkLocalNetwork() {
			return model.Medium
		}
		return model.Low
	}
	conn.Close()

	// External network is accessible
	return model.Medium
}

// checkLocalNetwork checks if localhost/private network is accessible
func checkLocalNetwork() bool {
	// Try to connect to localhost
	conn, err := net.DialTimeout("tcp", "127.0.0.1:80", 1*time.Second)
	if err == nil {
		conn.Close()
		return true
	}

	// Try common private network addresses
	privateAddresses := []string{
		"192.168.1.1:80",
		"10.0.0.1:80",
		"172.16.0.1:80",
	}

	for _, addr := range privateAddresses {
		conn, err := net.DialTimeout("tcp", addr, 500*time.Millisecond)
		if err == nil {
			conn.Close()
			return true
		}
	}

	return false
}

// CheckSpecificHost checks if a specific host is reachable
func CheckSpecificHost(host string, port int, timeout time.Duration) bool {
	addr := fmt.Sprintf("%s:%d", host, port)
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

// ScanNetworkDetailed performs detailed network scanning and returns information about open ports and connections
func ScanNetworkDetailed() (model.RiskLevel, []model.RiskDetail) {
	details := []model.RiskDetail{}

	// Check basic network access
	basicRisk := ScanNetwork()
	if basicRisk == model.Low {
		return model.Low, details
	}

	// Get detailed information
	openPorts := detectOpenPorts()
	activeConns := detectActiveConnections()

	if len(openPorts) == 0 && len(activeConns) == 0 {
		// Network accessible but no detailed info available
		detail := model.RiskDetail{
			Type:        basicRisk,
			Category:    "network",
			Description: "External network access is available",
		}
		details = append(details, detail)
		return basicRisk, details
	}

	// Build detail with ports and connections
	detail := model.RiskDetail{
		Type:        basicRisk,
		Category:    "network",
		Description: fmt.Sprintf("Network access: %d open ports, %d active connections", len(openPorts), len(activeConns)),
		Details: model.RiskSpecificInfo{
			OpenPorts:   convertToPortDetails(openPorts),
			ActiveConns: convertToConnectionDetails(activeConns),
		},
	}

	details = append(details, detail)
	return basicRisk, details
}


// OpenPortInfo contains information about an open port
type OpenPortInfo struct {
	Protocol string
	Port     int
	Address  string
	State    string
}

// ConnectionInfo contains information about a network connection
type ConnectionInfo struct {
	Protocol   string
	LocalAddr  string
	RemoteAddr string
	State      string
}

// detectOpenPorts detects open ports on the local system
func detectOpenPorts() []OpenPortInfo {
	var ports []OpenPortInfo
	
	// Use lsof to get listening ports
	cmd := exec.Command("lsof", "-i", "-P", "-n", "-sTCP:LISTEN")
	output, err := cmd.Output()
	if err != nil {
		return ports
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		
		if len(fields) < 10 {
			continue
		}
		
		// Extract port from local address
		localAddr := fields[8]
		port := 0
		if strings.Contains(localAddr, ":") {
			portStr := strings.Split(localAddr, ":")[1]
			port, _ = strconv.Atoi(portStr)
		}
		
		if port > 0 {
			ports = append(ports, OpenPortInfo{
				Protocol: fields[7], // TCP/UDP
				Port:     port,
				Address:  localAddr,
				State:    "LISTEN",
			})
		}
	}
	
	return ports
}

// detectActiveConnections detects active network connections
func detectActiveConnections() []ConnectionInfo {
	var conns []ConnectionInfo
	
	cmd := exec.Command("netstat", "-an")
	output, err := cmd.Output()
	if err != nil {
		return conns
	}

	scanner := bufio.NewScanner(strings.NewReader(string(output)))
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		
		if len(fields) < 6 {
			continue
		}
		
		protocol := fields[0]
		localAddr := fields[3]
		remoteAddr := fields[4]
		state := fields[5]
		
		// Only include ESTABLISHED connections
		if state == "ESTABLISHED" {
			conns = append(conns, ConnectionInfo{
				Protocol:   protocol,
				LocalAddr:  localAddr,
				RemoteAddr: remoteAddr,
				State:      state,
			})
		}
	}
	
	return conns
}

// convertToPortDetails converts OpenPortInfo to model.PortDetail
func convertToPortDetails(ports []OpenPortInfo) []model.PortDetail {
	details := []model.PortDetail{}
	
	for _, port := range ports {
		service := "unknown"
		riskReason := "Open port may expose services"
		
		// Common service ports
		services := map[int]string{
			22:    "ssh",
			80:    "http",
			443:   "https",
			3306:  "mysql",
			5432:  "postgresql",
			6379:  "redis",
			27017: "mongodb",
			8080:  "http-proxy",
		}
		
		if s, ok := services[port.Port]; ok {
			service = s
		}
		
		// Determine risk
		if port.Port < 1024 {
			riskReason = "Privileged port - requires root access"
		}
		
		details = append(details, model.PortDetail{
			Port:       port.Port,
			Protocol:   port.Protocol,
			Service:    service,
			Process:    "N/A",
			RiskReason: riskReason,
		})
	}
	
	return details
}

// convertToConnectionDetails converts ConnectionInfo to model.ConnectionDetail
func convertToConnectionDetails(conns []ConnectionInfo) []model.ConnectionDetail {
	details := []model.ConnectionDetail{}
	
	for _, conn := range conns {
		details = append(details, model.ConnectionDetail{
			Protocol:    conn.Protocol,
			LocalAddr:   conn.LocalAddr,
			RemoteAddr:  conn.RemoteAddr,
			State:       conn.State,
		})
	}
	
	return details
}
