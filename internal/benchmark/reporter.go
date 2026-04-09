package benchmark

import (
	"fmt"
	"sort"

	"github.com/charmbracelet/lipgloss"
)

var (
	styleTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("12"))
	styleOK      = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	styleWarn    = lipgloss.NewStyle().Foreground(lipgloss.Color("3"))
	styleErr     = lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	styleBold    = lipgloss.NewStyle().Bold(true)
	styleDim     = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

func PrintStats(s *Stats) {
	fmt.Println(styleTitle.Render("── Benchmark Sonuçları ────────────────────────────"))
	fmt.Println()

	// Özet
	fmt.Printf("  %-20s %s\n", "Toplam İstek:", styleBold.Render(fmt.Sprintf("%d", s.Total)))
	fmt.Printf("  %-20s %s\n", "Başarılı:", styleOK.Render(fmt.Sprintf("%d", s.Successes)))
	if s.Failures > 0 {
		fmt.Printf("  %-20s %s\n", "Başarısız:", styleErr.Render(fmt.Sprintf("%d", s.Failures)))
	}
	fmt.Printf("  %-20s %s\n", "Toplam Süre:", s.TotalTime.Round(1000000).String())
	fmt.Printf("  %-20s %s\n", "RPS:", styleBold.Render(fmt.Sprintf("%.2f", s.RPS)))
	fmt.Printf("  %-20s %d bytes (%.2f KB/s)\n", "Transfer:",
		s.SizeTotal,
		float64(s.SizeTotal)/s.TotalTime.Seconds()/1024,
	)
	fmt.Println()

	// Latency
	fmt.Println(styleTitle.Render("── Latency İstatistikleri ─────────────────────────"))
	fmt.Println()
	fmt.Printf("  %-20s %s\n", "En Hızlı:", styleOK.Render(s.Fastest.Round(100000).String()))
	fmt.Printf("  %-20s %s\n", "En Yavaş:", styleErr.Render(s.Slowest.Round(100000).String()))
	fmt.Printf("  %-20s %s\n", "Ortalama:", s.Average.Round(100000).String())
	fmt.Println()
	fmt.Println(styleDim.Render("  Dağılım:"))
	fmt.Printf("  %-20s %s\n", "P50:", s.P50.Round(100000).String())
	fmt.Printf("  %-20s %s\n", "P75:", s.P75.Round(100000).String())
	fmt.Printf("  %-20s %s\n", "P90:", s.P90.Round(100000).String())
	fmt.Printf("  %-20s %s\n", "P95:", s.P95.Round(100000).String())
	fmt.Printf("  %-20s %s\n", "P99:", s.P99.Round(100000).String())
	fmt.Println()

	// Status code dağılımı
	if len(s.StatusCodes) > 0 {
		fmt.Println(styleTitle.Render("── HTTP Status Dağılımı ────────────────────────────"))
		fmt.Println()
		codes := make([]int, 0, len(s.StatusCodes))
		for code := range s.StatusCodes {
			codes = append(codes, code)
		}
		sort.Ints(codes)
		for _, code := range codes {
			count := s.StatusCodes[code]
			style := styleOK
			if code >= 500 {
				style = styleErr
			} else if code >= 400 {
				style = styleWarn
			}
			bar := buildBar(count, s.Total, 30)
			fmt.Printf("  %s  %s %s\n",
				style.Render(fmt.Sprintf("[%d]", code)),
				bar,
				styleDim.Render(fmt.Sprintf("%d (%.1f%%)", count, float64(count)/float64(s.Total)*100)),
			)
		}
		fmt.Println()
	}

	// Hata özeti
	if len(s.ErrorMsgs) > 0 {
		fmt.Println(styleTitle.Render("── Hatalar ─────────────────────────────────────────"))
		fmt.Println()
		for msg, count := range s.ErrorMsgs {
			fmt.Printf("  %s  %s\n", styleErr.Render(fmt.Sprintf("x%d", count)), msg)
		}
		fmt.Println()
	}
}

func buildBar(count, total, width int) string {
	if total == 0 {
		return ""
	}
	filled := int(float64(count) / float64(total) * float64(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < filled {
			bar += "█"
		} else {
			bar += "░"
		}
	}
	return styleDim.Render(bar)
}
