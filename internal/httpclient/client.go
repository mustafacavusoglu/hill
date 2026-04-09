package httpclient

import (
	"bytes"
	"crypto/tls"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/http2"
)

type Request struct {
	Method  string
	URL     string
	Headers map[string]string
	Body    []byte
	Timeout time.Duration
}

type Response struct {
	StatusCode int
	Status     string
	Headers    http.Header
	Body       []byte
	Duration   time.Duration
	Proto      string
	BodySize   int64
}

func Execute(req Request) (*Response, error) {
	timeout := req.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	tr := &http.Transport{
		MaxIdleConnsPerHost: 100,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: false},
	}
	_ = http2.ConfigureTransport(tr)

	client := &http.Client{
		Transport: tr,
		Timeout:   timeout,
	}

	var bodyReader io.Reader
	if len(req.Body) > 0 {
		bodyReader = bytes.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, err
	}

	for k, v := range req.Headers {
		httpReq.Header.Set(k, v)
	}

	start := time.Now()
	httpResp, err := client.Do(httpReq)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	body, err := io.ReadAll(httpResp.Body)
	if err != nil {
		return nil, err
	}
	duration := time.Since(start)

	return &Response{
		StatusCode: httpResp.StatusCode,
		Status:     httpResp.Status,
		Headers:    httpResp.Header,
		Body:       body,
		Duration:   duration,
		Proto:      httpResp.Proto,
		BodySize:   int64(len(body)),
	}, nil
}
