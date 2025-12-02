package test

import (
	"testing"

	clientinternal "secure_notes/client/internalpkg"
	"secure_notes/test/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestGenerateAESKey tests AES key generation
func TestGenerateAESKey(t *testing.T) {
	t.Run("Generate valid AES key", func(t *testing.T) {
		key, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)
		assert.Len(t, key, 32, "AES key should be 32 bytes (256 bits)")

		// Generate another key and ensure they are different
		key2, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)
		assert.NotEqual(t, key, key2, "Generated keys should be random and different")
	})
}

// TestFileEncryptionDecryption tests file encryption and decryption
func TestFileEncryptionDecryption(t *testing.T) {
	t.Run("Encrypt and decrypt file successfully", func(t *testing.T) {
		// Generate AES key
		aesKey, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)

		// Original plaintext
		originalText := []byte("This is a secret note content that needs to be encrypted!")

		// Encrypt
		ciphertext, err := clientinternal.EncryptFile(aesKey, originalText)
		require.NoError(t, err)
		assert.NotEqual(t, originalText, ciphertext, "Ciphertext should be different from plaintext")
		assert.Greater(t, len(ciphertext), len(originalText), "Ciphertext should be longer (includes GCM tag)")

		// Decrypt
		decryptedText, err := clientinternal.DecryptFile(aesKey, ciphertext)
		require.NoError(t, err)
		assert.Equal(t, originalText, decryptedText, "Decrypted text should match original")
	})

	t.Run("Decrypt with wrong key fails", func(t *testing.T) {
		// Generate two different AES keys
		aesKey1, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)

		aesKey2, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)

		originalText := []byte("Secret data")

		// Encrypt with key1
		ciphertext, err := clientinternal.EncryptFile(aesKey1, originalText)
		require.NoError(t, err)

		// Try to decrypt with key2 (should fail)
		_, err = clientinternal.DecryptFile(aesKey2, ciphertext)
		assert.Error(t, err, "Decryption with wrong key should fail")
	})
}

// TestRSAEncryptionDecryption tests RSA encryption/decryption of AES keys
func TestRSAEncryptionDecryption(t *testing.T) {
	t.Run("Encrypt and decrypt AES key with RSA successfully", func(t *testing.T) {
		// Generate test keys
		testKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		// Generate AES key to encrypt
		aesKey, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)

		// Encrypt AES key with RSA public key
		encryptedAESKey, err := clientinternal.EncryptAESKeyRSA(aesKey, testKeys.PublicKey)
		require.NoError(t, err)
		assert.NotEqual(t, aesKey, encryptedAESKey, "Encrypted key should be different")
		assert.Greater(t, len(encryptedAESKey), len(aesKey), "Encrypted key should be larger")

		// Decrypt AES key with RSA private key
		decryptedAESKey, err := clientinternal.DecryptAESKeyRSA(encryptedAESKey, testKeys.PrivateKey)
		require.NoError(t, err)
		assert.Equal(t, aesKey, decryptedAESKey, "Decrypted AES key should match original")
	})

	t.Run("Decrypt with wrong private key fails", func(t *testing.T) {
		// Generate two different key pairs
		testKeys1, err := testutils.GenerateTestKeys()
		require.NoError(t, err)
		testKeys2, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		aesKey, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)

		// Encrypt with public key 1
		encryptedAESKey, err := clientinternal.EncryptAESKeyRSA(aesKey, testKeys1.PublicKey)
		require.NoError(t, err)

		// Try to decrypt with private key 2
		_, err = clientinternal.DecryptAESKeyRSA(encryptedAESKey, testKeys2.PrivateKey)
		assert.Error(t, err, "Decryption with wrong private key should fail")
	})
}

// TestEndToEndEncryption tests complete encryption workflow
func TestEndToEndEncryption(t *testing.T) {
	t.Run("Complete encryption workflow", func(t *testing.T) {
		// Setup: Generate RSA keys for receiver
		receiverKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		// Step 1: Sender generates AES key
		aesKey, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)

		// Step 2: Sender encrypts file content with AES
		originalContent := []byte("Confidential document content")
		encryptedContent, err := clientinternal.EncryptFile(aesKey, originalContent)
		require.NoError(t, err)

		// Step 3: Encrypt AES key with receiver's public key
		encryptedAESKey, err := clientinternal.EncryptAESKeyRSA(aesKey, receiverKeys.PublicKey)
		require.NoError(t, err)

		// === Simulation: Data is stored/transmitted ===

		// Step 4: Receiver decrypts AES key with their private key
		decryptedAESKey, err := clientinternal.DecryptAESKeyRSA(encryptedAESKey, receiverKeys.PrivateKey)
		require.NoError(t, err)
		assert.Equal(t, aesKey, decryptedAESKey, "AES key should be correctly recovered")

		// Step 5: Receiver decrypts file content with AES key
		decryptedContent, err := clientinternal.DecryptFile(decryptedAESKey, encryptedContent)
		require.NoError(t, err)
		assert.Equal(t, originalContent, decryptedContent, "Original content should be recovered")
	})
}

// TestKeyProtection tests that encryption keys are properly protected
func TestKeyProtection(t *testing.T) {
	t.Run("AES key is properly protected", func(t *testing.T) {
		testKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		aesKey, err := clientinternal.GenerateAESKey()
		require.NoError(t, err)

		// Encrypt AES key
		encryptedKey, err := clientinternal.EncryptAESKeyRSA(aesKey, testKeys.PublicKey)
		require.NoError(t, err)

		// Verify that encrypted key doesn't contain original key data
		assert.NotContains(t, encryptedKey, aesKey, "Encrypted key should not contain plaintext key")

		// Verify key can only be decrypted with correct private key
		decryptedKey, err := clientinternal.DecryptAESKeyRSA(encryptedKey, testKeys.PrivateKey)
		require.NoError(t, err)
		assert.Equal(t, aesKey, decryptedKey, "Key should decrypt correctly with right private key")
	})
}
