# Dango 🥟

一个简单的 Go Web 示例项目，演示使用 Go 标准库 `net/http` 构建 Web 服务。

## 功能

- 首页 HTML 页面
- `/health` 健康检查接口
- `/api/hello` JSON 问候接口
- `/api/time` 返回当前服务器时间

## 运行

```bash
go run .
```

默认监听 `:8080`，可通过环境变量 `PORT` 修改端口：

```bash
PORT=3000 go run .
```

## 构建

```bash
go build -o dango .
./dango
```

## 接口示例

```bash
curl http://localhost:8080/health
# {"status":"ok"}

curl http://localhost:8080/api/hello
# {"message":"Hello from Dango! 🥟","timestamp":"..."}

curl http://localhost:8080/api/time
# {"time":"..."}
```

## 许可

MIT
