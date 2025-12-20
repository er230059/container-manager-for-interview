# Container Manager

這是一個用於管理容器的後端服務，使用 Go + Gin 開發。

## 專案功能

- 用戶註冊與認證
- Docker 容器管理
- 檔案上傳

## 系統需求

- Go 1.25
- Docker
- PostgreSQL

## 啟動流程

### 設定環境變數

預設會讀取 `config.yaml` 的設定值，可使用以下環境變數覆蓋

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

### 執行測試

執行單元測試：
```bash
go test ./internal/...
```

執行整合測試：
```bash
go test ./integration_tests/...
```

### 編譯與執行

```bash
make build

./build/main
```

或直接透過 go run 執行：

```bash
go run cmd/main.go
```

## 專案架構

### 目錄結構

```text
.
├── cmd/                # 程式入口點 (main.go)
├── ddl/                # 資料庫定義檔 (SQL 腳本)
├── integration_tests/  # 整合測試
├── internal/           # 內部核心邏輯
│   ├── application/    # 應用層：協調領域模型與處理業務流程
│   ├── domain/         # 領域層：定義實體 (Entity) 與介面 (Interface)
│   ├── errors/         # 錯誤定義
│   ├── infrastructure/ # 基礎設施層：具體實作 (如 Docker API, 資料庫儲存)
│   └── server/         # 路由設定、處理程序 (Handler) 與中間件
|       ├── handler/    
|       └── middleware/
├── pkg/                # 共用模組 (Config Loader, SQL Client)
├── config.yml
├── Makefile
├── go.mod
└── README.md
```

### 非同步處理機制

由於建立 Container 可能耗時較長，該任務會回傳 Job ID 後以非同步的方式建立 Container，範例情境如下:

1. 建立 Container

```bash
curl --location 'http://127.0.0.1:8080/containers' \
--header 'Content-Type: application/json' \
--header 'Authorization: Bearer eyJhb...' \
--data '{
    "cmd": ["tail", "-f", "/dev/null"],
    "env": [],
    "image": "alpine"
}'

{ "job_id": "a8b42d45-b67e-4b77-88b9-a573631a06ee" }
```

2. 使用 Job ID 查詢

```bash
curl --location 'http://127.0.0.1:8080/jobs/a8b42d45-b67e-4b77-88b9-a573631a06ee' \
--header 'Authorization: Bearer eyJhb...'

{ "id":"a8b42d45-b67e-4b77-88b9-a573631a06ee","type":"container_creation","status":"completed","result":{"container_id":"b63595e69fa5377cb565ece4b962118a544e82c0c101e6ccd5c1cb12b79e6f65"},"created_at":"2025-12-20T12:14:09.576918Z","updated_at":"2025-12-20T12:14:11.847488Z" }
```

### 並發控制

對於同一個 container 做啟動、停止、刪除這三個操作時，相同的操作會被合併僅執行一次。例如同時刪除相同的 container 兩次，則系統只會對 Docker 送出一次刪除指令。如果是不同的操作，則只有其一會被執行，另一個 request 會拿到 HTTP 409 Conflict 的錯誤。

例如:

Request A 啟動 Container:
```bash
curl --location -i --request PATCH 'http://127.0.0.1:8080/containers/b63595e69fa5377cb565ece4b962118a544e82c0c101e6ccd5c1cb12b79e6f65/start' \
--header 'Authorization: Bearer eyJhb...'
```
Request B 停止 Container:
```bash
curl --location -i --request PATCH 'http://127.0.0.1:8080/containers/b63595e69fa5377cb565ece4b962118a544e82c0c101e6ccd5c1cb12b79e6f65/stop' \
--header 'Authorization: Bearer eyJhb...'
```

同時送出這兩個 request，則 response 分別會是:

Response A:
```bash
HTTP/1.1 200 OK
Date: Sat, 20 Dec 2025 04:21:35 GMT
Content-Length: 0
```

Response B:
```bash
HTTP/1.1 409 Conflict
Content-Type: application/json; charset=utf-8
Date: Sat, 20 Dec 2025 04:21:35 GMT
Content-Length: 40

{"error":"conflict container operation"}
```

底層實作採用 `singleflight`、`sync.Map` 及 `sync.Mutex` 套件達成，對於相同 container 的相同操作，會使用 `singleflight` 限制同時只會對 Docker API 操作一次，並把該次的結果作為所有相同操作的回傳值，同時對該 container id 加上 mutex lock 避免其他操作。

相關實作位於 `internal/application/container_service.go`
