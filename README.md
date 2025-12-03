# Secure Notes Project

## Mô tả
Dự án ghi chú mã hóa, chia sẻ bảo mật, hỗ trợ nhiều client (CLI) và backend API server.

## Cấu trúc thư mục
```
project-root/
├── client/      # CLI client (Go, tách biệt go.mod)
├── server/      # Backend API server (Go, tách biệt go.mod)
├── docs/        # Tài liệu, sơ đồ, API spec
├── test/        # Test chung (nếu có)
├── Makefile     # Lệnh build/test nhanh (nếu dùng)
└── README.md    # File này
```

## Hướng dẫn build & chạy

### Server
Xem chi tiết ở `server/README.md`.
```bash
cd server
go mod tidy
# Chạy migration DB (Postgres)
# ...
go run cmd/main.go
```

### Client CLI
Xem chi tiết ở `client/README.md`.
```bash
cd client
go mod tidy
go build -o notescli cmd/main.go
./notescli --help
```

## Tài liệu
- Sơ đồ, API, migration: xem trong `docs/`
- Hướng dẫn chi tiết từng phần: xem `client/README.md`, `server/README.md`
