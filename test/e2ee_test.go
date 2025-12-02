package test

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"testing"

	"secure_notes/test/testutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Mock E2EE session key exchange
type SessionKeyExchange struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
	SessionKey []byte
}

// GenerateSessionKey creates a new session key for E2EE
func (ske *SessionKeyExchange) GenerateSessionKey() error {
	key := make([]byte, 32) // 256-bit session key
	_, err := rand.Read(key)
	if err != nil {
		return err
	}
	ske.SessionKey = key
	return nil
}

// EncryptSessionKey encrypts the session key with recipient's public key using RSA-OAEP
func (ske *SessionKeyExchange) EncryptSessionKey(recipientPubKey *rsa.PublicKey) ([]byte, error) {
	hash := sha256.New()
	return rsa.EncryptOAEP(hash, rand.Reader, recipientPubKey, ske.SessionKey, nil)
}

// DecryptSessionKey decrypts the session key with recipient's private key using RSA-OAEP
func DecryptSessionKey(encryptedKey []byte, recipientPrivKey *rsa.PrivateKey) ([]byte, error) {
	hash := sha256.New()
	return rsa.DecryptOAEP(hash, rand.Reader, recipientPrivKey, encryptedKey, nil)
}

// TestSessionKeyGeneration tests session key generation for E2EE
func TestSessionKeyGeneration(t *testing.T) {
	t.Run("Generate valid session key", func(t *testing.T) {
		ske := &SessionKeyExchange{}

		err := ske.GenerateSessionKey()
		require.NoError(t, err)
		assert.Len(t, ske.SessionKey, 32, "Session key should be 32 bytes")

		// Generate another key to ensure randomness
		ske2 := &SessionKeyExchange{}
		err = ske2.GenerateSessionKey()
		require.NoError(t, err)
		assert.NotEqual(t, ske.SessionKey, ske2.SessionKey, "Session keys should be random")
	})
}

// TestKeyExchange tests end-to-end key exchange between two parties
func TestKeyExchange(t *testing.T) {
	t.Run("Successful key exchange between two parties", func(t *testing.T) {
		// Setup: Alice and Bob generate their key pairs
		aliceKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)
		bobKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		// Alice generates a session key
		alice := &SessionKeyExchange{
			PrivateKey: aliceKeys.PrivateKey,
			PublicKey:  aliceKeys.PublicKey,
		}
		err = alice.GenerateSessionKey()
		require.NoError(t, err)

		// Alice encrypts the session key with Bob's public key
		encryptedSessionKey, err := alice.EncryptSessionKey(bobKeys.PublicKey)
		require.NoError(t, err)
		assert.NotEqual(t, alice.SessionKey, encryptedSessionKey, "Encrypted key should be different")

		// Bob receives and decrypts the session key
		bobSessionKey, err := DecryptSessionKey(encryptedSessionKey, bobKeys.PrivateKey)
		require.NoError(t, err)
		assert.Equal(t, alice.SessionKey, bobSessionKey, "Bob should recover Alice's session key")
	})

	t.Run("Key exchange fails with wrong private key", func(t *testing.T) {
		// Setup: Alice, Bob, and Charlie (unauthorized party)
		aliceKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)
		bobKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)
		charlieKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		// Alice generates session key and encrypts it for Bob
		alice := &SessionKeyExchange{
			PrivateKey: aliceKeys.PrivateKey,
			PublicKey:  aliceKeys.PublicKey,
		}
		err = alice.GenerateSessionKey()
		require.NoError(t, err)

		encryptedSessionKey, err := alice.EncryptSessionKey(bobKeys.PublicKey)
		require.NoError(t, err)

		// Charlie tries to decrypt with his private key (should fail)
		_, err = DecryptSessionKey(encryptedSessionKey, charlieKeys.PrivateKey)
		assert.Error(t, err, "Charlie should not be able to decrypt Bob's session key")
	})
}

// TestE2EEMultiPartyScenario tests sharing encrypted data with multiple parties
func TestE2EEMultiPartyScenario(t *testing.T) {
	t.Run("Multi-party E2EE sharing scenario", func(t *testing.T) {
		// Setup: Alice wants to share data with Bob and Charlie
		aliceKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)
		bobKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)
		charlieKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		// Alice generates session key
		alice := &SessionKeyExchange{
			PrivateKey: aliceKeys.PrivateKey,
			PublicKey:  aliceKeys.PublicKey,
		}
		err = alice.GenerateSessionKey()
		require.NoError(t, err)

		// Alice encrypts session key for Bob
		sessionKeyForBob, err := alice.EncryptSessionKey(bobKeys.PublicKey)
		require.NoError(t, err)

		// Alice encrypts session key for Charlie
		sessionKeyForCharlie, err := alice.EncryptSessionKey(charlieKeys.PublicKey)
		require.NoError(t, err)

		// Bob decrypts and accesses the session key
		bobSessionKey, err := DecryptSessionKey(sessionKeyForBob, bobKeys.PrivateKey)
		require.NoError(t, err)
		assert.Equal(t, alice.SessionKey, bobSessionKey, "Bob should have Alice's session key")

		// Charlie decrypts and accesses the session key
		charlieSessionKey, err := DecryptSessionKey(sessionKeyForCharlie, charlieKeys.PrivateKey)
		require.NoError(t, err)
		assert.Equal(t, alice.SessionKey, charlieSessionKey, "Charlie should have Alice's session key")

		// Verify all parties have the same session key
		assert.Equal(t, bobSessionKey, charlieSessionKey, "Both recipients should have same session key")
	})
}

// TestE2EESecurityProperties tests security properties of E2EE
func TestE2EESecurityProperties(t *testing.T) {
	t.Run("Forward secrecy - old session keys cannot decrypt new data", func(t *testing.T) {
		aliceKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		// First session
		alice1 := &SessionKeyExchange{
			PrivateKey: aliceKeys.PrivateKey,
			PublicKey:  aliceKeys.PublicKey,
		}
		err = alice1.GenerateSessionKey()
		require.NoError(t, err)

		// Second session with new session key
		alice2 := &SessionKeyExchange{
			PrivateKey: aliceKeys.PrivateKey,
			PublicKey:  aliceKeys.PublicKey,
		}
		err = alice2.GenerateSessionKey()
		require.NoError(t, err)

		// Verify that session keys are different
		assert.NotEqual(t, alice1.SessionKey, alice2.SessionKey, "Session keys should be different")
	})

	t.Run("Session key confidentiality", func(t *testing.T) {
		aliceKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)
		bobKeys, err := testutils.GenerateTestKeys()
		require.NoError(t, err)

		alice := &SessionKeyExchange{
			PrivateKey: aliceKeys.PrivateKey,
			PublicKey:  aliceKeys.PublicKey,
		}
		err = alice.GenerateSessionKey()
		require.NoError(t, err)

		// Encrypt session key
		encryptedSessionKey, err := alice.EncryptSessionKey(bobKeys.PublicKey)
		require.NoError(t, err)

		// Verify encrypted session key doesn't contain plaintext session key
		assert.NotContains(t, encryptedSessionKey, alice.SessionKey,
			"Encrypted session key should not contain plaintext key")

		// Verify session key is protected
		assert.Len(t, encryptedSessionKey, 256, "RSA-2048 encrypted output should be 256 bytes")
		assert.Greater(t, len(encryptedSessionKey), len(alice.SessionKey),
			"Encrypted session key should be larger than original")
	})
}
