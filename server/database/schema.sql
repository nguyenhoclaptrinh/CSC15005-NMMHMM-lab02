-- =====================================================
-- SECURE NOTES DATABASE SCHEMA (Updated Design)
-- End-to-End Encryption with Advanced Security Features
-- =====================================================

-- =====================================================
-- TABLE: users
-- Lưu thông tin tài khoản người dùng
-- =====================================================
CREATE TABLE IF NOT EXISTS users (
    id              TEXT PRIMARY KEY,           -- UUID (an toàn hơn auto-increment)
    username        VARCHAR(50) NOT NULL UNIQUE,-- Tên đăng nhập
    password_hash   VARCHAR(255) NOT NULL,      -- Argon2id hash để xác thực
    kdf_salt        VARCHAR(64) NOT NULL,       -- Salt cho KDF ở client (sinh K_Master)
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_username ON users(username);

-- =====================================================
-- TABLE: user_keys
-- Lưu khóa công khai Diffie-Hellman của user
-- =====================================================
CREATE TABLE IF NOT EXISTS user_keys (
    user_id         TEXT PRIMARY KEY,           -- Khóa chính + Khóa ngoại
    public_key      TEXT NOT NULL,              -- DH public key (Base64)
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- =====================================================
-- TABLE: notes
-- Lưu ghi chú đã mã hóa của user
-- =====================================================
CREATE TABLE IF NOT EXISTS notes (
    id              TEXT PRIMARY KEY,           -- UUID
    user_id         TEXT NOT NULL,              -- Chủ sở hữu note
    title_enc       TEXT NOT NULL,              -- Tiêu đề đã mã hóa
    content_enc     BLOB NOT NULL,              -- Nội dung file đã mã hóa
    key_enc         TEXT NOT NULL,              -- K_Note được mã hóa bởi K_Master (Key Wrapping)
    iv_meta         TEXT NOT NULL,              -- JSON: {"iv_title": "...", "iv_content": "...", "iv_key": "..."}
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_notes_created_at ON notes(created_at);

-- =====================================================
-- TABLE: shared_links
-- Quản lý link chia sẻ với các tính năng nâng cao
-- =====================================================
CREATE TABLE IF NOT EXISTS shared_links (
    id                  TEXT PRIMARY KEY,       -- UUID (chính là ID trên URL /share/{id})
    owner_id            TEXT NOT NULL,          -- Người tạo link
    content_enc         BLOB NOT NULL,          -- Nội dung đã mã hóa riêng cho share
    sender_public_key   TEXT,                   -- DH Public Key của người gửi (cho async encryption)
    expires_at          TIMESTAMP,              -- Thời điểm hết hạn (NULL = không hết hạn)
    max_views           INTEGER,                -- Số lượt xem tối đa (NULL = không giới hạn)
    current_views       INTEGER DEFAULT 0,      -- Số lượt đã xem
    has_password        BOOLEAN DEFAULT FALSE,  -- Link có yêu cầu mật khẩu không
    access_hash         VARCHAR(64),            -- SHA256 hash của password truy cập
    is_active           BOOLEAN DEFAULT TRUE,   -- Trạng thái link (false = đã thu hồi)
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    last_accessed_at    TIMESTAMP,              -- Lần truy cập gần nhất
    
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_shared_links_owner ON shared_links(owner_id);
CREATE INDEX idx_shared_links_expires ON shared_links(expires_at);
CREATE INDEX idx_shared_links_active ON shared_links(is_active);

-- =====================================================
-- TABLE: refresh_tokens
-- Quản lý Refresh Tokens cho authentication
-- =====================================================
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id              TEXT PRIMARY KEY,           -- UUID
    user_id         TEXT NOT NULL,              -- Token của user nào
    token_hash      VARCHAR(64) NOT NULL,       -- SHA256 hash của refresh token gốc
    expires_at      TIMESTAMP NOT NULL,         -- Thời gian hết hạn
    is_revoked      BOOLEAN DEFAULT FALSE,      -- Đã thu hồi chưa
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    revoked_at      TIMESTAMP,                  -- Thời điểm thu hồi
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_user ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_expires ON refresh_tokens(expires_at);

-- =====================================================
-- TABLE: token_blacklist
-- Blacklist JWT tokens đã logout hoặc bị thu hồi
-- =====================================================
CREATE TABLE IF NOT EXISTS token_blacklist (
    jti             VARCHAR(36) PRIMARY KEY,    -- JWT ID (unique identifier)
    expires_at      TIMESTAMP NOT NULL,         -- Token expiry time
    blacklisted_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    reason          VARCHAR(50)                 -- Lý do: "logout", "security", "admin_revoke"
);

CREATE INDEX idx_token_blacklist_expires ON token_blacklist(expires_at);

-- =====================================================
-- CLEANUP EXPIRED DATA (Chạy định kỳ)
-- =====================================================

-- Xóa shared_links đã hết hạn
-- DELETE FROM shared_links WHERE expires_at < CURRENT_TIMESTAMP AND expires_at IS NOT NULL;

-- Xóa refresh_tokens đã hết hạn
-- DELETE FROM refresh_tokens WHERE expires_at < CURRENT_TIMESTAMP;

-- Xóa JWT blacklist đã hết hạn
-- DELETE FROM token_blacklist WHERE expires_at < CURRENT_TIMESTAMP;

-- =====================================================
-- SAMPLE QUERIES
-- =====================================================

-- Kiểm tra link còn valid không
-- SELECT * FROM shared_links 
-- WHERE id = ? 
--   AND is_active = TRUE 
--   AND (expires_at IS NULL OR expires_at > CURRENT_TIMESTAMP)
--   AND (max_views IS NULL OR current_views < max_views);

-- Tăng view count
-- UPDATE shared_links SET current_views = current_views + 1, last_accessed_at = CURRENT_TIMESTAMP WHERE id = ?;

-- Thu hồi tất cả refresh tokens của user
-- UPDATE refresh_tokens SET is_revoked = TRUE, revoked_at = CURRENT_TIMESTAMP WHERE user_id = ?;

-- Kiểm tra JWT có trong blacklist không
-- SELECT COUNT(*) FROM token_blacklist WHERE jti = ? AND expires_at > CURRENT_TIMESTAMP;
