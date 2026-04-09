package benchmark

import (
	"context"
	"crypto/tls"
	"net/http"
	"sync"
	"time"

	"golang.org/x/net/http2"
)

// Config benchmark yapılandırmasını tanımlar
type Config struct {
	URL       string
	Method    string
	Headers   map[string]string
	Body      []byte
	N         int
	C         int
	Timeout   time.Duration
	RateLimit float64 // istekler/saniye, 0=sınırsız
}

// BenchmarkRunner yük testini orkestre eder
type BenchmarkRunner struct {
	config Config
}

func NewRunner(cfg Config) *BenchmarkRunner {
	if cfg.Method == "" {
		cfg.Method = "GET"
	}
	if cfg.Timeout == 0 {
		cfg.Timeout = 20 * time.Second
	}
	if cfg.C <= 0 {
		cfg.C = 1
	}
	if cfg.N <= 0 {
		cfg.N = 1
	}
	return &BenchmarkRunner{config: cfg}
}

func (r *BenchmarkRunner) Run(ctx context.Context) (*Stats, error) {
	cfg := r.config

	// HTTP client oluştur (hey pattern)
	maxIdle := cfg.C
	if maxIdle > 100 {
		maxIdle = 100
	}
	tr := &http.Transport{
		MaxIdleConnsPerHost: maxIdle,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
	}
	_ = http2.ConfigureTransport(tr)

	client := &http.Client{
		Transport: tr,
		Timeout:   cfg.Timeout,
	}

	tasks := make(chan Task, cfg.C)
	results := make(chan Result, cfg.N)

	var wg sync.WaitGroup
	wg.Add(cfg.C)

	// Worker'ları başlat
	for i := 0; i < cfg.C; i++ {
		w := newWorker(i, client, tasks, results, &wg)
		go w.run()
	}

	start := time.Now()

	// Task üretici goroutine
	go func() {
		task := Task{
			Method:  cfg.Method,
			URL:     cfg.URL,
			Headers: cfg.Headers,
			Body:    cfg.Body,
		}

		var ticker *time.Ticker
		if cfg.RateLimit > 0 {
			interval := time.Duration(float64(time.Second) / cfg.RateLimit)
			ticker = time.NewTicker(interval)
			defer ticker.Stop()
		}

		for i := 0; i < cfg.N; i++ {
			select {
			case <-ctx.Done():
				close(tasks)
				return
			default:
			}

			if ticker != nil {
				<-ticker.C
			}
			tasks <- task
		}
		close(tasks)
	}()

	// Worker'ların bitmesini bekle ve results kanalını kapat
	go func() {
		wg.Wait()
		close(results)
	}()

	// Sonuçları topla
	collected := make([]Result, 0, cfg.N)
	for result := range results {
		collected = append(collected, result)
	}

	totalDuration := time.Since(start)
	return ComputeStats(collected, totalDuration), nil
}
