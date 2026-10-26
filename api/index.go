package handler

import (
	"net/http"

	"github.com/dan-ai-studio/dango/internal/server"
)

var router = server.New()

// Handler 是 Vercel Go Serverless Function 的入口。
// 配合 vercel.json 的 rewrite，所有请求都转发到这里，
// 再由 Gin 引擎统一分发（首页 Vue 页面 + /api/* 接口）。
func Handler(w http.ResponseWriter, r *http.Request) {
	router.ServeHTTP(w, r)
}
