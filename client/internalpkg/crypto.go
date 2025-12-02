package clientinternal

import (
    "crypto/rsa"
    "errors"
    "crypto/rand"
    "crypto/cipher"
    "crypto/aes"
    "crypto/sha256"
    "io"
    "fmt"
)

// Sinh AES key ngẫu nhiên
func GenerateAESKey() ([]byte, error) {
    // TODO: generate crypto/rand AES key (e.g., 32 bytes for AES-256)
    key := make([]byte, 32)  // AES-256, 32 bytes
    _, err := rand.Read(key)
    if err != nil {
         return nil, errors.New("Failed to generate an AES key!")
    }
    return key, nil
}

// Sinh IV random (Initialization Vector) cho AES-GCM, 12 bytes
func GenerateIV() ([]byte, error){
    iv := make([]byte, 12) 
    if _, err := io.ReadFull(rand.Reader, iv); err != nil {
        return nil, fmt.Errorf("failed to generate IV: %w", err)
    }
    return iv, nil
}

// Mã hóa file bằng AES-GCM
func EncryptFile(aesKey []byte, plaintext []byte) ([]byte, error) {
    // TODO: AES-GCM encrypt
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

// Giải mã file AES-GCM
func DecryptFile(aesKey []byte, ciphertext []byte) ([]byte, error) {
    // TODO: AES-GCM decrypt
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

// Mã hóa AES key bằng public key RSA người nhận
func EncryptAESKeyRSA(aesKey []byte, pubKey *rsa.PublicKey) ([]byte, error) {
    // TODO: RSA-OAEP encrypt
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

// Giải mã AES key bằng private key RSA
func DecryptAESKeyRSA(encKey []byte, privKey *rsa.PrivateKey) ([]byte, error) {
    // TODO: RSA-OAEP decrypt
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
