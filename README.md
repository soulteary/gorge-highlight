# gorge-highlight

Go 语法高亮微服务，为 Phorge 提供代码高亮能力。使用 [Chroma](https://github.com/alecthomas/chroma) 引擎替代 Pygments 子进程调用，通过 HTTP API 输出与 Pygments CSS 类名兼容的 HTML，无需修改 Phorge 前端样式即可无缝切换。

## 特性

- 基于 Chroma 引擎，支持 200+ 种编程语言
- 输出与 Pygments CSS 类名完全兼容（`k`, `nf`, `nb`, `s2` 等），复用 Phorge 已有样式表
- 100+ 条语言别名映射，与 Phorge PHP 端 `PhutilPygmentsSyntaxHighlighter` 保持一致
- 三级语言检测：显式指定 -> 内容分析 -> 兜底降级，保证永不失败
- 可选 Token 认证，支持请求头和查询参数两种方式
- 可配置源码大小限制，防止大文件攻击
- 静态编译，零外部依赖，Docker 镜像极轻量
- 内置健康检查端点，适配容器编排

## 快速开始

### 本地运行

```bash
go build -o gorge-highlight ./cmd/server
./gorge-highlight
```

服务默认监听 `:8140`。

### Docker 运行

```bash
docker build -t gorge-highlight .
docker run -p 8140:8140 gorge-highlight
```

### 带配置运行

```bash
export SERVICE_TOKEN="my-secret-token"
export MAX_BYTES=2097152
./gorge-highlight
```

## 配置

支持两种配置方式：环境变量（优先）和 JSON 配置文件。

### 环境变量

| 变量 | 默认值 | 说明 |
|---|---|---|
| `LISTEN_ADDR` | `:8140` | 服务监听地址 |
| `SERVICE_TOKEN` | (空) | 服务间认证 Token，为空则不启用认证 |
| `MAX_BYTES` | `1048576` (1 MiB) | 单次请求源码最大字节数 |
| `TIMEOUT_SEC` | `15` | 请求超时秒数 |
| `HIGHLIGHT_CONFIG_FILE` | (无) | JSON 配置文件路径，设置后从文件加载配置 |

### JSON 配置文件

设置 `HIGHLIGHT_CONFIG_FILE` 环境变量指向一个 JSON 文件：

```json
{
  "listenAddr": ":8140",
  "serviceToken": "my-token",
  "maxBytes": 1048576,
  "timeoutSec": 15
}
```

## API

所有 `/api/highlight/*` 端点在启用 `SERVICE_TOKEN` 时需要认证。认证方式：

- 请求头：`X-Service-Token: <token>`
- 查询参数：`?token=<token>`

### POST /api/highlight/render

高亮渲染源代码。

**请求体**：

```json
{
  "source": "print('hello')",
  "language": "python"
}
```

- `source`：待高亮的源代码文本
- `language`：语言标识符（可选，留空则自动检测）

**成功响应** (200)：

```json
{
  "data": {
    "html": "<span class=\"nb\">print</span><span class=\"p\">(</span><span class=\"s1\">&#39;hello&#39;</span><span class=\"p\">)</span>",
    "language": "python"
  }
}
```

**错误响应**：

| 状态码 | 错误码 | 说明 |
|---|---|---|
| 400 | `ERR_BAD_REQUEST` | 请求体格式错误 |
| 401 | `ERR_UNAUTHORIZED` | Token 缺失或无效 |
| 413 | `ERR_TOO_LARGE` | 源码超过 MaxBytes 限制 |
| 500 | `ERR_HIGHLIGHT_FAILED` | 高亮处理内部错误 |

### GET /api/highlight/languages

返回所有支持的语言列表。

**响应** (200)：

```json
{
  "data": ["abap", "actionscript", "ada", "..."]
}
```

### GET /healthz

健康检查端点，不需要认证。

**响应** (200)：

```json
{"status": "ok"}
```

## 项目结构

```
gorge-highlight/
├── cmd/server/main.go              # 服务入口
├── internal/
│   ├── config/config.go            # 配置加载（环境变量 / JSON 文件）
│   ├── highlight/highlight.go      # Chroma 高亮引擎核心逻辑
│   └── httpapi/handlers.go         # HTTP 路由、中间件与处理器
├── Dockerfile                      # 多阶段 Docker 构建
├── go.mod
└── go.sum
```

## 开发

```bash
# 运行全部测试
go test ./...

# 运行测试（带详细输出）
go test -v ./...

# 构建二进制
go build -o gorge-highlight ./cmd/server
```

## 技术栈

- **语言**：Go 1.26
- **语法高亮**：[Chroma](https://github.com/alecthomas/chroma) v2.14.0
- **HTTP 框架**：[Echo](https://echo.labstack.com/) v4.12.0
- **许可证**：Apache License 2.0
