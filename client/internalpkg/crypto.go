package clientinternal

import (
    "crypto/rsa"
    "errors"
)

// Sinh AES key ngẫu nhiên
func GenerateAESKey() ([]byte, error) {
    // TODO: generate crypto/rand AES key (e.g., 32 bytes for AES-256)
    return nil, errors.New("not implemented")
}

// Mã hóa file bằng AES-GCM
func EncryptFile(aesKey []byte, plaintext []byte) ([]byte, error) {
    // TODO: AES-GCM encrypt
    return nil, errors.New("not implemented")
}

// Giải mã file AES-GCM
func DecryptFile(aesKey []byte, ciphertext []byte) ([]byte, error) {
    // TODO: AES-GCM decrypt
    return nil, errors.New("not implemented")
}

// Mã hóa AES key bằng public key RSA người nhận
func EncryptAESKeyRSA(aesKey []byte, pubKey *rsa.PublicKey) ([]byte, error) {
    // TODO: RSA-OAEP encrypt
    return nil, errors.New("not implemented")
}

// Giải mã AES key bằng private key RSA
func DecryptAESKeyRSA(encKey []byte, privKey *rsa.PrivateKey) ([]byte, error) {
    // TODO: RSA-OAEP decrypt
    return nil, errors.New("not implemented")
}
