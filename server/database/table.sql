-- Secure Notes Database Schema
-- Tạo các bảng cho ứng dụng secure notes

-- USER table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    public_key TEXT NOT NULL,
    username TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL, -- hashed password
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

-- NOTE table
CREATE TABLE IF NOT EXISTS notes (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    expire_at DATETIME,
    filename TEXT NOT NULL, -- tên file bao gồm extension
    aes_key_encrypted TEXT NOT NULL, -- AES key đã được mã hóa
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);
    
-- SHARE table
CREATE TABLE IF NOT EXISTS shares (
    id TEXT PRIMARY KEY,
    note_id TEXT NOT NULL,
    expire_at DATETIME,
    shared_to_user_id TEXT NOT NULL,
    aes_key_encrypted TEXT NOT NULL, -- AES key đã được mã hóa cho người nhận
    url_token TEXT UNIQUE NOT NULL, -- token để truy cập qua URL
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,

    FOREIGN KEY (note_id) REFERENCES notes(id) ON DELETE CASCADE,
    FOREIGN KEY (shared_to_user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes để tối ưu truy vấn
CREATE INDEX IF NOT EXISTS idx_notes_user_id ON notes(user_id);
CREATE INDEX IF NOT EXISTS idx_notes_expire_at ON notes(expire_at);
CREATE INDEX IF NOT EXISTS idx_shares_note_id ON shares(note_id);
CREATE INDEX IF NOT EXISTS idx_shares_shared_to_user_id ON shares(shared_to_user_id);
CREATE INDEX IF NOT EXISTS idx_shares_url_token ON shares(url_token);
CREATE INDEX IF NOT EXISTS idx_shares_expire_at ON shares(expire_at);

-- Triggers để tự động update updated_at
CREATE TRIGGER IF NOT EXISTS update_users_timestamp 
    AFTER UPDATE ON users
BEGIN
    UPDATE users SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_notes_timestamp 
    AFTER UPDATE ON notes
BEGIN
    UPDATE notes SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

CREATE TRIGGER IF NOT EXISTS update_shares_timestamp 
    AFTER UPDATE ON shares
BEGIN
    UPDATE shares SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;