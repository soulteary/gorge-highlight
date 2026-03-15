# go-highlight

Go syntax highlighting service for Phorge, replacing the Pygments subprocess with a Chroma-based HTTP API.

Produces HTML with Pygments-compatible CSS class names, so existing Phorge stylesheets work without modification.

## Configuration

### Environment Variables

| Variable | Default | Description |
|---|---|---|
| `LISTEN_ADDR` | `:8140` | Server listen address |
| `SERVICE_TOKEN` | (empty) | Token for `X-Service-Token` auth |
| `MAX_BYTES` | `1048576` | Maximum source code size (bytes) |
| `TIMEOUT_SEC` | `15` | Request timeout |
| `HIGHLIGHT_CONFIG_FILE` | (none) | Path to JSON config file |

### Config File (optional)

```json
{
  "listenAddr": ":8140",
  "serviceToken": "my-token",
  "maxBytes": 1048576,
  "timeoutSec": 15
}
```

## API

### POST /api/highlight/render

Highlight source code.

**Request**:
```json
{
  "source": "print('hello')",
  "language": "python"
}
```

**Response**:
```json
{
  "data": {
    "html": "<span class=\"nb\">print</span>...",
    "language": "python"
  }
}
```

### GET /api/highlight/languages

Returns list of supported language names.

### GET /healthz

Health check endpoint.

## Development

```bash
go test ./...
go build ./cmd/server
```

## Docker

```bash
docker build -t github.com/soulteary/gorge-highlight .
docker run -p 8140:8140 github.com/soulteary/gorge-highlight
```
