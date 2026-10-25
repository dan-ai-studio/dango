package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

// Response 通用 JSON 响应结构
type Response struct {
	Message   string    `json:"message"`
	Timestamp time.Time `json:"timestamp"`
}

func main() {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /", homeHandler)
	mux.HandleFunc("GET /health", healthHandler)
	mux.HandleFunc("GET /api/hello", helloHandler)
	mux.HandleFunc("GET /api/time", timeHandler)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	addr := ":" + port
	log.Printf("🚀 dango server listening on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// homeHandler 首页
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	fmt.Fprint(w, `<!DOCTYPE html>
<html lang="zh-CN">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Dango - Go Web 示例</title>
  <style>
    body { font-family: -apple-system, BlinkMacSystemFont, "Segoe UI", Roboto, sans-serif; max-width: 720px; margin: 60px auto; padding: 0 20px; color: #222; }
    h1 { color: #00ADD8; }
    code { background: #f4f4f4; padding: 2px 6px; border-radius: 4px; }
    ul { line-height: 1.8; }
    a { color: #00ADD8; }
  </style>
</head>
<body>
  <h1>🥟 Dango</h1>
  <p>一个简单的 Go Web 示例项目。</p>
  <h2>可用接口</h2>
  <ul>
    <li><a href="/health">/health</a> - 健康检查</li>
    <li><a href="/api/hello">/api/hello</a> - JSON 问候</li>
    <li><a href="/api/time">/api/time</a> - 当前服务器时间</li>
  </ul>
  <p>源码：<a href="https://github.com/dan-ai-studio/dango">github.com/dan-ai-studio/dango</a></p>
</body>
</html>`)
}

// healthHandler 健康检查
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// helloHandler JSON 问候
func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(Response{
		Message:   "Hello from Dango! 🥟",
		Timestamp: time.Now(),
	})
}

// timeHandler 返回当前服务器时间
func timeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"time": time.Now().Format(time.RFC3339),
	})
}
