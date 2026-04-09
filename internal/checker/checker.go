package checker

import (
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// CheckResult ağ kontrolü sonuçlarını içerir
type CheckResult struct {
	Host string

	IPs        []string
	DNSLatency time.Duration
	DNSError   error
	PTRRecords []string
	MXRecords  []string

	TCPPorts []TCPProbe

	PingReachable bool
	PingLatency   time.Duration
	PingError     error
}

type TCPProbe struct {
	Port      int
	Reachable bool
	Latency   time.Duration
	Error     error
}

// Run tüm ağ kontrollerini paralel olarak çalıştırır
func Run(host string, customPort int) (*CheckResult, error) {
	result := &CheckResult{Host: host}
	var mu sync.Mutex
	var wg sync.WaitGroup

	// DNS
	wg.Add(1)
	go func() {
		defer wg.Done()
		ips, lat, ptrs, mxs, err := resolveDNS(host)
		mu.Lock()
		result.IPs = ips
		result.DNSLatency = lat
		result.DNSError = err
		result.PTRRecords = ptrs
		result.MXRecords = mxs
		mu.Unlock()
	}()

	// TCP port kontrolü
	ports := []int{80, 443}
	if customPort > 0 {
		ports = []int{customPort}
	}
	for _, port := range ports {
		port := port
		wg.Add(1)
		go func() {
			defer wg.Done()
			reachable, lat, err := checkTCP(host, port, 5*time.Second)
			mu.Lock()
			result.TCPPorts = append(result.TCPPorts, TCPProbe{
				Port:      port,
				Reachable: reachable,
				Latency:   lat,
				Error:     err,
			})
			mu.Unlock()
		}()
	}

	// ICMP ping
	wg.Add(1)
	go func() {
		defer wg.Done()
		reachable, lat, err := pingICMP(host, 3*time.Second)
		mu.Lock()
		result.PingReachable = reachable
		result.PingLatency = lat
		result.PingError = err
		mu.Unlock()
	}()

	wg.Wait()
	return result, nil
}

var (
	styleTitle  = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleOK     = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	styleErr    = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	styleWarn   = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	styleBold   = lipgloss.NewStyle().Bold(true)
	styleDim    = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
	styleLabel  = lipgloss.NewStyle().Foreground(lipgloss.Color("14"))
)

// PrintResult check sonuçlarını formatlı yazdırır
func PrintResult(r *CheckResult) {
	fmt.Println(styleTitle.Render("── hill check: " + r.Host + " ─────────────────────────"))
	fmt.Println()

	// DNS
	fmt.Println(styleBold.Render("  DNS:"))
	if r.DNSError != nil {
		fmt.Printf("    %s %v\n", styleErr.Render("✗"), r.DNSError)
	} else {
		fmt.Printf("    %s Çözümlendi (%s)\n", styleOK.Render("✓"), r.DNSLatency.Round(100000))
		for _, ip := range r.IPs {
			fmt.Printf("      %s %s\n", styleDim.Render("→"), ip)
		}
		if len(r.PTRRecords) > 0 {
			fmt.Printf("    %s PTR: %v\n", styleDim.Render("•"), r.PTRRecords)
		}
		if len(r.MXRecords) > 0 {
			fmt.Printf("    %s MX: %v\n", styleDim.Render("•"), r.MXRecords)
		}
	}
	fmt.Println()

	// TCP
	fmt.Println(styleBold.Render("  TCP:"))
	for _, probe := range r.TCPPorts {
		if probe.Reachable {
			fmt.Printf("    %s Port %s açık (%s)\n",
				styleOK.Render("✓"),
				styleLabel.Render(fmt.Sprintf("%d", probe.Port)),
				probe.Latency.Round(100000),
			)
		} else {
			msg := ""
			if probe.Error != nil {
				msg = " — " + probe.Error.Error()
			}
			fmt.Printf("    %s Port %s kapalı%s\n",
				styleErr.Render("✗"),
				styleLabel.Render(fmt.Sprintf("%d", probe.Port)),
				msg,
			)
		}
	}
	fmt.Println()

	// ICMP
	fmt.Println(styleBold.Render("  ICMP Ping:"))
	if r.PingError != nil {
		if errors.Is(r.PingError, ErrICMPPermission) {
			fmt.Printf("    %s ICMP kullanılamıyor (sudo gerekli)\n", styleWarn.Render("⚠"))
		} else {
			fmt.Printf("    %s %v\n", styleErr.Render("✗"), r.PingError)
		}
	} else if r.PingReachable {
		fmt.Printf("    %s Erişilebilir (%s)\n", styleOK.Render("✓"), r.PingLatency.Round(100000))
	} else {
		fmt.Printf("    %s Yanıt yok\n", styleErr.Render("✗"))
	}
	fmt.Println()
}
