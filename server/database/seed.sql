-- Seed data cho cơ sở dữ liệu Secure Notes
-- Dữ liệu mẫu để test và phát triển

-- Dữ liệu mẫu cho bảng users (sử dụng UUID string cho id)
INSERT INTO users (id, public_key, username, password) VALUES 
('550e8400-e29b-41d4-a716-446655440001', 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5K', 'admin', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'),
('550e8400-e29b-41d4-a716-446655440002', 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7L', 'user1', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'),
('550e8400-e29b-41d4-a716-446655440003', 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9M', 'user2', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'),
('550e8400-e29b-41d4-a716-446655440004', 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1N', 'testuser', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi');

-- Dữ liệu mẫu cho bảng notes (sử dụng UUID string cho id và user_id)
INSERT INTO notes (id, user_id, expire_at, filename, aes_key_encrypted) VALUES 
('660e8400-e29b-41d4-a716-446655440001', '550e8400-e29b-41d4-a716-446655440001', datetime('now', '+30 days'), 'important_note.txt', 'gAAAAABhZ2J3...encrypted_aes_key_1...'),
('660e8400-e29b-41d4-a716-446655440002', '550e8400-e29b-41d4-a716-446655440001', datetime('now', '+7 days'), 'meeting_notes.md', 'gAAAAABhZ2J4...encrypted_aes_key_2...'),
('660e8400-e29b-41d4-a716-446655440003', '550e8400-e29b-41d4-a716-446655440002', datetime('now', '+14 days'), 'personal_diary.txt', 'gAAAAABhZ2J5...encrypted_aes_key_3...'),
('660e8400-e29b-41d4-a716-446655440004', '550e8400-e29b-41d4-a716-446655440002', NULL, 'permanent_note.txt', 'gAAAAABhZ2J6...encrypted_aes_key_4...'),
('660e8400-e29b-41d4-a716-446655440005', '550e8400-e29b-41d4-a716-446655440003', datetime('now', '+60 days'), 'project_plan.md', 'gAAAAABhZ2J7...encrypted_aes_key_5...'),
('660e8400-e29b-41d4-a716-446655440006', '550e8400-e29b-41d4-a716-446655440004', datetime('now', '+1 day'), 'temp_note.txt', 'gAAAAABhZ2J8...encrypted_aes_key_6...');

-- Dữ liệu mẫu cho bảng shares (sử dụng UUID string cho tất cả foreign key)
INSERT INTO shares (id, note_id, expire_at, shared_to_user_id, aes_key_encrypted, url_token) VALUES 
('770e8400-e29b-41d4-a716-446655440001', '660e8400-e29b-41d4-a716-446655440001', datetime('now', '+7 days'), '550e8400-e29b-41d4-a716-446655440002', 'gAAAAABhZ2K1...shared_aes_key_1...', 'abc123def456ghi789'),
('770e8400-e29b-41d4-a716-446655440002', '660e8400-e29b-41d4-a716-446655440001', datetime('now', '+7 days'), '550e8400-e29b-41d4-a716-446655440003', 'gAAAAABhZ2K2...shared_aes_key_2...', 'def456ghi789jkl012'),
('770e8400-e29b-41d4-a716-446655440003', '660e8400-e29b-41d4-a716-446655440003', datetime('now', '+3 days'), '550e8400-e29b-41d4-a716-446655440001', 'gAAAAABhZ2K3...shared_aes_key_3...', 'ghi789jkl012mno345'),
('770e8400-e29b-41d4-a716-446655440004', '660e8400-e29b-41d4-a716-446655440005', datetime('now', '+30 days'), '550e8400-e29b-41d4-a716-446655440004', 'gAAAAABhZ2K4...shared_aes_key_4...', 'jkl012mno345pqr678'),
('770e8400-e29b-41d4-a716-446655440005', '660e8400-e29b-41d4-a716-446655440002', NULL, '550e8400-e29b-41d4-a716-446655440004', 'gAAAAABhZ2K5...shared_aes_key_5...', 'mno345pqr678stu901');

-- Thêm một số dữ liệu test cho expired records
INSERT INTO notes (id, user_id, expire_at, filename, aes_key_encrypted) VALUES 
('660e8400-e29b-41d4-a716-446655440007', '550e8400-e29b-41d4-a716-446655440001', datetime('now', '-1 day'), 'expired_note.txt', 'gAAAAABhZ2J9...expired_aes_key...');

INSERT INTO shares (id, note_id, expire_at, shared_to_user_id, aes_key_encrypted, url_token) VALUES 
('770e8400-e29b-41d4-a716-446655440006', '660e8400-e29b-41d4-a716-446655440007', datetime('now', '-1 hour'), '550e8400-e29b-41d4-a716-446655440002', 'gAAAAABhZ2K6...expired_shared_key...', 'expired_token_123');
