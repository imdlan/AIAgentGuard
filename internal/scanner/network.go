package scanner

import (
	"fmt"
	"net"
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
