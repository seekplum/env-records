# 环境使用记录管理服务

这是一个使用 Go 语言和 SQLite3 实现的轻量级 Web 服务，用于管理和查询环境（如服务器 IP）的使用占用情况。该服务提供了一套简单的 API，方便团队内部追踪资源的分配状态。

---

## 1. 功能特性

- **环境查询**：实时获取当前所有环境的占用者、标识符及占用时间。
- **批量添加/更新**：支持一次性提交多条记录；若标识符（Ident）已存在，则自动更新现有记录。
- **批量删除**：支持根据标识符列表批量释放环境资源。
- **数据持久化**：使用 SQLite3 本地数据库（`env_records.db`）存储，无需复杂的数据库配置。
- **自定义端口**：支持在启动时通过命令行参数指定监听端口。

---

## 2. 快速开始

### 2.1 环境要求

- 已安装 **Go** (建议 1.18+)。
- **无需 GCC**（已切换至现代纯 Go 驱动 `modernc.org/sqlite`）

### 2.2 安装与运行

```bash
# 初始化项目
go mod tidy

# 运行服务 (默认 9301 端口)
go run main.go -port 9301
```

### 2.3 测试接口

1. 添加记录 (POST)

```bash
curl -X POST http://12.0.0.1:9301/env \
    -H "Content-Type: application/json" \
    -d '[{"ident": "192.168.1.7", "owner": "Tom", "date": "2026-04-29 23:37:35 CST"}, {"ident": "192.168.1.8", "owner": "Jerry", "date": "2026-04-30 10:00:00 CST"}]'
```

2. 查询记录 (GET)

```bash
curl http://12.0.0.1:9301/env
```

3. 删除记录 (DELETE)

```bash
curl -X DELETE http://12.0.0.1:9301/env \
    -H "Content-Type: application/json" \
    -d '{"idents": ["192.168.1.7", "192.168.1.8"]}'
```

## 本地构建

```bash
bash -x bin/build.sh
```

## 执行单元测试

```bash
go test ./...
```
