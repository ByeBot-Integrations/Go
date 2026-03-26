# Ray CAPTCHA - Go Demo

Example Go application demonstrating Ray CAPTCHA integration.

## Quick Start

```bash
go run main.go
```

Open http://localhost:4343

## Configuration

Edit the constants at the top of `main.go`:

```go
const (
    baseURL = "https://challenge.byebot.de"  // Ray CAPTCHA server URL
    apiKey  = ""                                  // Your API key
    siteKey = ""                                  // Your site key
    port    = "4343"                              // Server port
)
```

## How It Works

1. HTML form includes the captcha widget via `<div class="captcha-widget" data-sitekey="...">`
2. Widget script (`/ray/widget.js`) renders the captcha and adds a hidden `byebot-token` field on success
3. Form submits to `/submit`, server validates token via POST to `/validate_token`
4. Server returns success/error result

## Files

- `main.go` - HTTP server with form handling and token validation
- `templates/index.html` - Login form with captcha widget
