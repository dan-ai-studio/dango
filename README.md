# Dango 🥟

Gin + Vue 3 的完整任务管理 Demo：后端 Gin 提供 REST API，前端 Vue 3 单页应用（CDN 引入，无需构建），首页和接口由同一个 Go 二进制同时提供。

## 功能

- 任务 CRUD：新增、列表（搜索 `?q=` + 按状态过滤 `?done=`）、重命名、删除
- 一键切换完成状态 `PATCH /api/tasks/:id/toggle`
- 统计 `GET /api/stats`（全部 / 待办 / 已完成）
- 健康检查 `GET /health`、`GET /api/health`
- 数据存内存（重启重置，Vercel 上每次冷启动也会重置，仅适合演示）

## 目录结构

```
main.go                    本地/VPS 入口（监听 $PORT）
api/index.go               Vercel Serverless 入口（复用同一 Gin 引擎）
internal/server/server.go  路由 + 中间件（CORS/日志/恢复）
internal/server/store.go   内存任务存储
internal/server/web/index.html  Vue 3 前端（embed 打进二进制）
vercel.json                Vercel rewrite：所有请求转给 /api/index
```

## 本地运行

```bash
go run .
# 打开 http://localhost:8080
```

指定端口：

```bash
PORT=3000 go run .
```

## 接口示例

```bash
curl http://localhost:8080/api/health
# {"status":"ok","service":"dango"}

curl http://localhost:8080/api/tasks
# {"items":[...]}

curl -X POST http://localhost:8080/api/tasks \
  -H "Content-Type: application/json" \
  -d '{"title":"写个新任务"}'

curl -X PATCH http://localhost:8080/api/tasks/1/toggle

curl -X PUT http://localhost:8080/api/tasks/1 \
  -H "Content-Type: application/json" \
  -d '{"title":"新标题","done":true}'

curl -X DELETE http://localhost:8080/api/tasks/1

curl "http://localhost:8080/api/tasks?q=gin&done=false"
curl http://localhost:8080/api/stats
```

## 部署到 Vercel

Application Preset 选 `Go`，三个命令保持默认（Build: None，Output: N/A，Install: `go mod download`），直接 Deploy 即可。`vercel.json` 已配好全量转发。

注意：Vercel 是 Serverless，内存数据不持久化，仅适合演示；要持久化请接数据库（Postgres/Redis）或部署到 VPS（`go build -o dango .` + systemd 常驻）。

## 许可

MIT
