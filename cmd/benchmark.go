package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/mustafacavusoglu/hill/internal/benchmark"
	"github.com/spf13/cobra"
)

var benchmarkCmd = &cobra.Command{
	Use:   "benchmark <url>",
	Short: "HTTP yük testi (hey benzeri)",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		n, _ := cmd.Flags().GetInt("number")
		c, _ := cmd.Flags().GetInt("concurrency")
		method, _ := cmd.Flags().GetString("method")
		data, _ := cmd.Flags().GetString("data")
		headers, _ := cmd.Flags().GetStringArray("header")
		timeout, _ := cmd.Flags().GetDuration("timeout")
		qps, _ := cmd.Flags().GetFloat64("qps")

		cfg := benchmark.Config{
			URL:       args[0],
			Method:    method,
			Headers:   parseHeaders(headers),
			Body:      []byte(data),
			N:         n,
			C:         c,
			Timeout:   timeout,
			RateLimit: qps,
		}

		fmt.Printf("hill benchmark: %d istek, %d eşzamanlı → %s\n\n", n, c, args[0])

		start := time.Now()
		runner := benchmark.NewRunner(cfg)
		stats, err := runner.Run(context.Background())
		if err != nil {
			fmt.Fprintf(os.Stderr, "Benchmark hatası: %v\n", err)
			os.Exit(1)
		}
		_ = start

		benchmark.PrintStats(stats)
		return nil
	},
}

func init() {
	benchmarkCmd.Flags().IntP("number", "n", 200, "Toplam istek sayısı")
	benchmarkCmd.Flags().IntP("concurrency", "c", 50, "Eşzamanlı bağlantı sayısı")
	benchmarkCmd.Flags().StringP("method", "m", "GET", "HTTP metodu")
	benchmarkCmd.Flags().StringP("data", "d", "", "Request body")
	benchmarkCmd.Flags().StringArrayP("header", "H", nil, "Header")
	benchmarkCmd.Flags().DurationP("timeout", "t", 20*time.Second, "İstek timeout")
	benchmarkCmd.Flags().Float64P("qps", "q", 0, "Saniyedeki istek sayısı limiti (0=sınırsız)")
	rootCmd.AddCommand(benchmarkCmd)
}
