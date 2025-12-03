``` mermaid
sequenceDiagram
    autonumber
    actor User as Người dùng
    participant Client as Client App (CLI)
    participant Server as API Server
    participant DB as Database

    User->>Client: Nhập Username, Password
    Client->>Server: Gửi Username
    Server->>DB: Lấy Salt của User
    DB-->>Server: Trả về Salt
    Server-->>Client: Trả về Salt
    
    rect rgba(13, 113, 113, 1)
    note right of Client: Xử lý tại Client (Bảo mật)
    Client->>Client: KDF(Password + Salt) -> Tạo K_Master
    Client->>Client: Hash(Password + Salt) -> PasswordHash
    end
    
    Client->>Server: Gửi PasswordHash (Login Request)
    Server->>DB: Kiểm tra PasswordHash
    alt Password Đúng
        Server->>Client: Trả về JWT Token (Login Success)
        note right of Client: Lưu K_Master vào bộ nhớ tạm (RAM)
    else Password Sai
        Server-->>Client: Báo lỗi 401
    end
```