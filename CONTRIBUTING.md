# Contributing to HLS Key Server

首先，感謝您考慮為 HLS Key Server 做出貢獻！我們歡迎所有形式的貢獻。

## 目錄

- [行為準則](#行為準則)
- [如何貢獻](#如何貢獻)
- [開發流程](#開發流程)
- [程式碼規範](#程式碼規範)
- [提交訊息規範](#提交訊息規範)
- [測試要求](#測試要求)
- [Pull Request 流程](#pull-request-流程)

## 行為準則

本專案遵循貢獻者公約（Contributor Covenant）。參與本專案即表示您同意遵守其條款。

### 我們的承諾

為了營造開放且友善的環境，我們承諾：

- 使用友善且包容的語言
- 尊重不同的觀點和經驗
- 優雅地接受建設性批評
- 關注對社群最有利的事情
- 對其他社群成員表現同理心

## 如何貢獻

### 回報 Bug

如果您發現 bug，請：

1. **檢查是否已有相同問題** - 搜尋 [Issues](https://github.com/vincent119/hls-key-server-go/issues)
2. **建立新 Issue** - 如果問題尚未回報，請建立新的 issue
3. **提供詳細資訊**：
   - 清楚的標題和描述
   - 重現步驟
   - 預期行為 vs 實際行為
   - Go 版本、作業系統
   - 相關日誌或截圖

**Bug Report 範本**：

```markdown
## 描述
簡短描述問題

## 重現步驟
1. 步驟一
2. 步驟二
3. 觀察到的錯誤

## 預期行為
應該發生什麼

## 實際行為
實際發生什麼

## 環境
- Go 版本：1.24.0
- 作業系統：macOS 14.0
- 專案版本：v1.0.0

## 日誌
```
相關錯誤日誌
```
```

### 提議新功能

如果您想提議新功能：

1. **開啟 Discussion** - 先在 [Discussions](https://github.com/vincent119/hls-key-server-go/discussions) 討論
2. **說明用例** - 解釋為什麼需要這個功能
3. **考慮替代方案** - 是否有其他實作方式
4. **等待反饋** - 獲得維護者認可後再開始開發

**Feature Request 範本**：

```markdown
## 功能描述
簡短描述您想要的功能

## 動機
為什麼需要這個功能？解決什麼問題？

## 建議實作
您認為應該如何實作

## 替代方案
考慮過哪些其他方案

## 額外資訊
其他相關資訊或截圖
```

### 改善文件

文件貢獻同樣重要！

- 修正錯字或文法錯誤
- 改善現有文件的清晰度
- 新增遺漏的文件
- 翻譯文件（未來可能）

## 開發流程

### 1. Fork 並 Clone

```bash
# Fork 專案到您的 GitHub 帳號
# 然後 clone 您的 fork

git clone https://github.com/YOUR_USERNAME/hls-key-server-go.git
cd hls-key-server-go

# 新增上游遠端
git remote add upstream https://github.com/vincent119/hls-key-server-go.git
```

### 2. 建立分支

```bash
# 確保主分支是最新的
git checkout main
git pull upstream main

# 建立功能分支
git checkout -b feature/amazing-feature
# 或修復分支
git checkout -b fix/bug-description
```

分支命名規範：

- `feature/功能名稱` - 新功能
- `fix/問題描述` - Bug 修復
- `docs/文件主題` - 文件更新
- `refactor/重構範圍` - 程式碼重構
- `test/測試範圍` - 測試相關
- `chore/任務描述` - 雜項任務

### 3. 設定開發環境

```bash
# 安裝依賴
go mod download

# 執行測試確保環境正常
make test

# 執行 linter
make lint
```

### 4. 進行變更

開發時請遵循：

- [程式碼規範](#程式碼規範)
- [Uber Go Style Guide](https://github.com/uber-go/guide/blob/master/style.md)
- 專案現有的架構模式

### 5. 測試

```bash
# 執行所有測試
make test

# 執行特定套件測試
go test -v ./internal/service/...

# 執行 race detector
go test -race ./...

# 查看覆蓋率
go test -cover ./...
```

### 6. Commit 變更

```bash
# 暫存變更
git add .

# 提交（遵循 commit 訊息規範）
git commit -m "feat: add amazing feature"
```

### 7. 推送並建立 Pull Request

```bash
# 推送到您的 fork
git push origin feature/amazing-feature

# 在 GitHub 上建立 Pull Request
```

## 程式碼規範

### Go 編碼標準

遵循 **Uber Go Style Guide** 和 **Effective Go**：

#### 1. 命名規範

```go
// ✅ Good - 使用駝峰式命名
type UserService struct {
    repository UserRepository
    logger     *zap.Logger
}

// ❌ Bad - 使用下劃線
type user_service struct {
    user_repository UserRepository
}

// ✅ Good - 介面命名（-er 後綴）
type Reader interface {
    Read(p []byte) (n int, err error)
}

// ✅ Good - 縮寫全大寫或全小寫
userID := "123"
httpServer := &http.Server{}
```

#### 2. 錯誤處理

```go
// ✅ Good - 使用 sentinel errors
var ErrNotFound = errors.New("not found")

// ✅ Good - 錯誤包裝保留堆疊
if err != nil {
    return fmt.Errorf("failed to load config: %w", err)
}

// ✅ Good - 使用 errors.Is/As
if errors.Is(err, ErrNotFound) {
    // handle not found
}
```

#### 3. Context 使用

```go
// ✅ Good - context 作為第一個參數
func (s *Service) GetUser(ctx context.Context, id string) (*User, error) {
    // implementation
}

// ❌ Bad - context 不是第一個參數
func (s *Service) GetUser(id string, ctx context.Context) (*User, error) {
```

#### 4. 依賴注入

```go
// ✅ Good - 建構子注入
func NewService(repo Repository, logger *zap.Logger) *Service {
    return &Service{
        repo:   repo,
        logger: logger,
    }
}

// ❌ Bad - 使用 init() 或全域變數
var globalLogger *zap.Logger

func init() {
    globalLogger = zap.NewProduction()
}
```

#### 5. 介面定義

```go
// ✅ Good - 在使用方定義小介面
type KeyRepository interface {
    Get(ctx context.Context, name string) ([]byte, error)
    List(ctx context.Context) []string
}

// ❌ Bad - 大而全的介面
type Repository interface {
    Get(...)
    List(...)
    Create(...)
    Update(...)
    Delete(...)
    // 10+ more methods
}
```

### 格式化

```bash
# 使用 gofmt
make fmt

# 或
go fmt ./...

# 使用 goimports（推薦）
goimports -w .
```

### Linting

```bash
# 執行 golangci-lint
make lint

# 自動修復部分問題
golangci-lint run --fix
```

必須通過的檢查：

- `govet` - Go 官方靜態分析
- `errcheck` - 錯誤檢查
- `staticcheck` - 進階靜態分析
- `gosimple` - 簡化建議
- `unused` - 未使用程式碼檢測

## 提交訊息規範

遵循 **Conventional Commits** 規範：

### 格式

```
<type>(<scope>): <subject>

<body>

<footer>
```

### Type 類型

- `feat`: 新功能
- `fix`: Bug 修復
- `docs`: 文件變更
- `style`: 格式化（不影響程式碼功能）
- `refactor`: 重構（既非新功能也非 bug 修復）
- `perf`: 效能改善
- `test`: 測試相關
- `chore`: 建置或輔助工具變更
- `ci`: CI/CD 相關

### Scope 範圍

- `handler`: Handler 層
- `service`: Service 層
- `repository`: Repository 層
- `config`: 配置相關
- `auth`: 認證相關
- `api`: API 相關

### 範例

```bash
# 新功能
git commit -m "feat(service): add key caching mechanism"

# Bug 修復
git commit -m "fix(handler): prevent path traversal in key name"

# 文件
git commit -m "docs(readme): update API usage examples"

# 重構
git commit -m "refactor(repository): simplify key loading logic"

# 效能
git commit -m "perf(service): optimize JWT validation"

# Breaking change
git commit -m "feat(api)!: change authentication endpoint

BREAKING CHANGE: /auth endpoint moved to /api/v1/auth"
```

## 測試要求

### 必須條件

所有 Pull Request 必須：

1. ✅ **通過所有現有測試**
2. ✅ **新程式碼有對應測試**
3. ✅ **測試覆蓋率不降低**
4. ✅ **通過 race detector**

### 測試類型

#### 單元測試

```go
// 測試檔案命名: *_test.go
func TestNewService(t *testing.T) {
    logger := zap.NewNop()
    repo := &mockRepository{}

    service := NewService(repo, logger)

    if service == nil {
        t.Fatal("expected non-nil service")
    }
}
```

#### Table-Driven Tests（推薦）

```go
func TestService_GetKey(t *testing.T) {
    tests := []struct {
        name        string
        keyName     string
        setupMock   func(*mockRepository)
        expectedErr error
    }{
        {
            name:    "existing key",
            keyName: "test.key",
            setupMock: func(m *mockRepository) {
                m.keys["test.key"] = []byte("data")
            },
            expectedErr: nil,
        },
        {
            name:        "non-existent key",
            keyName:     "missing.key",
            setupMock:   func(m *mockRepository) {},
            expectedErr: ErrKeyNotFound,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mock := &mockRepository{keys: make(map[string][]byte)}
            tt.setupMock(mock)
            service := NewService(mock, zap.NewNop())

            _, err := service.GetKey(context.Background(), tt.keyName)

            if !errors.Is(err, tt.expectedErr) {
                t.Errorf("expected error %v, got %v", tt.expectedErr, err)
            }
        })
    }
}
```

#### HTTP 整合測試

```go
func TestHandler_GetKey(t *testing.T) {
    gin.SetMode(gin.TestMode)
    router := setupTestRouter()

    req := httptest.NewRequest("POST", "/api/v1/hls/key", nil)
    w := httptest.NewRecorder()

    router.ServeHTTP(w, req)

    if w.Code != http.StatusOK {
        t.Errorf("expected status 200, got %d", w.Code)
    }
}
```

### 執行測試

```bash
# 所有測試
make test

# 特定套件
go test -v ./internal/service/

# 覆蓋率報告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Race detector
go test -race ./...

# Benchmark
go test -bench=. -benchmem ./...
```

## Pull Request 流程

### 提交前檢查清單

- [ ] 程式碼遵循專案風格指南
- [ ] 已執行 `make lint` 且無錯誤
- [ ] 已執行 `make test` 且全部通過
- [ ] 新增或更新了相關測試
- [ ] 新增或更新了相關文件
- [ ] Commit 訊息遵循規範
- [ ] 已從 main 分支 rebase 最新程式碼

### PR 描述範本

```markdown
## 變更類型
- [ ] Bug fix
- [ ] New feature
- [ ] Breaking change
- [ ] Documentation update

## 描述
簡短描述此 PR 的變更內容

## 相關 Issue
Fixes #123

## 變更內容
- 變更項目 1
- 變更項目 2

## 測試
描述您如何測試此變更

## 截圖（如適用）
新增相關截圖

## 檢查清單
- [ ] 我的程式碼遵循專案風格指南
- [ ] 我已執行自我 review
- [ ] 我已新增必要的註解
- [ ] 我已更新文件
- [ ] 我的變更不會產生新的警告
- [ ] 我已新增測試證明修復有效或功能運作
- [ ] 新舊單元測試都在本地通過
```

### Review 流程

1. **自動檢查** - CI/CD 會自動執行測試和 linter
2. **程式碼審查** - 維護者會審查您的程式碼
3. **討論** - 如有需要，會在 PR 中討論
4. **修改** - 根據反饋進行修改
5. **合併** - 通過審查後合併到 main

### 回應 Review 意見

- 對建設性批評保持開放態度
- 清楚解釋您的設計決策
- 如果不同意，禮貌地討論
- 及時回應 review 意見
- 完成修改後標記為已解決

## 版本發布

版本號遵循 **Semantic Versioning (SemVer)**：

- `MAJOR.MINOR.PATCH`
- `MAJOR`: 不相容的 API 變更
- `MINOR`: 向後相容的新功能
- `PATCH`: 向後相容的 bug 修復

範例：`v1.2.3`

## 獲得幫助

如果您需要幫助：

1. **查看文件** - README.md、ARCHITECTURE.md
2. **搜尋 Issues** - 看看是否有人遇到相同問題
3. **開啟 Discussion** - 在 Discussions 提問
4. **聯絡維護者** - 透過 Issue 或 email

## 致謝

感謝所有貢獻者！您的貢獻讓這個專案更好。

貢獻者名單：[CONTRIBUTORS.md](CONTRIBUTORS.md)（未來建立）

## 授權

貢獻到本專案即表示您同意您的貢獻將依照 [MIT License](LICENSE) 授權。

---

**再次感謝您的貢獻！** 🎉

有任何問題，歡迎隨時開啟 Issue 或 Discussion。
