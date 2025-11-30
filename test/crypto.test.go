package test

import (
	"testing"

	clientinternal "secure_notes/client/internalpkg"
)

// Mã hóa/giải mã: Xác minh mã hóa/giải mã chính xác và đảm bảo khóa mã hóa được bảo vệ.
func TestCrypto_Placeholder(t *testing.T) {
	t.Skip("crypto tests are placeholders until crypto implementations are added")

	// Example future tests:
	// - GenerateAESKey, ensure length
	// - EncryptFile -> DecryptFile roundtrip
	// - EncryptAESKeyRSA -> DecryptAESKeyRSA roundtrip
	_ = clientinternal.GenerateAESKey
	_ = clientinternal.EncryptFile
	_ = clientinternal.DecryptFile
	_ = clientinternal.EncryptAESKeyRSA
	_ = clientinternal.DecryptAESKeyRSA
}
