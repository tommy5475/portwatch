package scanner

import (
	"fmt"
	"net"
	"time"
)

// Port represents a network port with its protocol
type Port struct {
	Number   int
	Protocol string // "tcp" or "udp"
}

// OpenPort represents an open port with additional metadata
type OpenPort struct {
	Port
	Process   string
	Timestamp time.Time
}

// Scanner handles port scanning operations
type Scanner struct {
	timeout time.Duration
}

// New creates a new Scanner with the specified timeout
func New(timeout time.Duration) *Scanner {
	return &Scanner{
		timeout: timeout,
	}
}

// ScanTCPPort checks if a TCP port is open on localhost
func (s *Scanner) ScanTCPPort(port int) (bool, error) {
	address := fmt.Sprintf("127.0.0.1:%d", port)
	conn, err := net.DialTimeout("tcp", address, s.timeout)
	if err != nil {
		return false, nil
	}
	defer conn.Close()
	return true, nil
}

// ScanTCPRange scans a range of TCP ports and returns the open ones
func (s *Scanner) ScanTCPRange(startPort, endPort int) ([]OpenPort, error) {
	if startPort < 1 || endPort > 65535 || startPort > endPort {
		return nil, fmt.Errorf("invalid port range: %d-%d", startPort, endPort)
	}

	var openPorts []OpenPort
	for port := startPort; port <= endPort; port++ {
		if open, err := s.ScanTCPPort(port); err == nil && open {
			openPorts = append(openPorts, OpenPort{
				Port: Port{
					Number:   port,
					Protocol: "tcp",
				},
				Timestamp: time.Now(),
			})
		}
	}

	return openPorts, nil
}

// ScanCommonPorts scans commonly used ports
func (s *Scanner) ScanCommonPorts() ([]OpenPort, error) {
	commonPorts := []int{20, 21, 22, 23, 25, 53, 80, 110, 143, 443, 465, 587, 993, 995, 3000, 3306, 5432, 6379, 8000, 8080, 8443, 9000}
	var openPorts []OpenPort

	for _, port := range commonPorts {
		if open, err := s.ScanTCPPort(port); err == nil && open {
			openPorts = append(openPorts, OpenPort{
				Port: Port{
					Number:   port,
					Protocol: "tcp",
				},
				Timestamp: time.Now(),
			})
		}
	}

	return openPorts, nil
}
