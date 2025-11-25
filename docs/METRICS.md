# Prometheus Metrics Endpoint

## Overview

The server now exposes a `/metrics` endpoint for Prometheus to scrape metrics. This endpoint is protected with HTTP Basic Authentication.

## Configuration

The metrics endpoint authentication is configured in `config/config.yaml`:

```yaml
metric:
  user: "admin"
  password: "sshhsuuwgwhysgs"
```

## Access

The metrics endpoint is available at:

```text
GET /api/v1/metrics
```

### Authentication

The endpoint requires HTTP Basic Authentication. Example using curl:

```bash
curl -u admin:sshhsuuwgwhysgs http://localhost:9090/api/v1/metrics
```

### Prometheus Configuration

Add the following to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'hls-key-server'
    basic_auth:
      username: 'admin'
      password: 'sshhsuuwgwhysgs'
    static_configs:
      - targets: ['localhost:9090']
    metrics_path: '/api/v1/metrics'
```

## Security

- Uses constant-time comparison to prevent timing attacks
- Credentials are read from configuration file
- Failed authentication attempts are logged
- Returns `401 Unauthorized` with `WWW-Authenticate` header on auth failure

## Implementation Details

- **Handler**: `internal/handler/metrics_handler.go`
- **Route**: `internal/routes/api/v1/metrics.go`
- **Tests**: `internal/handler/metrics_handler_test.go`
- **Custom Metrics**: `internal/pkg/metrics/metrics.go`
- **Middleware**: `internal/handler/middleware/prometheus.go`
- **Dependencies**: Uses `github.com/prometheus/client_golang/prometheus/promhttp`

## Available Custom Metrics

### HTTP Metrics

- `hls_http_requests_total` - Total HTTP requests by method, path, and status
- `hls_http_request_duration_seconds` - HTTP request duration histogram
- `hls_concurrent_connections` - Current number of concurrent connections

### Key Management Metrics

- `hls_key_requests_total` - Total key requests by key name and status
- `hls_key_cache_hits_total` - Total number of cache hits
- `hls_key_cache_misses_total` - Total number of cache misses
- `hls_active_keys` - Number of currently active keys
- `hls_key_reload_duration_seconds` - Duration of key reload operations
- `hls_key_file_size_bytes` - Size of key files in bytes

### Authentication Metrics

- `hls_auth_attempts_total` - Total authentication attempts by result (success/invalid_header/invalid_credentials)
- `hls_token_generations_total` - Total JWT tokens generated
- `hls_token_validations_total` - Total JWT token validations by result

### Error Metrics

- `hls_errors_total` - Total errors by type

### System Metrics

- `hls_server_uptime_seconds` - Server uptime counter
- `hls_api_version_info` - API version and mode information


