# Secure Notes Server

## Mô tả
Backend API server cho ứng dụng ghi chú mã hóa, hỗ trợ xác thực, quản lý ghi chú, chia sẻ bảo mật, và quản lý token.

## Cách chạy
```bash
cd server
# Cài dependencies
go mod tidy
# Chạy migration (yêu cầu PostgreSQL)
# Ví dụ với golang-migrate:
migrate -path migrations -database "postgres://user:pass@localhost:5432/notesdb?sslmode=disable" up
# Chạy server
cd cmd
go run main.go
```

## Cấu trúc thư mục
- `cmd/`         : Entrypoint server
- `internal/`    : Business logic (auth, notes, share, storage, ...)
- `migrations/`  : File SQL migration
- `configs/`     : Cấu hình mẫu

## Biến môi trường gợi ý
- `DB_URL`       : Kết nối database
- `JWT_SECRET`   : Secret ký JWT
- `PORT`         : Cổng chạy server

## Tài liệu API
Xem thêm ở thư mục `docs/` hoặc file OpenAPI nếu có.
