package checker

import (
	"context"
	"net"
	"time"
)

func resolveDNS(host string) (ips []string, latency time.Duration, ptrs []string, mxs []string, err error) {
	resolver := net.DefaultResolver

	start := time.Now()
	addrs, lookupErr := resolver.LookupHost(context.Background(), host)
	latency = time.Since(start)

	if lookupErr != nil {
		return nil, latency, nil, nil, lookupErr
	}
	ips = addrs

	// PTR kayıtları (reverse DNS)
	for _, ip := range addrs {
		names, e := resolver.LookupAddr(context.Background(), ip)
		if e == nil {
			ptrs = append(ptrs, names...)
		}
	}

	// MX kayıtları
	mxRecords, e := resolver.LookupMX(context.Background(), host)
	if e == nil {
		for _, mx := range mxRecords {
			mxs = append(mxs, mx.Host)
		}
	}

	return ips, latency, ptrs, mxs, nil
}
