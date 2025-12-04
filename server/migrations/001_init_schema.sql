-- Database schema for secure note-sharing application
-- SQLite3 compatible schema (using TEXT for UUID)

-- ============================================================
-- TABLE 1: users - User accounts
-- ============================================================
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),  -- UUID text (32 chars hex)
    username TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,                    -- Argon2id hash để xác thực đăng nhập
    kdf_salt TEXT NOT NULL,                         -- Salt dùng cho KDF ở Client để sinh K_Master
    created_at TEXT DEFAULT (datetime('now')),
    last_login TEXT
);

CREATE INDEX idx_users_username ON users(username);

-- ============================================================
-- TABLE 2: user_keys - DH public keys for E2EE
-- ============================================================
CREATE TABLE IF NOT EXISTS user_keys (
    user_id TEXT PRIMARY KEY,                       -- Khóa chính + Khóa ngoại
    public_key TEXT NOT NULL,                       -- Khóa công khai Diffie-Hellman (Base64)
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- ============================================================
-- TABLE 3: notes - Encrypted notes storage
-- ============================================================
CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    user_id TEXT NOT NULL,
    title_enc TEXT NOT NULL,           -- Tiêu đề đã mã hóa (Server không biết user lưu gì)
    content_enc BLOB NOT NULL,         -- Nội dung file đã mã hóa
    key_enc TEXT NOT NULL,             -- Key Wrapping: K_Note được mã hóa bởi K_Master
    iv_meta TEXT NOT NULL,             -- Lưu IV và Auth Tag cho AES-GCM (JSON string)
    created_at TEXT DEFAULT (datetime('now')),
    updated_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_notes_user_id ON notes(user_id);
CREATE INDEX idx_notes_created_at ON notes(created_at);

-- ============================================================
-- TABLE 4: shared_links - Time-sensitive share links
-- ============================================================
CREATE TABLE IF NOT EXISTS shared_links (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))), -- Chính là chuỗi ID trên URL
    owner_id TEXT NOT NULL,                        -- Người tạo link (để kiểm quyền xóa)
    content_enc BLOB NOT NULL,                     -- Nội dung file mã hóa dành riêng cho chia sẻ
    sender_public_key TEXT,                        -- Public Key A của người gửi (Dùng cho DH Async)
    expires_at TEXT,                               -- Time-sensitive: Thời điểm link hết hạn
    max_views INTEGER,                             -- Quota: Số lượt xem tối đa
    current_views INTEGER NOT NULL DEFAULT 0,      -- Đếm số lượt đã xem
    has_password INTEGER NOT NULL DEFAULT 0,       -- SQLite: 0=false, 1=true
    access_hash TEXT,                              -- Hash SHA256 của mật khẩu truy cập (nếu có)
    is_active INTEGER NOT NULL DEFAULT 1,          -- SQLite: 0=false, 1=true
    created_at TEXT DEFAULT (datetime('now')),
    last_accessed_at TEXT,
    FOREIGN KEY (owner_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_shared_links_id ON shared_links(id);
CREATE INDEX idx_shared_links_owner_id ON shared_links(owner_id);
CREATE INDEX idx_shared_links_expires_at ON shared_links(expires_at);
CREATE INDEX idx_shared_links_is_active ON shared_links(is_active);

-- ============================================================
-- TABLE 5: refresh_tokens - Long-lived tokens for JWT refresh
-- ============================================================
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16)))),
    user_id TEXT NOT NULL,
    token_hash TEXT NOT NULL,                    -- Hash SHA256 của Refresh Token gốc (bảo vệ nếu DB lộ)
    expires_at TEXT NOT NULL,                    -- Thời gian hết hạn (thường 7 ngày)
    is_revoked INTEGER NOT NULL DEFAULT 0,       -- SQLite: 0=false, 1=true
    created_at TEXT DEFAULT (datetime('now')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_refresh_tokens_token_hash ON refresh_tokens(token_hash);
CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_expires_at ON refresh_tokens(expires_at);

-- ============================================================
-- TABLE 6: token_blacklist - Revoked JWT tokens
-- ============================================================
CREATE TABLE IF NOT EXISTS token_blacklist (
    jti TEXT PRIMARY KEY,            -- JWT ID (Unique Identifier trong payload của token)
    expires_at TEXT NOT NULL,        -- Thời gian hết hạn của token
    created_at TEXT DEFAULT (datetime('now'))
);

CREATE INDEX idx_token_blacklist_expires_at ON token_blacklist(expires_at);

-- ============================================================
-- TRIGGERS - Automatic cleanup and updates (SQLite3)
-- ============================================================

-- Trigger: Update updated_at on notes table
CREATE TRIGGER IF NOT EXISTS update_notes_timestamp
AFTER UPDATE ON notes
FOR EACH ROW
BEGIN
    UPDATE notes SET updated_at = datetime('now') WHERE id = NEW.id;
END;

-- Trigger: Update updated_at on user_keys table
CREATE TRIGGER IF NOT EXISTS update_user_keys_timestamp
AFTER UPDATE ON user_keys
FOR EACH ROW
BEGIN
    UPDATE user_keys SET updated_at = datetime('now') WHERE user_id = NEW.user_id;
END;

-- ============================================================
-- VIEWS - Useful queries
-- ============================================================

-- View: Active shared links (not expired, not maxed out)
CREATE VIEW IF NOT EXISTS active_shared_links AS
SELECT 
    sl.id,
    sl.owner_id,
    sl.expires_at,
    sl.max_views,
    sl.current_views,
    sl.has_password,
    sl.is_active,
    sl.created_at,
    sl.last_accessed_at
FROM shared_links sl
WHERE 
    sl.is_active = 1
    AND (sl.expires_at IS NULL OR sl.expires_at > datetime('now'))
    AND (sl.max_views IS NULL OR sl.current_views < sl.max_views);

-- ============================================================
-- CLEANUP QUERIES (Run periodically via cron or application)
-- ============================================================

-- SQLite doesn't support stored procedures, so these are standalone queries
-- Run these periodically from your application or cron job

-- 1. Delete expired refresh tokens
-- DELETE FROM refresh_tokens WHERE expires_at < datetime('now');

-- 2. Delete expired blacklisted tokens
-- DELETE FROM token_blacklist WHERE expires_at < datetime('now');

-- 3. Deactivate expired share links
-- UPDATE shared_links 
-- SET is_active = 0 
-- WHERE expires_at < datetime('now') AND is_active = 1;

-- 4. Deactivate share links that reached max views
-- UPDATE shared_links
-- SET is_active = 0
-- WHERE max_views IS NOT NULL 
--   AND current_views >= max_views 
--   AND is_active = 1;
