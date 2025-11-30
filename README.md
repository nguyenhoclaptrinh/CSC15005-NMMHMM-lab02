# Secure Notes - Lab02

Dự án mẫu lưu trữ ghi chú bảo mật, mã hóa đầu-cuối.

## Cấu trúc thư mục

- `client/` — chương trình dòng lệnh (CLI) và các hàm tiện ích client
- `server/` — API server sử dụng Gin
- `server/config/` — cấu hình server
- `server/storage/` — lưu trữ file mã hóa và truy cập database
- `test/` — các test và hàm hỗ trợ kiểm thử

## Hướng dẫn sử dụng

### Build & chạy server

```bash
# Build file thực thi server 
go build -o bin/server ./server/main.go

# Hoặc chạy trực tiếp
go run ./server/main.go
```

### Chạy client

```bash
go run ./client/main.go
```

### Chạy test

```bash
go test ./test -v
```

## Lưu ý phát triển

- Hầu hết các hàm xử lý và logic mã hóa mới chỉ là khung mẫu, có đánh dấu `TODO`.
- Các package nội bộ nằm ở `client/internalpkg` và `server/internalpkg` (có thể đổi tên thành `internal/` để đúng chuẩn Go).
- Một số tính năng SQLite cần CGO; trên Windows nên dùng `CGO_ENABLED=0` nếu chưa cài toolchain GCC 64-bit.
- Xem comment trong code để biết hướng dẫn triển khai chi tiết.
- Sử dụng các thư viện mã hóa chuẩn như `crypto/aes`, `crypto/rand`, `golang.org/x/crypto/scrypt` để đảm bảo an toàn.
- Tuân thủ các nguyên tắc bảo mật khi xử lý khóa và dữ liệu nhạy cảm.
- Cấu trúc API server sử dụng Gin, tham khảo tài liệu chính thức để mở rộng.
- Sử dụng các công cụ quản lý phụ thuộc như Go Modules để quản lý thư viện bên thứ ba.
- Đảm bảo viết test đầy đủ cho các hàm mã hóa và lưu trữ để đảm bảo tính đúng đắn và an toàn.   
