-- =====================================================
-- SEED DATA - Secure Notes Database
-- Dữ liệu mẫu để test và development
-- =====================================================

-- =====================================================
-- TABLE: users
-- Password: "password123" (Argon2id hash)
-- KDF Salt: Random hex string (dùng để derive K_Master ở client)
-- =====================================================
INSERT INTO users (id, username, password_hash, kdf_salt) VALUES 
-- Admin user
('550e8400-e29b-41d4-a716-446655440001', 
 'admin', 
 '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0MTIz$kCXUjmjg5xZqjJz8vZ4kY8h3FqL7rV9xN2mK8pQ1wYc',
 'a1b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef123456'),

-- Regular users
('550e8400-e29b-41d4-a716-446655440002', 
 'alice', 
 '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0MTIz$kCXUjmjg5xZqjJz8vZ4kY8h3FqL7rV9xN2mK8pQ1wYc',
 'b2c3d4e5f6789012345678901234567890abcdef1234567890abcdef1234567'),

('550e8400-e29b-41d4-a716-446655440003', 
 'bob', 
 '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0MTIz$kCXUjmjg5xZqjJz8vZ4kY8h3FqL7rV9xN2mK8pQ1wYc',
 'c3d4e5f6789012345678901234567890abcdef1234567890abcdef12345678'),

('550e8400-e29b-41d4-a716-446655440004', 
 'charlie', 
 '$argon2id$v=19$m=65536,t=3,p=4$c2FsdHNhbHRzYWx0MTIz$kCXUjmjg5xZqjJz8vZ4kY8h3FqL7rV9xN2mK8pQ1wYc',
 'd4e5f6789012345678901234567890abcdef1234567890abcdef123456789');

-- =====================================================
-- TABLE: user_keys
-- Diffie-Hellman public keys (Base64 encoded)
-- =====================================================
INSERT INTO user_keys (user_id, public_key) VALUES 
('550e8400-e29b-41d4-a716-446655440001', 
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5KzVmGFqxqHqJ3kL+8VhQvX0M7jKZqPmN3wYlTuRV8'),

('550e8400-e29b-41d4-a716-446655440002', 
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7LnYpHFrtRsM4nO+9WhRwY1N8kPtZrQoP4yXmWsU9J'),

('550e8400-e29b-41d4-a716-446655440003', 
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9MqTzKJuvRtL5pP+0XiSxZ2O9lQvYsRpQ5zYnXtV0L'),

('550e8400-e29b-41d4-a716-446655440004', 
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1NrWuLKwxStN6qQ+1YjTyA3P0mRwZtSpR6AZoYuW1M');

-- =====================================================
-- TABLE: notes
-- Encrypted notes with title, content, and wrapped keys
-- IV Meta: JSON containing all IVs for AES-GCM operations
-- =====================================================
INSERT INTO notes (id, user_id, title_enc, content_enc, key_enc, iv_meta) VALUES 
-- Admin's notes
('660e8400-e29b-41d4-a716-446655440001', 
 '550e8400-e29b-41d4-a716-446655440001',
 'U2FsdGVkX1+encrypted_title_data_1',
 X'53616c7465645f5f656e637279707465645f636f6e74656e745f646174615f31',
 'gAAAAABhZ2J3...wrapped_K_Note_1...',
 '{"iv_title":"abc123def456","iv_content":"def456ghi789","iv_key":"ghi789jkl012"}'),

('660e8400-e29b-41d4-a716-446655440002', 
 '550e8400-e29b-41d4-a716-446655440001',
 'U2FsdGVkX1+encrypted_title_data_2',
 X'53616c7465645f5f656e637279707465645f636f6e74656e745f646174615f32',
 'gAAAAABhZ2J4...wrapped_K_Note_2...',
 '{"iv_title":"jkl012mno345","iv_content":"mno345pqr678","iv_key":"pqr678stu901"}'),

-- Alice's notes
('660e8400-e29b-41d4-a716-446655440003', 
 '550e8400-e29b-41d4-a716-446655440002',
 'U2FsdGVkX1+encrypted_title_data_3',
 X'53616c7465645f5f656e637279707465645f636f6e74656e745f646174615f33',
 'gAAAAABhZ2J5...wrapped_K_Note_3...',
 '{"iv_title":"stu901vwx234","iv_content":"vwx234yza567","iv_key":"yza567bcd890"}'),

('660e8400-e29b-41d4-a716-446655440004', 
 '550e8400-e29b-41d4-a716-446655440002',
 'U2FsdGVkX1+encrypted_title_data_4',
 X'53616c7465645f5f656e637279707465645f636f6e74656e745f646174615f34',
 'gAAAAABhZ2J6...wrapped_K_Note_4...',
 '{"iv_title":"bcd890efg123","iv_content":"efg123hij456","iv_key":"hij456klm789"}'),

-- Bob's notes
('660e8400-e29b-41d4-a716-446655440005', 
 '550e8400-e29b-41d4-a716-446655440003',
 'U2FsdGVkX1+encrypted_title_data_5',
 X'53616c7465645f5f656e637279707465645f636f6e74656e745f646174615f35',
 'gAAAAABhZ2J7...wrapped_K_Note_5...',
 '{"iv_title":"klm789nop012","iv_content":"nop012qrs345","iv_key":"qrs345tuv678"}'),

-- Charlie's notes
('660e8400-e29b-41d4-a716-446655440006', 
 '550e8400-e29b-41d4-a716-446655440004',
 'U2FsdGVkX1+encrypted_title_data_6',
 X'53616c7465645f5f656e637279707465645f636f6e74656e745f646174615f36',
 'gAAAAABhZ2J8...wrapped_K_Note_6...',
 '{"iv_title":"tuv678wxy901","iv_content":"wxy901zab234","iv_key":"zab234cde567"}');

-- =====================================================
-- TABLE: shared_links
-- Public shared links with advanced security features
-- =====================================================
INSERT INTO shared_links (id, owner_id, content_enc, sender_public_key, expires_at, max_views, current_views, has_password, access_hash, is_active) VALUES 
-- Link không hết hạn, không limit views, không password
('770e8400-e29b-41d4-a716-446655440001', 
 '550e8400-e29b-41d4-a716-446655440001',
 X'53616c7465645f5f7368617265645f636f6e74656e745f31',
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA5KzVm...',
 NULL,
 NULL,
 0,
 FALSE,
 NULL,
 TRUE),

-- Link hết hạn sau 7 ngày, max 10 views, có password
('770e8400-e29b-41d4-a716-446655440002', 
 '550e8400-e29b-41d4-a716-446655440002',
 X'53616c7465645f5f7368617265645f636f6e74656e745f32',
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA7LnYp...',
 datetime('now', '+7 days'),
 10,
 3,
 TRUE,
 'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',
 TRUE),

-- Link hết hạn sau 24h, max 1 view (one-time link)
('770e8400-e29b-41d4-a716-446655440003', 
 '550e8400-e29b-41d4-a716-446655440003',
 X'53616c7465645f5f7368617265645f636f6e74656e745f33',
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA9MqTz...',
 datetime('now', '+1 day'),
 1,
 0,
 FALSE,
 NULL,
 TRUE),

-- Link đã bị revoke (is_active = FALSE)
('770e8400-e29b-41d4-a716-446655440004', 
 '550e8400-e29b-41d4-a716-446655440004',
 X'53616c7465645f5f7368617265645f636f6e74656e745f34',
 'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA1NrWu...',
 datetime('now', '+30 days'),
 NULL,
 0,
 FALSE,
 NULL,
 FALSE),

-- Link đã hết hạn (expired)
('770e8400-e29b-41d4-a716-446655440005', 
 '550e8400-e29b-41d4-a716-446655440001',
 X'53616c7465645f5f7368617265645f636f6e74656e745f35',
 NULL,
 datetime('now', '-1 day'),
 5,
 2,
 FALSE,
 NULL,
 TRUE),

-- Link đã đạt max views (current_views >= max_views)
('770e8400-e29b-41d4-a716-446655440006', 
 '550e8400-e29b-41d4-a716-446655440002',
 X'53616c7465645f5f7368617265645f636f6e74656e745f36',
 NULL,
 datetime('now', '+7 days'),
 3,
 3,
 FALSE,
 NULL,
 TRUE);

-- =====================================================
-- TABLE: refresh_tokens
-- Sample refresh tokens for testing
-- =====================================================
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, is_revoked) VALUES 
-- Active tokens
('880e8400-e29b-41d4-a716-446655440001',
 '550e8400-e29b-41d4-a716-446655440001',
 'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3',
 datetime('now', '+30 days'),
 FALSE),

('880e8400-e29b-41d4-a716-446655440002',
 '550e8400-e29b-41d4-a716-446655440002',
 'b665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3',
 datetime('now', '+30 days'),
 FALSE),

-- Revoked token
('880e8400-e29b-41d4-a716-446655440003',
 '550e8400-e29b-41d4-a716-446655440003',
 'c665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3',
 datetime('now', '+30 days'),
 TRUE),

-- Expired token
('880e8400-e29b-41d4-a716-446655440004',
 '550e8400-e29b-41d4-a716-446655440004',
 'd665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3',
 datetime('now', '-1 day'),
 FALSE);

-- =====================================================
-- TABLE: token_blacklist
-- Blacklisted JWT tokens (logged out or revoked)
-- =====================================================
INSERT INTO token_blacklist (jti, expires_at, reason) VALUES 
-- Recently blacklisted (logout)
('990e8400-e29b-41d4-a716-446655440001',
 datetime('now', '+1 day'),
 'logout'),

-- Security revoke
('990e8400-e29b-41d4-a716-446655440002',
 datetime('now', '+2 days'),
 'security'),

-- Admin force logout
('990e8400-e29b-41d4-a716-446655440003',
 datetime('now', '+3 days'),
 'admin_revoke'),

-- Already expired (should be cleaned up)
('990e8400-e29b-41d4-a716-446655440004',
 datetime('now', '-1 hour'),
 'logout');

-- =====================================================
-- SUMMARY
-- =====================================================
-- Users: 4 (admin, alice, bob, charlie)
-- User Keys: 4 DH public keys
-- Notes: 6 encrypted notes
-- Shared Links: 6 (with various states: active, expired, revoked, max_views_reached)
-- Refresh Tokens: 4 (active, revoked, expired)
-- Token Blacklist: 4 (various reasons)
-- =====================================================

