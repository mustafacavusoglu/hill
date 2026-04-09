package checker

import (
	"encoding/binary"
	"errors"
	"net"
	"os"
	"time"
)

// ErrICMPPermission ICMP için gerekli yetkinin olmadığını belirtir
var ErrICMPPermission = errors.New("ICMP ping için yükseltilmiş yetki gerekiyor (sudo)")

func pingICMP(host string, timeout time.Duration) (reachable bool, latency time.Duration, err error) {
	// IP adresini çöz
	addrs, resolveErr := net.LookupHost(host)
	if resolveErr != nil || len(addrs) == 0 {
		return false, 0, resolveErr
	}
	ip := addrs[0]

	conn, connErr := net.ListenPacket("ip4:icmp", "0.0.0.0")
	if connErr != nil {
		if os.IsPermission(connErr) {
			return false, 0, ErrICMPPermission
		}
		return false, 0, connErr
	}
	defer conn.Close()

	dst, resolveAddrErr := net.ResolveIPAddr("ip4", ip)
	if resolveAddrErr != nil {
		return false, 0, resolveAddrErr
	}

	// ICMP Echo Request paketi oluştur
	msg := makeICMPEcho(1, 1)

	if err := conn.SetDeadline(time.Now().Add(timeout)); err != nil {
		return false, 0, err
	}

	start := time.Now()
	if _, writeErr := conn.WriteTo(msg, dst); writeErr != nil {
		return false, 0, writeErr
	}

	buf := make([]byte, 256)
	for {
		n, _, readErr := conn.ReadFrom(buf)
		if readErr != nil {
			return false, time.Since(start), readErr
		}
		// ICMP reply kontrolü: tip 0 = Echo Reply
		if n >= 1 && buf[0] == 0 {
			return true, time.Since(start), nil
		}
	}
}

// makeICMPEcho basit bir ICMP Echo Request paketi oluşturur
func makeICMPEcho(id, seq uint16) []byte {
	msg := make([]byte, 8)
	msg[0] = 8 // Echo Request
	msg[1] = 0 // Code
	binary.BigEndian.PutUint16(msg[4:], id)
	binary.BigEndian.PutUint16(msg[6:], seq)
	// Checksum
	cs := icmpChecksum(msg)
	binary.BigEndian.PutUint16(msg[2:], cs)
	return msg
}

func icmpChecksum(data []byte) uint16 {
	var sum uint32
	for i := 0; i < len(data)-1; i += 2 {
		sum += uint32(data[i])<<8 | uint32(data[i+1])
	}
	if len(data)%2 != 0 {
		sum += uint32(data[len(data)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xffff) + (sum >> 16)
	}
	return ^uint16(sum)
}
