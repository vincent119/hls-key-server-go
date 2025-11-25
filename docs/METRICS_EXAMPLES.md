# Prometheus 自定義 Metrics 使用範例

## 快速啟動

1. 啟動服務器：

```bash
./bin/server
```

1. 訪問 metrics endpoint（需要 basic auth）：

```bash
curl -u admin:sshhsuuwgwhysgs http://localhost:9090/api/v1/metrics
```

## 可用的自定義 Metrics 範例

### HTTP 請求追蹤

```prometheus
# 總請求數（按方法、路徑和狀態碼分類）
hls_http_requests_total{method="POST",path="/api/v1/hls/key",status="200"} 150
hls_http_requests_total{method="POST",path="/api/v1/auth/token",status="200"} 50
hls_http_requests_total{method="GET",path="/healthz",status="200"} 500

# 請求持續時間（histogram）
hls_http_request_duration_seconds_bucket{method="POST",path="/api/v1/hls/key",le="0.005"} 145
hls_http_request_duration_seconds_sum{method="POST",path="/api/v1/hls/key"} 0.75
hls_http_request_duration_seconds_count{method="POST",path="/api/v1/hls/key"} 150

# 並發連接數
hls_concurrent_connections 5
```

### 密鑰管理 Metrics

```prometheus
# 密鑰請求總數
hls_key_requests_total{key_name="stream.key",status="success"} 145
hls_key_requests_total{key_name="test.key",status="error"} 5

# 緩存命中/未命中
hls_key_cache_hits_total 1200
hls_key_cache_misses_total 50

# 活動密鑰數量
hls_active_keys 3

# 密鑰重載持續時間
hls_key_reload_duration_seconds_bucket{le="0.01"} 10
hls_key_reload_duration_seconds_sum 0.08
hls_key_reload_duration_seconds_count 10

# 密鑰文件大小
hls_key_file_size_bytes{key_name="stream.key"} 16384
```

### 認證 Metrics

```prometheus
# 認證嘗試次數
hls_auth_attempts_total{result="success"} 45
hls_auth_attempts_total{result="invalid_header"} 5
hls_auth_attempts_total{result="invalid_credentials"} 3

# JWT Token 生成總數
hls_token_generations_total 45

# Token 驗證
hls_token_validations_total{result="success"} 145
hls_token_validations_total{result="expired"} 10
hls_token_validations_total{result="invalid"} 5
```

### 錯誤追蹤

```prometheus
# 各類型錯誤統計
hls_errors_total{type="key_retrieval"} 5
hls_errors_total{type="token_generation"} 2
hls_errors_total{type="validation"} 3
```

### 系統 Metrics

```prometheus
# API 版本資訊
hls_api_version_info{version="1.2.3",mode="production"} 1

# 服務器運行時間（秒）
hls_server_uptime_seconds 86400
```

## Prometheus 查詢範例

### 計算請求成功率

```promql
sum(rate(hls_http_requests_total{status="200"}[5m]))
/
sum(rate(hls_http_requests_total[5m])) * 100
```

### 平均請求延遲

```promql
rate(hls_http_request_duration_seconds_sum[5m])
/
rate(hls_http_request_duration_seconds_count[5m])
```

### 錯誤率

```promql
sum(rate(hls_errors_total[5m]))
```

### 認證成功率

```promql
sum(rate(hls_auth_attempts_total{result="success"}[5m]))
/
sum(rate(hls_auth_attempts_total[5m])) * 100
```

### 密鑰請求 QPS

```promql
sum(rate(hls_key_requests_total{status="success"}[1m]))
```

## Grafana 儀表板建議

### 面板建議

1. **HTTP 流量監控**
   - 總請求數時間序列圖
   - 請求延遲熱圖
   - 並發連接數

2. **密鑰管理**
   - 活動密鑰數量
   - 密鑰請求成功/失敗率
   - 緩存命中率

3. **認證監控**
   - 認證嘗試分佈（成功/失敗）
   - Token 生成率
   - Token 驗證成功率

4. **系統健康**
   - 錯誤率趨勢
   - 服務器運行時間
   - 請求成功率

## 告警規則範例

```yaml
groups:
  - name: hls_alerts
    rules:
      # 高錯誤率告警
      - alert: HighErrorRate
        expr: rate(hls_errors_total[5m]) > 10
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High error rate detected"
          description: "Error rate is {{ $value }} errors/sec"

      # 認證失敗率過高
      - alert: HighAuthFailureRate
        expr: |
          sum(rate(hls_auth_attempts_total{result!="success"}[5m]))
          /
          sum(rate(hls_auth_attempts_total[5m])) > 0.1
        for: 5m
        labels:
          severity: warning
        annotations:
          summary: "High authentication failure rate"

      # 密鑰請求失敗率過高
      - alert: HighKeyRequestFailureRate
        expr: |
          sum(rate(hls_key_requests_total{status="error"}[5m]))
          /
          sum(rate(hls_key_requests_total[5m])) > 0.05
        for: 5m
        labels:
          severity: critical
        annotations:
          summary: "High key request failure rate"
```
