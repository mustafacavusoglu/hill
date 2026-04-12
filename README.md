# hill

**hill** — a Postman + [hey](https://github.com/rakyll/hey) hybrid for the terminal. Send REST API requests, run load tests, and check network connectivity — all from one tool.

> [Türkçe dokümantasyon için tıklayın](README.tr.md)

```
hill get https://api.example.com/users
hill post https://api.example.com/users -d '{"name":"ali"}' -H "Authorization: Bearer token"
hill benchmark https://api.example.com/users -n 2000 -c 100
hill check api.example.com
hill                          # → TUI mode
```

---

## Installation

```bash
git clone https://github.com/mustafacavusoglu/hill
cd hill
go build -o hill .
sudo mv hill /usr/local/bin/
```

Requires Go 1.21+.

---

## TUI Mode

Running `hill` with no arguments opens a full-screen interactive interface.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  hill  │ [F1] Request │ [F2] Response │ [F3] History │ [?] Help                 │
├──────────────────────────────────────┬──────────────────────────────────────────┤
│  REQUEST                     [active]│  RESPONSE                                │
│                                      │                                          │
│  GET   https://api.github.com/users  │  200 OK  (312ms)  HTTP/2.0  1.4 KB       │
│                                      │                                          │
│  Body (JSON):                        │  Content-Type: application/json          │
│  ┌────────────────────────────────┐  │  X-RateLimit-Remaining: 59               │
│  │ {                              │  │                                          │
│  │   "filter": "active",          │  │  [                                       │
│  │   "limit": 10                  │  │    {                                     │
│  │ }                              │  │      "login": "torvalds",                │
│  │                                │  │      "id": 1024,                         │
│  └────────────────────────────────┘  │      "type": "User"                      │
│                                      │    }                                     │
│  [ctrl+r] Send  [ctrl+m] Method      │  ]                                       │
│  [tab] Switch field                  │  [j/k] Scroll  [c] Copy                  │
├──────────────────────────────────────┴──────────────────────────────────────────┤
│  HISTORY                                                                        │
│  ▶ GET    https://api.github.com/users            200  312ms  20:14:32          │
│    POST   https://api.example.com/users           201   89ms  20:13:11          │
│    GET    https://httpbin.org/get                 200  145ms  20:12:05          │
│    DELETE https://api.example.com/users/42        404   44ms  20:11:58          │
└─────────────────────────────────────────────────────────────────────────────────┘
  [ctrl+r] Send  [ctrl+m] Method  [F1-F3] Panels  [q] Quit
```

### Panels

| Panel | Shortcut | Description |
|-------|----------|-------------|
| **Request** | `F1` | Edit HTTP method, URL, and JSON body |
| **Response** | `F2` | View response status, headers, and body |
| **History** | `F3` | List of all previous requests |

### Keyboard Shortcuts

| Key | Action |
|-----|--------|
| `ctrl+r` | Send request |
| `ctrl+m` | Cycle HTTP method (GET → POST → PUT → DELETE → PATCH → HEAD) |
| `tab` | Next field (URL ↔ Body) |
| `shift+tab` | Previous field |
| `F1` | Focus Request panel |
| `F2` | Focus Response panel |
| `F3` | Focus History panel |
| `j` / `↓` | Scroll response body down |
| `k` / `↑` | Scroll response body up |
| `c` | Copy response body to clipboard |
| `enter` | Load selected history entry into Request panel |
| `q` / `ctrl+c` | Quit |

---

## CLI Commands

### `hill get`

Send an HTTP GET request.

```
hill get <url> [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--header` | `-H` | — | Add a header (repeatable) |
| `--timeout` | `-t` | `30s` | Request timeout |

**Examples:**

```bash
# Simple GET
hill get https://api.github.com/users

# With authorization header
hill get https://api.example.com/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1..."

# Multiple headers
hill get https://api.example.com/data \
  -H "Authorization: Bearer token" \
  -H "Accept: application/json" \
  -H "X-Request-ID: abc123"

# Custom timeout
hill get https://slow.example.com/endpoint -t 5s
```

**Output:**
```
200 OK  (312ms)  HTTP/2.0  1248 bytes

Content-Type: application/json
X-RateLimit-Limit: 60
X-RateLimit-Remaining: 59

{
  "login": "torvalds",
  "id": 1024,
  "type": "User"
}
```

---

### `hill post`

Send an HTTP POST request. If a body is provided and no `Content-Type` header is set, it defaults to `application/json`.

```
hill post <url> [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--data` | `-d` | — | Request body |
| `--header` | `-H` | — | Add a header (repeatable) |
| `--timeout` | `-t` | `30s` | Request timeout |

**Examples:**

```bash
# JSON body
hill post https://api.example.com/users \
  -d '{"name": "Alice", "email": "alice@example.com"}'

# With authorization header
hill post https://api.example.com/posts \
  -H "Authorization: Bearer token" \
  -d '{"title": "New post", "body": "Content...", "published": true}'

# Nested JSON
hill post https://api.example.com/orders \
  -d '{
    "customer": {"id": 42, "name": "Alice"},
    "items": [
      {"sku": "ABC-001", "qty": 2, "price": 19.99},
      {"sku": "XYZ-999", "qty": 1, "price": 49.90}
    ],
    "shipping": "express"
  }'

# Form data (override Content-Type)
hill post https://api.example.com/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=alice&password=secret"
```

---

### `hill put`

Send an HTTP PUT request. Accepts the same flags as `post`.

```
hill put <url> [flags]
```

**Examples:**

```bash
# Update a record
hill put https://api.example.com/users/42 \
  -H "Authorization: Bearer token" \
  -d '{"name": "Alice Updated", "role": "admin"}'
```

---

### `hill delete`

Send an HTTP DELETE request.

```
hill delete <url> [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--header` | `-H` | — | Add a header |
| `--timeout` | `-t` | `30s` | Request timeout |

**Examples:**

```bash
# Delete a record
hill delete https://api.example.com/users/42 \
  -H "Authorization: Bearer token"
```

---

### `hill benchmark`

HTTP load testing — uses a concurrent worker pool architecture (inspired by [hey](https://github.com/rakyll/hey)). Each request is instrumented with `net/http/httptrace` for DNS, TCP connection, and TTFB breakdowns.

```
hill benchmark <url> [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--number` | `-n` | `200` | Total number of requests |
| `--concurrency` | `-c` | `50` | Number of concurrent workers |
| `--method` | `-m` | `GET` | HTTP method |
| `--data` | `-d` | — | Request body |
| `--header` | `-H` | — | Add a header (repeatable) |
| `--timeout` | `-t` | `20s` | Per-request timeout |
| `--qps` | `-q` | `0` | Rate limit in requests/sec (0 = unlimited) |

**Examples:**

```bash
# Basic load test
hill benchmark https://api.example.com/health -n 1000 -c 50

# Load test a POST endpoint
hill benchmark https://api.example.com/users \
  -n 500 -c 20 \
  -m POST \
  -d '{"name": "test"}' \
  -H "Authorization: Bearer token"

# Rate-limited test (max 100 req/s)
hill benchmark https://api.example.com/search -n 1000 -c 10 -q 100

# Endpoint with slow responses
hill benchmark https://api.example.com/report -n 100 -c 5 -t 60s
```

**Output:**
```
hill benchmark: 1000 requests, 50 concurrent → https://api.example.com/health

── Summary ─────────────────────────────────────────
  Total Requests:      1000
  Succeeded:           998
  Failed:              2
  Total Duration:      4.832s
  RPS:                 206.97
  Transfer:            142080 bytes (28.69 KB/s)

── Latency ─────────────────────────────────────────
  Fastest:             18.2ms
  Slowest:             1.204s
  Average:             241.7ms

  Distribution:
  P50:                 198.4ms
  P75:                 312.1ms
  P90:                 489.3ms
  P95:                 621.8ms
  P99:                 987.2ms

── HTTP Status Codes ────────────────────────────────
  [200]  ████████████████████████████░░  998 (99.8%)
  [503]  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    2 (0.2%)
```

---

### `hill check`

Network connectivity check for an IP address or hostname. Runs DNS resolution, TCP reachability, and ICMP ping **in parallel**.

```
hill check <ip-or-host> [flags]
```

**Flags:**

| Flag | Short | Default | Description |
|------|-------|---------|-------------|
| `--port` | `-p` | `80, 443` | Custom port to probe |

**Examples:**

```bash
# Full check for a domain (DNS + TCP 80/443 + ICMP)
hill check google.com

# IP address check
hill check 8.8.8.8

# Custom port (e.g. PostgreSQL, Redis, MySQL)
hill check db.internal    -p 5432
hill check cache.internal -p 6379
hill check mysql.internal -p 3306

# Check multiple hosts (shell loop)
for host in api.example.com db.example.com cache.example.com; do
  hill check $host
  echo
done
```

**Output:**
```
── hill check: google.com ──────────────────────────

  DNS:
    ✓ Resolved (13.7ms)
      → 142.250.185.78
      → 2a00:1450:4017:821::200e
    • PTR: [fra16s52-in-f14.1e100.net.]
    • MX:  [smtp.google.com.]

  TCP:
    ✓ Port 80  open (29.5ms)
    ✓ Port 443 open (28.9ms)

  ICMP Ping:
    ⚠ ICMP unavailable (requires elevated privileges)
```

> **Note:** ICMP ping requires root on macOS and Linux. Run `sudo hill check <host>` to enable it. Without permission the tool continues gracefully and shows a warning instead of failing.

---

## Project Structure

```
hill/
├── main.go
├── go.mod
├── cmd/
│   ├── root.go          # cobra root; no args → TUI
│   ├── get.go           # hill get, post, put, delete
│   ├── benchmark.go     # hill benchmark
│   └── check.go         # hill check
├── internal/
│   ├── httpclient/
│   │   ├── client.go    # HTTP/2-enabled HTTP engine
│   │   └── formatter.go # JSON pretty-print, colored output
│   ├── benchmark/
│   │   ├── worker.go    # Worker pool with httptrace instrumentation
│   │   ├── runner.go    # BenchmarkRunner orchestration
│   │   ├── result.go    # Stats and percentile calculation
│   │   └── reporter.go  # Colored benchmark report
│   ├── checker/
│   │   ├── dns.go       # DNS resolution, PTR, MX records
│   │   ├── tcp.go       # TCP connectivity probe
│   │   ├── icmp.go      # ICMP ping
│   │   └── checker.go   # Parallel orchestration + output
│   └── tui/
│       ├── model.go     # Bubbletea root model
│       ├── styles.go    # Lipgloss style definitions
│       ├── keys.go      # All keyboard bindings
│       └── panels/
│           ├── request.go   # Method / URL / body panel
│           ├── response.go  # Response viewer panel
│           └── history.go   # Request history table
```

--- 

## Dependencies

| Package | Purpose |
|---------|---------|
| [cobra](https://github.com/spf13/cobra) | CLI framework |
| [bubbletea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal colors and layout |
| [bubbles](https://github.com/charmbracelet/bubbles) | textinput, textarea, viewport, table, spinner |
| [golang.org/x/net/http2](https://pkg.go.dev/golang.org/x/net/http2) | HTTP/2 transport |
