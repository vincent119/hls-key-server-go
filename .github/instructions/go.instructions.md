---
description: 'Go 開發與 Copilot/Agent 產生規範指引（整合 Uber Go Style Guide）'
applyTo: '**/*.go,**/go.mod,**/go.sum'
---

# Go 開發與 Copilot / Agent 指南（整合 Uber Guide 版本）

本檔延伸自 `.github/standards/copilot-common.md` 與 `.github/standards/copilot-vocabulary.yaml`，
並 **參照**：

- [Uber Go 風格指南（繁體中文維護版）](https://github.com/ianchen0119/uber_go_guide_tw)
- [Effective Go（官方）](https://go.dev/doc/effective_go#introduction)
- Go 版本目標：**Go 1.22+**（如需更動版本，PR 必述影響）

> 目標：統一格式、用詞與安全實務，確保自動產生或人工撰寫的程式碼皆符合 idiomatic Go，並可直接編譯、部署與維護。

---

## Copilot / Agent 產生守則

### 檔案與 package 規範
- 每檔僅 **一行** `package <name>` 宣告（置頂）。
  - 編輯檔案：保留原 package。
  - 新檔案：與資料夾既有 `.go` 同名 package。
- 可執行程式置於 `cmd/<app>/main.go`，library 不得含 `main()`。
- package 名稱：**全小寫、單字、無底線**（Uber）。

### Imports 與工具
- 產出前必可通過 `gofmt -s`、`goimports`、`go vet`。
- 自動清除未用 imports，避免循環依賴。
- 變更 `go.mod` 後提示 `go mod tidy`。
- 縮排：Tab；檔尾留 **單一** 換行；UTF-8（無 BOM）。
- Imports 排序：**標準庫 → 第三方 → 專案內部**；群組以空行分隔。

### 錯誤處理與流程
- 呼叫後**立即**檢查 `err`，採 **early return**。
- 包裝錯誤：`fmt.Errorf("context: %w", err)`；跨層使用 `errors.Is/As`。
- 訊息小寫開頭，尾端**不加標點**。
- 僅在**不可恢復初始化**時用 `panic`；避免在 library 使用。
- 禁止「只記錄不回傳」導致錯誤吞沒；**記錄與回傳擇一**，以邏輯層級決定。

### 並行與 I/O 安全
- 每個 goroutine 需有**退出機制**（`context`、`WaitGroup` 或關閉 channel）。
- Channel 緩衝預設 0 或 1（除非有量測證據）。
- 嚴禁 goroutine 泄漏；資源關閉要落在呼叫點 `defer Close()`。
- 不可重用已讀取的 `req.Body`；需 **clone**：
  ```go
  // 將來源位元組切片拷貝，確保可重播 Body
  buf := bytes.Clone(src)
  req.Body = io.NopCloser(bytes.NewReader(buf))
  req.GetBody = func() (io.ReadCloser, error) {
      return io.NopCloser(bytes.NewReader(buf)), nil
  }
  ```
- `io.Pipe`/multipart 必須**單執行緒順序寫入**；失敗用 `pw.CloseWithError(err)`、成功 `mw.Close()` 再 `pw.Close()`。
- 底層 slice/map 在**邊界（入/出）**時一律複製，避免別名共享。

### HTTP Client 設計
- `Client` 僅存設定（BaseURL、`*http.Client`、headers）；**不得**保存 `*http.Request` 或可變請求狀態。
- 方法介面：
  - 皆接收 `context.Context`。
  - 內部建 `http.Request` → `c.httpClient.Do(req)` → `defer resp.Body.Close()`。
  - 要求**逾時/重試/回退**策略明確（見「Net/HTTP 實務」）。

### JSON / Struct Tag
- 對外型別欄位加上 `json,yaml,mapstructure` tags；**選填**欄位 `omitempty`。
- 輸入端（decode）預設**拒絕未知欄位**：
  ```go
  dec := json.NewDecoder(r)
  dec.DisallowUnknownFields()
  ```
- 使用 `any` 取代 `interface{}`；但優先具體型別。
- 時間欄位採 **RFC3339**（UTC 優先）；必要時標注本地時區偏差。

### 測試與範例
- 採 **table-driven tests**；子測試用 `t.Run`。
- 輔助函式 `t.Helper()`；清理用 `t.Cleanup()`。
- 匯出 API 提供 `example_test.go`。
- 優先標準 `testing`；除非必要不引入 assert 套件。
- 需通過：`-race`、單元涵蓋率門檻（預設 80% 可調整；變更需 PR 說明）。
- 提供**基準**與**模糊測試（fuzz）**於關鍵路徑。

### 產出內容要求
- 輸出 **完整可編譯檔案**或明確 **diff**。
- 多檔變更列出：檔名 / 變更摘要 / 風險。
- 新增外部套件需附：`go get <module>@<version>` 與風險評估。

### 詞彙與術語
- 優先 `.github/standards/copilot-vocabulary.yaml`。
- 與現有命名衝突時以 vocabulary 為準。
- 與 Uber/Effective Go 不一致時，PR **必述理由**與替代方案。

---

## Go 一般開發規範（整合 Uber + Effective Go）

### 通用原則
- 清晰優於巧妙；主流程靠左排列；讓 **零值可用**。
- 結構自我說明；註解描述「**為何**」而非「做什麼」。

### 命名慣例
- package：全小寫、單字、無底線；避免 `util`、`common`。
- 變數/函式：小駝峰；匯出名稱首字母大寫。
- 介面以 `-er` 結尾（Reader/Writer）；**小介面**優先。
- 縮略詞大小寫一致：`HTTPServer`、`URLParser`。
- 建構子命名採 `NewType(...)`；必要時 `WithXxx` 選項，但避免過度抽象。

### 常數與列舉
- 群組 `const (...)`；**型別化常數**避免魔數。
- Enum 起始值**考慮零值可用性**，必要時保留 `Unknown`。

### 接收者與方法
- 以量測決定**指標/值**接收者（大型結構/需變異 → 指標；小值/不變 → 值）。
- 避免 `init()` 副作用與全域可變狀態。

### `context` 規範
- 對外 API **第一個參數**為 `ctx context.Context`。
- 禁用 `context.Background()` 直傳至深層；由呼叫者注入。
- 設定逾時/截止於**呼叫邊界**；尊重 `ctx.Done()`。
- 不將 `ctx` 保存於結構體。

### 並行進階
- 以 `errgroup`/`WaitGroup` + `ctx` 收斂；提供**背壓**與**取消**。
- 共享狀態以 `sync.Mutex/RWMutex` 或無鎖結構（經量測）保護。
- 在 `for range` 捕獲變數時，**重新宣告** loop 變數避免閉包陷阱。

### Net/HTTP 實務
- **重用 Transport**，設定逾時：
  ```go
  tr := &http.Transport{
      MaxIdleConns:        100,
      IdleConnTimeout:     90 * time.Second,
      TLSHandshakeTimeout: 10 * time.Second,
      ExpectContinueTimeout: 1 * time.Second,
  }
  c := &http.Client{
      Transport: tr,
      Timeout:   15 * time.Second, // 全域上限；更細粒度以 context 控制
  }
  ```
- 明確重試策略（**僅**冪等方法），具退避與上限；對 5xx/網路錯誤重試，對業務 4xx 不重試。
- 嚴格 `resp.Body.Close()`；讀取前先檢 HTTP 狀態碼。

### 日誌與可觀測性
- 使用**結構化日誌**（如 `zap`）；固定欄位：`trace_id`, `span_id`, `req_id`, `subsystem`。
- `logger.Error("msg", zap.Error(err))` 報告；避免把錯誤訊息再字串化拼接。
- 指標/追蹤採 OpenTelemetry；HTTP/DB 客戶端優先用已 instrument 的實作。

### 時間與時區
- **內部以 UTC 儲存與運算**；輸出呈現再格式化。
- JSON 時間使用 RFC3339（必要時 `time.RFC3339Nano`）。

### 安全性
- 僅用標準 `crypto/*`；禁自製密碼學。
- 外部輸入需驗證與長度限制；避免正則 ReDoS。
- 檔案 I/O 使用 `fs.FS` 與限制型讀取；防 Zip Slip。
- 納入 `gosec`（或等價 analyzer）於 CI；敏感資訊不得進日誌。

### 依賴與模組
- 模組遵循 **SemVer**；破壞性改動於 major path（`/v2`）。
- 嚴格釘版：`go.mod` 使用最小相依原則；避免 transitive 泄漏。
- 移除依賴需跑 `go mod tidy` 並附影響說明。

### 目錄結構建議
```
<app>/cmd/<app>/main.go    # 應用進入點；僅進行依賴注入、初始化與啟動，不含業務邏輯
<app>/internal/...         # 僅供本模組使用的封裝邏輯（封裝 service、repository、adapter 等，外部無法 import）
<app>/pkg/...              # 穩定且可重用的公開 API（供其他模組或工具匯入使用）
<app>/api/...              # 對外契約層（OpenAPI、protobuf、GraphQL schema）；生成的 stub/server/client 亦放此層
<app>/configs/...          # 預設設定、環境樣板與設定載入程式（config.yaml, env.example 等）
<app>/scripts/...          # 開發、測試、部署、CI/CD 用腳本（bash、make、python 等）
Dockerfile                 # 應用容器化描述；建議採 multi-stage，含 lint / test / build 階段
.dockerignore              # Docker build 忽略清單（排除編譯輸出、暫存檔與測試資料）
.gitignore                 # Git 忽略清單（node_modules、vendor、log、tmp 等）
Makefile                   # 常用開發指令（tidy、lint、test、bench、build、run 等封裝）
.golangci.yml              # 靜態分析與 Linter 設定（統一風格與品質門檻）
README.md                  # 專案說明：目的、架構、建置步驟、測試與部署指引
LICENSE                    # 授權條款；內部專案可標註版權與使用限制
go.mod                     # Go 模組定義與依賴版本管理
go.sum                     # 依賴模組驗證雜湊清單（確保可重建性）
```

### 產生器與 build
- 使用 `//go:build` 標籤管理條件編譯；禁止舊 `+build` 註解。
- `go generate` 指令須在檔頭註解，並可重複執行（可重入）。
- CGO 預設關閉；開啟需 PR 說明平台/效能/部署影響。

---

## Copilot / Agent 提示模板（Do/Don't）

**Do**
- 僅產生一個 `package` 宣告；imports 分群。
- 立即檢查 `err`，使用 `%w` 包裝。
- 所有公開 API 第一參數 `context.Context`。
- 在邊界複製 slice/map；為 struct 加 `json`/`yaml` tags。
- 撰寫 table-driven 測試 + `t.Helper()`，並加入一個基準測試。

**Don’t**
- 不要保存 `context.Context` 或 `*http.Request` 於 struct。
- 不要在 library 使用 `panic`；不要忽略 `Close()`。
- 不要以 `interface{}` 取代具體型別；不要暴露可變 slice/map。
- 不要在 loop 內 `defer`（除非小範圍且有量測）。

---

## 實作範例片段

> 具體展示關鍵規則如何落地。

```go
// Package client 提供與遠端服務互動的 HTTP 用戶端。
// 零值不可用，請使用 New 建構子建立。
package client

import (
	"context"          // 佈線取消與逾時的標準機制
	"encoding/json"    // 編解碼輸入/輸出
	"errors"           // 錯誤比對
	"fmt"              // 錯誤包裝與格式化
	"io"               // I/O 介面
	"net/http"         // HTTP 基礎
	"time"             // 逾時與回退間隔
)

// ErrNotFound：對應遠端 404 的語意錯誤（sentinel error）。
var ErrNotFound = errors.New("resource not found") // 小寫開頭，不加標點

// Client 僅保存不可變設定與共享 *http.Client；不保存請求狀態。
type Client struct {
	baseURL    string        // 基底位址（不可含尾斜線）
	httpClient *http.Client  // 可注入以便測試與重用 Transport
}

// New 建立可用的 Client；呼叫者可注入自定 *http.Client。
func New(baseURL string, hc *http.Client) *Client {
	if hc == nil {
		hc = &http.Client{Timeout: 15 * time.Second} // 安全預設
	}
	return &Client{baseURL: baseURL, httpClient: hc}
}

// Resource 對外輸出時含有 json 標籤，零值可用。
type Resource struct {
	ID        string    `json:"id"`
	Name      string    `json:"name,omitempty"`
	UpdatedAt time.Time `json:"updatedAt"` // RFC3339 UTC
}

// Get 透過 context 控制逾時/取消，正確關閉 Body 並轉換語意錯誤。
func (c *Client) Get(ctx context.Context, id string) (Resource, error) {
	var out Resource

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, c.baseURL+"/v1/resources/"+id, nil)
	if err != nil {
		return out, fmt.Errorf("new request: %w", err)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return out, fmt.Errorf("do request: %w", err)
	}
	defer resp.Body.Close() // 確保釋放連線

	switch resp.StatusCode {
	case http.StatusOK:
		dec := json.NewDecoder(resp.Body)
		dec.DisallowUnknownFields()
		if err := dec.Decode(&out); err != nil {
			return out, fmt.Errorf("decode: %w", err)
		}
		return out, nil
	case http.StatusNotFound:
		// 將 HTTP 狀態轉換為語意錯誤
		io.Copy(io.Discard, resp.Body) // 盡量讀完以便連線重用
		return out, ErrNotFound
	default:
		b, _ := io.ReadAll(io.LimitReader(resp.Body, 64<<10)) // 限流保護
		return out, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(b))
	}
}
```

---

## Review Checklist
- [ ] 僅一個 `package` 宣告（置頂）
- [ ] 通過 `gofmt -s` / `goimports` / `go vet`
- [ ] 無未使用 imports、無循環依賴
- [ ] `err` 立即檢查並以 `%w` 包裝；跨層以 `errors.Is/As`
- [ ] goroutine / channel 正確收斂；無泄漏
- [ ] I/O 操作安全（含 Close、Pipe、Body 重新可讀）
- [ ] JSON tag 一致、解碼拒絕未知欄位、零值可用
- [ ] 測試含 table-driven、-race、必要 fuzz/bench
- [ ] `go.mod` 依賴釘版；`go mod tidy` 後無不明變更
- [ ] 與 Uber / Effective Go 一致或於 PR 註明偏離理由

---

## CI 與工具建議（可直接採用）

### `Makefile`（節選）
```makefile
.PHONY: tidy lint test bench

tidy:
	go mod tidy

lint:
	golangci-lint run ./...

test:
	go test -race -count=1 ./...

bench:
	go test -run=NONE -bench=. -benchmem ./...
```

### `.golangci.yml`（節選）
```yaml
run:
  timeout: 5m
linters:
  enable:
    - errcheck
    - gocritic
    - gofumpt
    - govet
    - ineffassign
    - staticcheck
    - unparam
    - prealloc
issues:
  exclude-use-default: false
```

### PR 模板要點
- 目的與背景（為何要改）
- 變更摘要（做了什麼）
- 風險與回滾方案
- 測試證據（覆蓋率、基準、相容性）
- 偏離 Uber/Effective Go 的理由（若有）

---

**建議存放路徑：** `.github/instructions/go.instructions.md`
此設定將自動套用於所有 Go 檔案（`*.go`, `go.mod`, `go.sum`）。
