package benchmark

import (
	"sort"
	"time"
)

// Result tek bir HTTP isteğinin sonucunu temsil eder
type Result struct {
	Duration     time.Duration
	StatusCode   int
	BytesRead    int64
	Error        error
	DNSDuration  time.Duration
	ConnDuration time.Duration
	TTFBDuration time.Duration
}

// Stats tüm benchmark sonuçlarının istatistiklerini içerir
type Stats struct {
	Total     int
	Successes int
	Failures  int

	TotalTime time.Duration
	Fastest   time.Duration
	Slowest   time.Duration
	Average   time.Duration
	RPS       float64

	SizeTotal int64

	P50 time.Duration
	P75 time.Duration
	P90 time.Duration
	P95 time.Duration
	P99 time.Duration

	StatusCodes map[int]int
	ErrorMsgs   map[string]int
}

// ComputeStats sonuç listesinden istatistik hesaplar
func ComputeStats(results []Result, totalDuration time.Duration) *Stats {
	stats := &Stats{
		Total:       len(results),
		StatusCodes: make(map[int]int),
		ErrorMsgs:   make(map[string]int),
		TotalTime:   totalDuration,
	}

	if len(results) == 0 {
		return stats
	}

	var latencies []time.Duration
	stats.Fastest = results[0].Duration
	stats.Slowest = results[0].Duration

	for _, r := range results {
		stats.SizeTotal += r.BytesRead

		if r.Error != nil {
			stats.Failures++
			stats.ErrorMsgs[r.Error.Error()]++
			continue
		}

		stats.Successes++
		stats.StatusCodes[r.StatusCode]++
		latencies = append(latencies, r.Duration)

		if r.Duration < stats.Fastest {
			stats.Fastest = r.Duration
		}
		if r.Duration > stats.Slowest {
			stats.Slowest = r.Duration
		}
	}

	if len(latencies) == 0 {
		return stats
	}

	sort.Slice(latencies, func(i, j int) bool { return latencies[i] < latencies[j] })

	var total time.Duration
	for _, l := range latencies {
		total += l
	}
	stats.Average = total / time.Duration(len(latencies))

	if totalDuration > 0 {
		stats.RPS = float64(stats.Total) / totalDuration.Seconds()
	}

	stats.P50 = percentile(latencies, 50)
	stats.P75 = percentile(latencies, 75)
	stats.P90 = percentile(latencies, 90)
	stats.P95 = percentile(latencies, 95)
	stats.P99 = percentile(latencies, 99)

	return stats
}

func percentile(sorted []time.Duration, p int) time.Duration {
	if len(sorted) == 0 {
		return 0
	}
	idx := int(float64(len(sorted)-1) * float64(p) / 100.0)
	return sorted[idx]
}
