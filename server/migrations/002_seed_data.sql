-- Seed data for testing/development (SQLite3 compatible)
-- WARNING: Remove this file in production!

-- ============================================================
-- SEED USERS
-- ============================================================

-- Test User 1: alice (password: Alice@123)
INSERT INTO users (id, username, password_hash, kdf_salt, created_at) VALUES
(
    'a0eebc999c0b4ef8bb6d6bb9bd380a11',
    'alice',
    '$argon2id$v=19$m=65536,t=1,p=4$c2FsdDEyMzQ1Njc4OTAxMg$hashed_password_placeholder_alice',
    'c2FsdDEyMzQ1Njc4OTAxMg==',
    datetime('now')
);

-- Test User 2: bob (password: Bob@456)
INSERT INTO users (id, username, password_hash, kdf_salt, created_at) VALUES
(
    'b1ffbc999c0b4ef8bb6d6bb9bd380a22',
    'bob',
    '$argon2id$v=19$m=65536,t=1,p=4$Ym9ic2FsdDEyMzQ1Njc4OQ$hashed_password_placeholder_bob',
    'Ym9ic2FsdDEyMzQ1Njc4OQ==',
    datetime('now')
);

-- Test User 3: charlie (password: Charlie@789)
INSERT INTO users (id, username, password_hash, kdf_salt, created_at) VALUES
(
    'c2ggcc999c0b4ef8bb6d6bb9bd380a33',
    'charlie',
    '$argon2id$v=19$m=65536,t=1,p=4$Y2hhcmxpZXNhbHQxMjM0NQ$hashed_password_placeholder_charlie',
    'Y2hhcmxpZXNhbHQxMjM0NQ==',
    datetime('now')
);

-- ============================================================
-- SEED USER_KEYS (DH Public Keys)
-- ============================================================

-- Alice's DH public key
INSERT INTO user_keys (user_id, public_key, created_at) VALUES
(
    'a0eebc999c0b4ef8bb6d6bb9bd380a11',
    'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...alice_public_key_base64...',
    datetime('now')
);

-- Bob's DH public key
INSERT INTO user_keys (user_id, public_key, created_at) VALUES
(
    'b1ffbc999c0b4ef8bb6d6bb9bd380a22',
    'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...bob_public_key_base64...',
    datetime('now')
);

-- ============================================================
-- SEED NOTES (Encrypted Notes)
-- ============================================================

-- Alice's note 1
INSERT INTO notes (id, user_id, title_enc, content_enc, key_enc, iv_meta, created_at) VALUES
(
    'd0eebc999c0b4ef8bb6d6bb9bd380a44',
    'a0eebc999c0b4ef8bb6d6bb9bd380a11',
    'encrypted_title_alice_note1',
    X'1234567890abcdef',  -- SQLite BLOB hex format
    'encrypted_k_note_with_k_master_alice',
    '{"iv_file": "abc123def456", "iv_key": "xyz789uvw012"}',
    datetime('now')
);

-- Alice's note 2
INSERT INTO notes (id, user_id, title_enc, content_enc, key_enc, iv_meta, created_at) VALUES
(
    'd1ffbc999c0b4ef8bb6d6bb9bd380a55',
    'a0eebc999c0b4ef8bb6d6bb9bd380a11',
    'encrypted_title_alice_note2',
    X'abcdef1234567890',
    'encrypted_k_note_with_k_master_alice_2',
    '{"iv_file": "ghi456jkl789", "iv_key": "mno012pqr345"}',
    datetime('now')
);

-- Bob's note 1
INSERT INTO notes (id, user_id, title_enc, content_enc, key_enc, iv_meta, created_at) VALUES
(
    'e0eebc999c0b4ef8bb6d6bb9bd380a66',
    'b1ffbc999c0b4ef8bb6d6bb9bd380a22',
    'encrypted_title_bob_note1',
    X'fedcba9876543210',
    'encrypted_k_note_with_k_master_bob',
    '{"iv_file": "stu678vwx901", "iv_key": "yza234bcd567"}',
    datetime('now')
);

-- ============================================================
-- SEED SHARED_LINKS (Share Links)
-- ============================================================

-- Alice shares note to Bob (expires in 7 days, max 5 views)
INSERT INTO shared_links (id, owner_id, content_enc, sender_public_key, expires_at, max_views, current_views, has_password, access_hash, is_active, created_at) VALUES
(
    'f0eebc999c0b4ef8bb6d6bb9bd380a77',
    'a0eebc999c0b4ef8bb6d6bb9bd380a11',
    X'1122334455667788',
    'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...alice_sender_public_key...',
    datetime('now', '+7 days'),
    5,
    0,
    0,  -- SQLite: 0=false
    NULL,
    1,  -- SQLite: 1=true
    datetime('now')
);

-- Bob shares note with password protection (expires in 1 hour, unlimited views)
INSERT INTO shared_links (id, owner_id, content_enc, sender_public_key, expires_at, max_views, current_views, has_password, access_hash, is_active, created_at) VALUES
(
    'f1ffbc999c0b4ef8bb6d6bb9bd380a88',
    'b1ffbc999c0b4ef8bb6d6bb9bd380a22',
    X'99aabbccddeeff00',
    'MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA...bob_sender_public_key...',
    datetime('now', '+1 hour'),
    NULL,  -- Unlimited views
    0,
    1,  -- SQLite: 1=true
    'a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3',  -- SHA256("test123")
    1,
    datetime('now')
);

-- Charlie shares note (one-time view, expires in 24 hours)
INSERT INTO shared_links (id, owner_id, content_enc, sender_public_key, expires_at, max_views, current_views, has_password, access_hash, is_active, created_at) VALUES
(
    'f2ggcc999c0b4ef8bb6d6bb9bd380a99',
    'c2ggcc999c0b4ef8bb6d6bb9bd380a33',
    X'1133557799bbddff',
    NULL,  -- No DH encryption (password-only protection)
    datetime('now', '+1 day'),
    1,     -- One-time view
    0,
    1,
    'e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855',  -- SHA256("secret")
    1,
    datetime('now')
);

-- ============================================================
-- SEED REFRESH_TOKENS
-- ============================================================

-- Active refresh token for alice
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, is_revoked, created_at) VALUES
(
    'g0eebc99-9c0b-4ef8-bb6d-6bb9bd380aaa',
    'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11',
    'abc123def456ghi789jkl012mno345pqr678stu901vwx234yza567bcd890efg123',
    datetime('now', '+7 days'),
    0,
    datetime('now')
);

-- Expired refresh token for bob (for testing cleanup)
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, is_revoked, created_at) VALUES
(
    'g1ffbc99-9c0b-4ef8-bb6d-6bb9bd380bbb',
    'b1ffbc99-9c0b-4ef8-bb6d-6bb9bd380a22',
    'xyz789uvw012rst345opq678lmn901ijk234fgh567cde890abc123def456ghi789',
    datetime('now', '-1 day'),  -- Expired 1 day ago
    0,
    datetime('now', '-8 days')
);

-- Revoked refresh token for charlie
INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, is_revoked, created_at) VALUES
(
    'g2ggcc99-9c0b-4ef8-bb6d-6bb9bd380ccc',
    'c2ggcc99-9c0b-4ef8-bb6d-6bb9bd380a33',
    'mno345pqr678stu901vwx234yza567bcd890efg123abc456def789ghi012jkl345',
    datetime('now', '+7 days'),
    1,   -- Revoked
    datetime('now')
);

-- ============================================================
-- SEED TOKEN_BLACKLIST
-- ============================================================

-- Blacklisted JWT from alice's logout
INSERT INTO token_blacklist (jti, expires_at, created_at) VALUES
(
    'jti-alice-logout-12345678-1234-1234-1234-123456789012',
    datetime('now', '+15 minutes'),
    datetime('now')
);

-- Expired blacklisted token (for testing cleanup)
INSERT INTO token_blacklist (jti, expires_at, created_at) VALUES
(
    'jti-old-token-87654321-4321-4321-4321-210987654321',
    datetime('now', '-1 hour'),  -- Expired
    datetime('now', '-2 hours')
);

-- ============================================================
-- VERIFICATION QUERIES (Optional - for manual testing)
-- ============================================================

-- Verify users
-- SELECT id, username, kdf_salt FROM users;

-- Verify notes
-- SELECT id, user_id, title_enc FROM notes;

-- Verify active share links
-- SELECT * FROM active_shared_links;

-- Verify refresh tokens
-- SELECT id, user_id, is_revoked, expires_at FROM refresh_tokens;

-- ============================================================
-- NOTES FOR DEVELOPERS
-- ============================================================

/*
PASSWORD HASHES (Placeholders above need to be generated):
- Use Argon2id with parameters: time=1, memory=64MB, parallelism=4
- Generate real hashes in your application code

EXAMPLE (pseudo-code):
  alice_hash = argon2id.hash("Alice@123", salt="salt123456789012", ...)
  bob_hash = argon2id.hash("Bob@456", salt="bobsalt123456789", ...)

DH PUBLIC KEYS:
- Base64 encoded DH public keys (2048-bit recommended)
- Generate using crypto library in your application

ENCRYPTED DATA:
- All title_enc, content_enc, key_enc are placeholders
- Real data should be AES-256-GCM encrypted

CLEANUP:
- Run `SELECT cleanup_expired_data();` periodically
- Or set up cron job to clean expired tokens/links
*/
