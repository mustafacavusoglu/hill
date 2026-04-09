# hill

**hill** — terminal için Postman + [hey](https://github.com/rakyll/hey) karışımı HTTP istemcisi. REST API istekleri atmak, yük testi yapmak ve ağ bağlantısı kontrol etmek için tek araç.

> [English documentation](README.md)

```
hill get https://api.example.com/users
hill post https://api.example.com/users -d '{"name":"ali"}' -H "Authorization: Bearer token"
hill benchmark https://api.example.com/users -n 2000 -c 100
hill check api.example.com
hill                          # → TUI modu
```

---

## Kurulum

```bash
git clone https://github.com/mustafacavusoglu/hill
cd hill
go build -o hill .
sudo mv hill /usr/local/bin/
```

Go 1.21+ gereklidir.

---

## TUI Modu

`hill` argümansız çalıştırıldığında tam ekran interaktif arayüz açılır.

```
┌─────────────────────────────────────────────────────────────────────────────────┐
│  hill  │ [F1] Request │ [F2] Response │ [F3] History │ [?] Yardım               │
├──────────────────────────────────────┬──────────────────────────────────────────┤
│  REQUEST                      [aktif]│  RESPONSE                                │
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
│  [ctrl+r] Gönder [ctrl+m] Method     │  ]                                       │
│  [tab] Alan geçiş                    │  [j/k] Scroll  [c] Kopyala               │
├──────────────────────────────────────┴──────────────────────────────────────────┤
│  HISTORY                                                                        │
│  ▶ GET    https://api.github.com/users            200  312ms  20:14:32          │
│    POST   https://api.example.com/users           201   89ms  20:13:11          │
│    GET    https://httpbin.org/get                 200  145ms  20:12:05          │
│    DELETE https://api.example.com/users/42        404   44ms  20:11:58          │
└─────────────────────────────────────────────────────────────────────────────────┘
  [ctrl+r] Gönder  [ctrl+m] Method  [F1-F3] Panel  [q] Çıkış
```

### TUI Panelleri

| Panel | Kısayol | Açıklama |
|-------|---------|----------|
| **Request** | `F1` | HTTP metodu, URL ve JSON body düzenleme |
| **Response** | `F2` | Yanıt status, headers ve body görüntüleme |
| **History** | `F3` | Yapılan tüm isteklerin listesi |

### TUI Klavye Kısayolları

| Kısayol | Eylem |
|---------|-------|
| `ctrl+r` | İsteği gönder |
| `ctrl+m` | HTTP metodunu değiştir (GET → POST → PUT → DELETE → PATCH → HEAD) |
| `tab` | Sonraki alana geç (URL ↔ Body) |
| `shift+tab` | Önceki alana geç |
| `F1` | Request paneline odaklan |
| `F2` | Response paneline odaklan |
| `F3` | History paneline odaklan |
| `j` / `↓` | Response body aşağı kaydır |
| `k` / `↑` | Response body yukarı kaydır |
| `c` | Response body'yi panoya kopyala |
| `enter` | History'den seçili isteği yükle |
| `q` / `ctrl+c` | Çıkış |

---

## CLI Komutları

### `hill get`

HTTP GET isteği gönderir.

```
hill get <url> [flags]
```

**Flags:**

| Flag | Kısa | Varsayılan | Açıklama |
|------|------|-----------|---------|
| `--header` | `-H` | — | Header ekle (tekrarlanabilir) |
| `--timeout` | `-t` | `30s` | İstek zaman aşımı |

**Örnekler:**

```bash
# Basit GET
hill get https://api.github.com/users

# Header ile
hill get https://api.example.com/profile \
  -H "Authorization: Bearer eyJhbGciOiJIUzI1..."

# Birden fazla header
hill get https://api.example.com/data \
  -H "Authorization: Bearer token" \
  -H "Accept: application/json" \
  -H "X-Request-ID: abc123"

# Timeout ile
hill get https://slow.example.com/endpoint -t 5s
```

**Çıktı:**
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

HTTP POST isteği gönderir. Body varsa ve `Content-Type` header'ı set edilmemişse otomatik `application/json` atanır.

```
hill post <url> [flags]
```

**Flags:**

| Flag | Kısa | Varsayılan | Açıklama |
|------|------|-----------|---------|
| `--data` | `-d` | — | Request body |
| `--header` | `-H` | — | Header ekle (tekrarlanabilir) |
| `--timeout` | `-t` | `30s` | İstek zaman aşımı |

**Örnekler:**

```bash
# JSON body ile POST
hill post https://api.example.com/users \
  -d '{"name": "Mustafa", "email": "m@example.com"}'

# Authorization header ile
hill post https://api.example.com/posts \
  -H "Authorization: Bearer token" \
  -d '{"title": "Yeni yazı", "body": "İçerik...", "published": true}'

# İç içe JSON
hill post https://api.example.com/orders \
  -d '{
    "customer": {"id": 42, "name": "Ali"},
    "items": [
      {"sku": "ABC-001", "qty": 2, "price": 19.99},
      {"sku": "XYZ-999", "qty": 1, "price": 49.90}
    ],
    "shipping": "express"
  }'

# Form data (Content-Type override)
hill post https://api.example.com/login \
  -H "Content-Type: application/x-www-form-urlencoded" \
  -d "username=ali&password=secret"
```

---

### `hill put`

HTTP PUT isteği gönderir. `post` ile aynı flags'lere sahiptir.

```
hill put <url> [flags]
```

**Örnekler:**

```bash
# Kayıt güncelle
hill put https://api.example.com/users/42 \
  -H "Authorization: Bearer token" \
  -d '{"name": "Mustafa Yeni", "role": "admin"}'
```

---

### `hill delete`

HTTP DELETE isteği gönderir.

```
hill delete <url> [flags]
```

**Flags:**

| Flag | Kısa | Varsayılan | Açıklama |
|------|------|-----------|---------|
| `--header` | `-H` | — | Header ekle |
| `--timeout` | `-t` | `30s` | İstek zaman aşımı |

**Örnekler:**

```bash
# Kayıt sil
hill delete https://api.example.com/users/42 \
  -H "Authorization: Bearer token"
```

---

### `hill benchmark`

HTTP yük testi — [hey](https://github.com/rakyll/hey) benzeri worker pool mimarisiyle çalışır. Her istek için DNS, TCP bağlantı ve TTFB süreleri `httptrace` ile ölçülür.

```
hill benchmark <url> [flags]
```

**Flags:**

| Flag | Kısa | Varsayılan | Açıklama |
|------|------|-----------|---------|
| `--number` | `-n` | `200` | Toplam istek sayısı |
| `--concurrency` | `-c` | `50` | Eşzamanlı bağlantı (worker) sayısı |
| `--method` | `-m` | `GET` | HTTP metodu |
| `--data` | `-d` | — | Request body |
| `--header` | `-H` | — | Header (tekrarlanabilir) |
| `--timeout` | `-t` | `20s` | Tek istek zaman aşımı |
| `--qps` | `-q` | `0` | Saniyedeki istek limiti (0 = sınırsız) |

**Örnekler:**

```bash
# Temel yük testi
hill benchmark https://api.example.com/health -n 1000 -c 50

# POST endpoint testi
hill benchmark https://api.example.com/users \
  -n 500 -c 20 \
  -m POST \
  -d '{"name": "test"}' \
  -H "Authorization: Bearer token"

# Rate limiting ile (saniyede max 100 istek)
hill benchmark https://api.example.com/search -n 1000 -c 10 -q 100

# Uzun timeout gerektiren endpoint
hill benchmark https://api.example.com/report -n 100 -c 5 -t 60s
```

**Çıktı:**
```
hill benchmark: 1000 istek, 50 eşzamanlı → https://api.example.com/health

── Benchmark Sonuçları ─────────────────────────────
  Toplam İstek:        1000
  Başarılı:            998
  Başarısız:           2
  Toplam Süre:         4.832s
  RPS:                 206.97
  Transfer:            142080 bytes (28.69 KB/s)

── Latency İstatistikleri ──────────────────────────
  En Hızlı:            18.2ms
  En Yavaş:            1.204s
  Ortalama:            241.7ms

  Dağılım:
  P50:                 198.4ms
  P75:                 312.1ms
  P90:                 489.3ms
  P95:                 621.8ms
  P99:                 987.2ms

── HTTP Status Dağılımı ────────────────────────────
  [200]  ████████████████████████████░░  998 (99.8%)
  [503]  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░    2 (0.2%)
```

---

### `hill check`

Bir IP adresi veya hostname için ağ bağlantı kontrolü yapar. DNS çözümleme, TCP erişilebilirlik ve ICMP ping testlerini **paralel** olarak çalıştırır.

```
hill check <ip-veya-host> [flags]
```

**Flags:**

| Flag | Kısa | Varsayılan | Açıklama |
|------|------|-----------|---------|
| `--port` | `-p` | `80, 443` | Özel port kontrolü |

**Örnekler:**

```bash
# Domain kontrolü (DNS + TCP 80/443 + ICMP)
hill check google.com

# IP adresi kontrolü
hill check 8.8.8.8

# Özel port (örn. PostgreSQL, Redis, MySQL)
hill check db.internal    -p 5432
hill check cache.internal -p 6379
hill check mysql.internal -p 3306

# Birden fazla host kontrol (shell loop)
for host in api.example.com db.example.com cache.example.com; do
  hill check $host
  echo
done
```

**Çıktı:**
```
── hill check: google.com ──────────────────────────

  DNS:
    ✓ Çözümlendi (13.7ms)
      → 142.250.185.78
      → 2a00:1450:4017:821::200e
    • PTR: [fra16s52-in-f14.1e100.net.]
    • MX:  [smtp.google.com.]

  TCP:
    ✓ Port 80  açık (29.5ms)
    ✓ Port 443 açık (28.9ms)

  ICMP Ping:
    ⚠ ICMP kullanılamıyor (sudo gerekli)
```

> **Not:** ICMP ping macOS ve Linux'ta root yetkisi gerektirir. `sudo hill check <host>` ile kullanılabilir. Yetki yoksa araç hata vermez, sadece uyarı gösterir.

---

## Proje Yapısı

```
hill/
├── main.go
├── go.mod
├── cmd/
│   ├── root.go          # cobra root; argümansız → TUI
│   ├── get.go           # hill get, post, put, delete
│   ├── benchmark.go     # hill benchmark
│   └── check.go         # hill check
├── internal/
│   ├── httpclient/
│   │   ├── client.go    # HTTP/2 destekli HTTP engine
│   │   └── formatter.go # JSON pretty-print, renkli çıktı
│   ├── benchmark/
│   │   ├── worker.go    # httptrace ile worker pool
│   │   ├── runner.go    # BenchmarkRunner orchestration
│   │   ├── result.go    # Stats ve percentile hesaplama
│   │   └── reporter.go  # Renkli benchmark raporu
│   ├── checker/
│   │   ├── dns.go       # DNS çözümleme, PTR, MX
│   │   ├── tcp.go       # TCP bağlantı testi
│   │   ├── icmp.go      # ICMP ping
│   │   └── checker.go   # Paralel orchestration + çıktı
│   └── tui/
│       ├── model.go     # Bubbletea root model
│       ├── styles.go    # Lipgloss stil tanımları
│       ├── keys.go      # Tüm klavye kısayolları
│       └── panels/
│           ├── request.go   # Method/URL/body paneli
│           ├── response.go  # Yanıt görüntüleme paneli
│           └── history.go   # Geçmiş istekler tablosu
```

---

## Bağımlılıklar

| Paket | Kullanım |
|-------|---------|
| [cobra](https://github.com/spf13/cobra) | CLI framework |
| [bubbletea](https://github.com/charmbracelet/bubbletea) | TUI framework (Elm architecture) |
| [lipgloss](https://github.com/charmbracelet/lipgloss) | Terminal renk ve stil |
| [bubbles](https://github.com/charmbracelet/bubbles) | textinput, textarea, viewport, table, spinner |
| [golang.org/x/net/http2](https://pkg.go.dev/golang.org/x/net/http2) | HTTP/2 transport |
