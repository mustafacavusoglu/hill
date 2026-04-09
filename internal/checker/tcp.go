package checker

import (
	"fmt"
	"net"
	"time"
)

func checkTCP(host string, port int, timeout time.Duration) (reachable bool, latency time.Duration, err error) {
	addr := fmt.Sprintf("%s:%d", host, port)
	start := time.Now()
	conn, dialErr := net.DialTimeout("tcp", addr, timeout)
	latency = time.Since(start)
	if dialErr != nil {
		return false, latency, dialErr
	}
	conn.Close()
	return true, latency, nil
}
