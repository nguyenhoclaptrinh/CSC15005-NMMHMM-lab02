# Secure Notes Client (CLI)

## Mô tả
Client dòng lệnh (CLI) cho ứng dụng ghi chú mã hóa. Hỗ trợ đăng ký, đăng nhập, upload/download ghi chú, chia sẻ bảo mật.

## Cách chạy
```bash
cd client
# Cài dependencies
go mod tidy
# Build client
cd cmd
go build -o notescli main.go
# Chạy thử
./notescli --help
```

## Cấu trúc thư mục
- `cmd/`         : Entrypoint CLI
- `internal/`    : Logic gọi API, mã hóa, xử lý file...
- `configs/`     : Cấu hình mẫu

## Biến môi trường gợi ý
- `API_URL`      : Địa chỉ server backend
- `TOKEN_PATH`   : File lưu token tạm thời

## Hướng dẫn sử dụng
- Đăng ký: `./notescli register`
- Đăng nhập: `./notescli login`
- Upload ghi chú: `./notescli upload <file>`
- Download: `./notescli download <note_id>`
- Chia sẻ: `./notescli share <note_id>`
