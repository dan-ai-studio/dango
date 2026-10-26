package server

import (
	"embed"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

//go:embed web/index.html
var frontendFS embed.FS

// New 构建 Gin 引擎，同时被 main.go（本地/VPS）和 api/index.go（Vercel）复用
func New() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())

	// 简单 CORS（demo 允许跨域调试）
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type")
		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	store := NewStore()

	// 首页：Vue 应用（embed 打进二进制，Vercel 也能用）
	r.GET("/", func(c *gin.Context) {
		data, _ := frontendFS.ReadFile("web/index.html")
		c.Data(http.StatusOK, "text/html; charset=utf-8", data)
	})

	// 健康检查
	health := func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok", "service": "dango"})
	}
	r.GET("/health", health)

	api := r.Group("/api")
	api.GET("/health", health)
	api.GET("/stats", func(c *gin.Context) {
		total, done, pending := store.Stats()
		c.JSON(http.StatusOK, gin.H{"total": total, "done": done, "pending": pending})
	})

	// 任务列表（支持 ?q=关键词&done=true/false）
	api.GET("/tasks", func(c *gin.Context) {
		q := c.Query("q")
		var donePtr *bool
		if v := c.Query("done"); v != "" {
			if b, err := strconv.ParseBool(v); err == nil {
				donePtr = &b
			}
		}
		c.JSON(http.StatusOK, gin.H{"items": store.List(q, donePtr)})
	})

	// 新建任务
	api.POST("/tasks", func(c *gin.Context) {
		var body struct {
			Title string `json:"title" binding:"required"`
		}
		if err := c.ShouldBindJSON(&body); err != nil || body.Title == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "title 不能为空"})
			return
		}
		c.JSON(http.StatusCreated, store.Create(body.Title))
	})

	// 更新任务
	api.PUT("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id 非法"})
			return
		}
		var body struct {
			Title *string `json:"title"`
			Done  *bool   `json:"done"`
		}
		if err := c.ShouldBindJSON(&body); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "参数非法"})
			return
		}
		t, ok := store.Update(id, body.Title, body.Done)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	// 一键切换完成状态
	api.PATCH("/tasks/:id/toggle", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id 非法"})
			return
		}
		t, ok := store.Toggle(id)
		if !ok {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}
		c.JSON(http.StatusOK, t)
	})

	// 删除任务
	api.DELETE("/tasks/:id", func(c *gin.Context) {
		id, err := strconv.ParseInt(c.Param("id"), 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "id 非法"})
			return
		}
		if !store.Delete(id) {
			c.JSON(http.StatusNotFound, gin.H{"error": "任务不存在"})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	return r
}
