package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
)

// Configuration - update these values for your setup
const (
	baseURL = "https://challenge.byebot.de"
	apiKey  = "" // Your API key here
	siteKey = "bd1cc81b04564d3f899e" // Just an example sitekey so the widget shows up
	port    = "4343"
)

type TokenValidation struct {
	APIKey string `json:"api_key"`
	Token  string `json:"token"`
}

type PageData struct {
	BaseURL string
	SiteKey string
}

type ResultData struct {
	Success  bool
	Username string
	Message  string
}

func main() {
	http.HandleFunc("/", handleIndex)
	http.HandleFunc("/submit", handleSubmit)

	log.Printf("Server starting on http://localhost:%s", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, "Template error", http.StatusInternalServerError)
		log.Printf("Template error: %v", err)
		return
	}

	tmpl.Execute(w, PageData{
		BaseURL: baseURL,
		SiteKey: siteKey,
	})
}

func handleSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Failed to parse form", http.StatusBadRequest)
		return
	}

	username := r.FormValue("username")
	token := r.FormValue("byebot-token")

	if token == "" {
		renderResult(w, ResultData{
			Success: false,
			Message: "Captcha token missing",
		})
		return
	}

	valid, err := validateToken(token)
	if err != nil {
		renderResult(w, ResultData{
			Success: false,
			Message: fmt.Sprintf("Validation error: %v", err),
		})
		return
	}

	if !valid {
		renderResult(w, ResultData{
			Success: false,
			Message: "Captcha validation failed",
		})
		return
	}

	renderResult(w, ResultData{
		Success:  true,
		Username: username,
		Message:  "Login successful",
	})
}

func validateToken(token string) (bool, error) {
	payload := TokenValidation{
		APIKey: apiKey,
		Token:  token,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false, err
	}

	resp, err := http.Post(
		baseURL+"/validate_token",
		"application/json",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return false, err
	}
	defer resp.Body.Close()

	io.Copy(io.Discard, resp.Body)

	return resp.StatusCode >= 200 && resp.StatusCode < 300, nil
}

func renderResult(w http.ResponseWriter, data ResultData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	statusColor := "#ef4444"
	if data.Success {
		statusColor = "#10b981"
	}

	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
    <title>Result</title>
    <style>
        body {
            font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
            background: linear-gradient(135deg, #0f0f1a 0%%, #1a1a2e 100%%);
            color: #fff;
            min-height: 100vh;
            margin: 0;
            display: flex;
            align-items: center;
            justify-content: center;
        }
        .card {
            background: #1e1e2e;
            border-radius: 12px;
            padding: 2rem;
            max-width: 400px;
            text-align: center;
            border: 1px solid rgba(255,255,255,0.1);
        }
        .status {
            font-size: 1.5rem;
            font-weight: 600;
            color: %s;
            margin-bottom: 1rem;
        }
        .message { color: #94a3b8; margin-bottom: 1.5rem; }
        a {
            display: inline-block;
            padding: 0.75rem 1.5rem;
            background: linear-gradient(135deg, #6366f1 0%%, #8b5cf6 100%%);
            color: #fff;
            text-decoration: none;
            border-radius: 6px;
            font-weight: 600;
        }
    </style>
</head>
<body>
    <div class="card">
        <div class="status">%s</div>
        <div class="message">%s</div>
        <a href="/">Back to Form</a>
    </div>
</body>
</html>`, statusColor, data.Message, getUserInfo(data))

	w.Write([]byte(html))
}

func getUserInfo(data ResultData) string {
	if data.Success && data.Username != "" {
		return fmt.Sprintf("Welcome, %s!", data.Username)
	}
	if !data.Success {
		return "Please try again."
	}
	return ""
}

