package serverpkg

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/cipher"
	"crypto/aes"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"math/big"
	"io"
	"fmt"
	"strings"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/hkdf"
)

// ============================================================
// AES-GCM ENCRYPTION (for file content)
// ============================================================

// GenerateAESKey generates random AES-256 key (K_Note)
func GenerateAESKey() ([]byte, error) {
	// TODO: Generate 32-byte key for AES-256
	key := make([]byte, 32)
	_, err := rand.Read(key)
	if err != nil {
         return nil, errors.New("Failed to generate an AES key!")
    }
	return key, err
}
func GenerateIV() ([]byte, error){
    iv := make([]byte, 12) 
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("failed to generate IV: %w", err)
    }
    return iv, nil
}
// EncryptFile encrypts plaintext using AES-256-GCM
func EncryptFile(aesKey []byte, plaintext []byte) ([]byte, error) {
	block, err := aes.NewCipher(aesKey) // Tao cipher block tu aesKey
    if err != nil {
        return nil, fmt.Errorf("Failed to create cipher block: %w", err)
    }

    gcm, err := cipher.NewGCM(block)  // Tao GCM mode tu cipher block
    if err != nil {
        return nil, fmt.Errorf("Failed to create GCM: %w", err)
    }
    
    iv, err := GenerateIV()
    if err != nil {
        return nil, err
    }

    if len(iv) != gcm.NonceSize() {
        return nil, errors.New("IV length does not match GCM nonce size")
    }

    cipherText := gcm.Seal(iv, iv, plaintext, nil)
    return cipherText, nil
}

// DecryptFile decrypts ciphertext using AES-256-GCM
func DecryptFile(aesKey []byte, ciphertext []byte) ([]byte, error) {
	// TODO: AES-GCM decrypt
	//   - Extract nonce (first 12 bytes)
	//   - Use cipher.NewGCM()
	//   - Open() verifies tag and decrypts
	block, err := aes.NewCipher(aesKey)
    if err != nil {
        return nil, fmt.Errorf("Failed to create cipher block: %w", err)
    }

    gcm, err := cipher.NewGCM(block)  
    if err != nil {
        return nil, fmt.Errorf("Failed to create GCM: %w", err)
    }

    nonceSize := gcm.NonceSize()

    if len(ciphertext) < nonceSize {
		return nil, errors.New("ciphertext too short, invalid format")
	}
    
    nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("Failed to decrypt: %w", err)
	}

	return plaintext, nil
}

// ============================================================
// RSA-OAEP KEY WRAPPING (optional, for legacy systems)
// ============================================================

// EncryptAESKeyRSA encrypts AES key using RSA-OAEP
func EncryptAESKeyRSA(aesKey []byte, pubKey *rsa.PublicKey) ([]byte, error) {
	// TODO: RSA-OAEP encrypt
	//   - Use rsa.EncryptOAEP()
	//   - Hash: sha256.New()
	encryptedKey, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pubKey,
		aesKey,
		nil, 
	)
	if err != nil {
		return nil, fmt.Errorf("RSA-OAEP encryption failed: %w", err)
	}

	return encryptedKey, nil
}

// DecryptAESKeyRSA decrypts AES key using RSA-OAEP
func DecryptAESKeyRSA(encKey []byte, privKey *rsa.PrivateKey) ([]byte, error) {
    decryptedKey, err := rsa.DecryptOAEP(
		sha256.New(),
		rand.Reader,
		privKey,
		encKey,
		nil, 
	)
	if err != nil {
		return nil, fmt.Errorf("RSA-OAEP decryption failed: %w", err)
	}

    return decryptedKey, nil
}

// ============================================================
// DIFFIE-HELLMAN KEY EXCHANGE (for E2EE sharing)
// ============================================================

// DHParams holds Diffie-Hellman parameters (p, g)
type DHParams struct {
	P *big.Int // Large prime modulus
	G *big.Int // Generator
}

// DHKeyPair holds private and public keys for DH
type DHKeyPair struct {
	Private *big.Int // Private key: a (secret)
	Public  *big.Int // Public key: A = g^a mod p
}

// GenerateDHParameters creates DH parameters using RFC 3526 Group 14 (2048-bit)
func GenerateDHParameters() (*DHParams, error) {
	// TODO: Use RFC 3526 Group 14 (2048-bit MODP group)
	//   p := new(big.Int)
	//   p.SetString("FFFFFFFF...FFFFFFFF", 16)
	//   g := big.NewInt(2)
	//   return &DHParams{P: p, G: g}, nil
	p :=  new(big.Int)
	p.SetString("FFFFFFFF"+"FFFFFFF"+"C90FDAA2"+"2168C234"+"C4C6628B"+"80DC1CD1"+
      "29024E08"+"8A67CC74"+"020BBEA6"+"3B139B22"+"514A0879"+"8E3404DD"+
      "EF9519B3"+ "CD3A431B"+"302B0A6D"+ "F25F1437"+"4FE1356D"+"6D51C245"+
      "E485B576"+ "625E7EC6"+ "F44C42E9"+ "A637ED6B" + "0BFF5CB6" + "F406B7ED" +
      "EE386BFB" + "5A899FA5" + "AE9F2411" + "7C4B1FE6" + "49286651" + "ECE45B3D" +
      "C2007CB8" + "A163BF05" + "98DA4836" + "1C55D39A" + "69163FA8" + "FD24CF5F" +
      "83655D23" + "DCA3AD96" + "1C62F356" + "208552BB" + "9ED52907" + "7096966D" +
      "670C354E" + "4ABC9804" + "F1746C08" + "CA18217C" + "32905E46" + "2E36CE3B" +
      "E39E772C" + "180E8603" + "9B2783A2" + "EC07A28F" + "B5C55DF0" + "6F4C52C9" +
      "DE2BCBF6" + "95581718" + "3995497C" + "EA956AE5" + "15D22618" + "98FA0510" +
      "15728E5A" + "8AACAA68" + "FFFFFFFF" + "FFFFFFFF", 16)
	
	g := big.NewInt(2)
	return &DHParams{P: p, G: g}, nil
}

// GenerateDHKeyPair generates private key a and public key A = g^a mod p
func GenerateDHKeyPair(params *DHParams) (*DHKeyPair, error) {
	// TODO: Generate random private key a
	//   a, err := rand.Int(rand.Reader, params.P)
	//   A := new(big.Int).Exp(params.G, a, params.P)
	//   return &DHKeyPair{Private: a, Public: A}, nil
	max := new(big.Int).Sub(params.P, big.NewInt(2))
	a, err := rand.Int(rand.Reader, max)   // make sure a in range [0, p - 2]
	if err != nil {
		return nil, fmt.Errorf("Failed to generate private key: %w", err)
	}
	A := new(big.Int).Exp(params.G, a, params.P)    // A = g^a mod p
	return &DHKeyPair{Private: a, Public: A}, nil
}

// ComputeSharedSecret calculates shared secret S = B^a mod p
func ComputeSharedSecret(theirPublic *big.Int, myPrivate *big.Int, params *DHParams) (*big.Int, error) {
	// TODO: Validate their public key (1 < B < p-1)
	//   S := new(big.Int).Exp(theirPublic, myPrivate, params.P)
	//   return S, nil
	if theirPublic.Cmp(big.NewInt(1)) <= 0 {
		return nil, errors.New("public key must be greater than 1")
	}
	
	if theirPublic.Cmp(new(big.Int).Sub(params.P, big.NewInt(1))) >= 0 {
		return nil, errors.New("public key must be less than p-1")
	}
	
	// S = B^a mod 
	S := new(big.Int).Exp(theirPublic, myPrivate, params.P)
	
	return S, nil
}

// DeriveSessionKey uses HKDF to derive AES key from DH shared secret
func DeriveSessionKey(sharedSecret *big.Int) ([]byte, error) {
	// TODO: Use HKDF-SHA256 to derive 32-byte key
	//   - Convert sharedSecret to bytes
	//   - HKDF(secretBytes, salt=nil, info="E2EE-Session-Key")
	//   - Return 32-byte AES-256 key
	secretBytes := sharedSecret.Bytes()
	hkdf := hkdf.New(sha256.New, secretBytes, nil, []byte("E2EE-Session-Key"))
	key := make([]byte, 32)
	if _, err := io.ReadFull(hkdf, key); err != nil {
		return nil, fmt.Errorf("Failed to derive session key: %w", err)
	}
	return key, nil
}

// VerifyKeyFingerprint creates human-readable fingerprint for public key
func VerifyKeyFingerprint(publicKey *big.Int) string {
	// TODO: Hash public key with SHA256
	//   hash := sha256.Sum256(publicKey.Bytes())
	//   Format as "A1:B2:C3:..." (first 16 bytes)
	hash := sha256.Sum256(publicKey.Bytes())
	hexStr := hex.EncodeToString(hash[:16])
	var result strings.Builder
	for i := 0; i < len(hexStr); i += 2 {
		if i > 0 {
			result.WriteString(":")
		}
		result.WriteString(strings.ToUpper(hexStr[i:i+2]))
	}
	
	return result.String()
}

// ============================================================
// KEY DERIVATION (Argon2id for K_Master)
// ============================================================

// DeriveKeyFromPassword uses Argon2id to derive K_Master from password
func DeriveKeyFromPassword(password string, salt []byte) ([]byte, error) {
	// TODO: Use Argon2id with strong parameters:
	//   key := argon2.IDKey(
	//       []byte(password), salt,
	//       1,       // time cost
	//       64*1024, // memory (64 MB)
	//       4,       // parallelism
	//       32,      // key length (AES-256)
	//   )
	//   return key, nil
	key := argon2.IDKey(
		[]byte(password), 
		salt,
		1,       // time cost
		64*1024, // memory (64 MB)
		4,       // parallelism
		32,      // key length (AES-256)
	)
	
	return key, nil
}

// GenerateSalt creates a cryptographically secure random salt
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	_, err := rand.Read(salt)
	return salt, err
}

// GenerateIV creates a random initialization vector for AES-GCM
func GenerateIV() ([]byte, error) {
	// AES-GCM standard nonce size is 12 bytes
	// CRITICAL: Never reuse IV with same key!
	iv := make([]byte, 12)
	_, err := rand.Read(iv)
	return iv, err
}

// ZeroizeKey securely wipes key from memory
func ZeroizeKey(key []byte) {
	// Overwrite key with zeros before freeing memory
	// Call this when: user logs out, after encryption/decryption
	for i := range key {
		key[i] = 0
	}
}
