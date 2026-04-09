package benchmark

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptrace"
	"sync"
	"time"
)

// Task bir worker'a gönderilen görev
type Task struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
}

// Worker bir goroutine içinde HTTP istekleri gönderir
type Worker struct {
	id      int
	client  *http.Client
	tasks   <-chan Task
	results chan<- Result
	wg      *sync.WaitGroup
}

func newWorker(id int, client *http.Client, tasks <-chan Task, results chan<- Result, wg *sync.WaitGroup) *Worker {
	return &Worker{
		id:      id,
		client:  client,
		tasks:   tasks,
		results: results,
		wg:      wg,
	}
}

func (w *Worker) run() {
	defer w.wg.Done()
	for task := range w.tasks {
		w.results <- w.executeWithTrace(task)
	}
}

func (w *Worker) executeWithTrace(task Task) Result {
	var (
		dnsStart  time.Time
		dnsDone   time.Time
		connStart time.Time
		connDone  time.Time
		ttfb      time.Time
	)

	trace := &httptrace.ClientTrace{
		DNSStart:             func(_ httptrace.DNSStartInfo) { dnsStart = time.Now() },
		DNSDone:              func(_ httptrace.DNSDoneInfo) { dnsDone = time.Now() },
		ConnectStart:         func(_, _ string) { connStart = time.Now() },
		ConnectDone:          func(_, _ string, _ error) { connDone = time.Now() },
		GotFirstResponseByte: func() { ttfb = time.Now() },
	}

	var bodyReader io.Reader
	if len(task.Body) > 0 {
		bodyReader = bytes.NewReader(task.Body)
	}

	req, err := http.NewRequest(task.Method, task.URL, bodyReader)
	if err != nil {
		return Result{Error: err}
	}
	req = req.WithContext(httptrace.WithClientTrace(req.Context(), trace))

	for k, v := range task.Headers {
		req.Header.Set(k, v)
	}

	start := time.Now()
	resp, err := w.client.Do(req)
	if err != nil {
		return Result{Duration: time.Since(start), Error: err}
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return Result{Duration: time.Since(start), StatusCode: resp.StatusCode, Error: err}
	}
	duration := time.Since(start)

	var dnsDur, connDur, ttfbDur time.Duration
	if !dnsStart.IsZero() && !dnsDone.IsZero() {
		dnsDur = dnsDone.Sub(dnsStart)
	}
	if !connStart.IsZero() && !connDone.IsZero() {
		connDur = connDone.Sub(connStart)
	}
	if !ttfb.IsZero() {
		ttfbDur = ttfb.Sub(start)
	}

	return Result{
		Duration:     duration,
		StatusCode:   resp.StatusCode,
		BytesRead:    int64(len(body)),
		DNSDuration:  dnsDur,
		ConnDuration: connDur,
		TTFBDuration: ttfbDur,
	}
}
