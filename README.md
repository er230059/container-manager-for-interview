# Container Manager

這是一個用於管理容器的後端服務，使用 Go + Gin 開發。

## 專案功能

- 用戶註冊與認證
- Docker 容器管理
- 檔案上傳

## 專案目錄結構

```text
.
├── cmd/                # 程式入口點 (main.go)
├── ddl/                # 資料庫定義檔 (SQL 腳本)
├── integration_tests/  # 整合測試
├── internal/           # 內部核心邏輯
│   ├── application/    # 應用層：協調領域模型與處理業務流程
│   ├── domain/         # 領域層：定義實體 (Entity) 與介面 (Interface)
│   ├── errors/         # 錯誤處理定義
│   ├── infrastructure/ # 基礎設施層：具體實作 (如 Docker API, 資料庫儲存)
│   └── server/         # 路由設定、處理程序 (Handler) 與中間件
|       ├── handler/    
|       └── middleware/
├── pkg/                # 共享工具包 (Config, PostgreSQL Client)
├── config.yml
├── Makefile
├── go.mod
└── README.md
```

## 前置需求

- Go 1.25
- Docker
- PostgreSQL

## 快速開始

### 設定環境變數

預設會讀取 `config.yaml` 的設定，該設定可使用以下環境變數覆蓋

| Config Key | 環境變數 | 說明 |
| :--- | :--- | :--- |
| `server.port` | `SERVER_PORT` | 服務監聽埠號 |
| `server.jwt_secret` | `SERVER_JWT_SECRET` | JWT 簽章密鑰 |
| `db.host` | `DB_HOST` | 資料庫主機 |
| `db.port` | `DB_PORT` | 資料庫埠號 |
| `db.user` | `DB_USER` | 資料庫使用者 |
| `db.password` | `DB_PASSWORD` | 資料庫密碼 |
| `db.name` | `DB_NAME` | 資料庫名稱 |
| `storage.base_path` | `STORAGE_BASE_PATH` | 檔案上傳儲存路徑 |

### 初始化資料庫

目錄 `ddl` 下的 SQL 腳本用來建立必要的資料表

### Swagger API 文件

透過以下指令產生 Swagger 文件：

```bash
make gen-doc
```

專案啟動後，您可以透過瀏覽器存取 Swagger UI 來查看與測試 API：

[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

### 編譯與執行

```bash
make build

./build/main
```

或直接透過 go run 執行：

```bash
go run cmd/main.go
```

## 開發與測試

執行單元測試：
```bash
go test ./internal/...
```

執行整合測試：
```bash
go test ./integration_tests/...
```
