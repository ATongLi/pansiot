/**
 * 加密服务单元测试
 */

package security

import (
	"strings"
	"testing"
)

/**
 * 测试加密解密功能
 */
func TestEncryptDecrypt(t *testing.T) {
	plaintext := []byte("Hello, World!")
	password := "test-password-123"

	// 加密
	ciphertext, err := EncryptData(plaintext, password)
	if err != nil {
		t.Fatalf("EncryptData failed: %v", err)
	}

	// 验证密文不为空
	if ciphertext == "" {
		t.Fatal("ciphertext is empty")
	}

	// 解密
	decrypted, err := DecryptData(ciphertext, password)
	if err != nil {
		t.Fatalf("DecryptData failed: %v", err)
	}

	// 验证明文匹配
	if string(decrypted) != string(plaintext) {
		t.Errorf("decrypted mismatch: got %q, want %q", string(decrypted), string(plaintext))
	}
}

/**
 * 测试错误密码解密
 */
func TestDecryptWrongPassword(t *testing.T) {
	plaintext := []byte("Hello, World!")
	password := "test-password-123"
	wrongPassword := "wrong-password"

	// 加密
	ciphertext, err := EncryptData(plaintext, password)
	if err != nil {
		t.Fatalf("EncryptData failed: %v", err)
	}

	// 使用错误密码解密
	_, err = DecryptData(ciphertext, wrongPassword)
	if err == nil {
		t.Fatal("expected error with wrong password, got nil")
	}

	if err != ErrInvalidPassword {
		t.Errorf("expected ErrInvalidPassword, got %v", err)
	}
}

/**
 * 测试签名生成和验证
 */
func TestSignature(t *testing.T) {
	data := []byte("test data")
	key := "secret-key"

	// 生成签名
	signature := GenerateSignature(data, key)
	if signature == "" {
		t.Fatal("signature is empty")
	}

	// 验证签名
	if !VerifySignature(data, signature, key) {
		t.Fatal("signature verification failed")
	}

	// 错误签名
	wrongSignature := "wrong" + signature
	if VerifySignature(data, wrongSignature, key) {
		t.Fatal("expected verification to fail with wrong signature")
	}
}

/**
 * 测试密码哈希
 */
func TestHashPassword(t *testing.T) {
	password := "my-password-123"

	// 哈希密码
	hash, err := HashPassword(password)
	if err != nil {
		t.Fatalf("HashPassword failed: %v", err)
	}

	// 验证密码
	if !VerifyPassword(password, hash) {
		t.Fatal("password verification failed")
	}

	// 错误密码
	if VerifyPassword("wrong-password", hash) {
		t.Fatal("expected verification to fail with wrong password")
	}
}

/**
 * 测试UUID生成
 */
func TestGenerateUUID(t *testing.T) {
	uuid := GenerateUUID()

	// 验证格式
	if len(uuid) != 36 { // 32个字符 + 4个连字符
		t.Errorf("uuid length mismatch: got %d, want 36", len(uuid))
	}

	// 验证连字符位置
	if uuid[8] != '-' || uuid[13] != '-' || uuid[18] != '-' || uuid[23] != '-' {
		t.Errorf("uuid format invalid: %s", uuid)
	}
}

/**
 * 测试UUID唯一性
 */
func TestGenerateUUIDUniqueness(t *testing.T) {
	uuids := make(map[string]bool)
	iterations := 1000

	for i := 0; i < iterations; i++ {
		uuid := GenerateUUID()
		if uuids[uuid] {
			t.Errorf("duplicate uuid generated: %s", uuid)
		}
		uuids[uuid] = true
	}
}

/**
 * 测试加密性能
 */
func BenchmarkEncrypt(b *testing.B) {
	plaintext := []byte(strings.Repeat("test data ", 100)) // ~1KB
	password := "test-password-123"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := EncryptData(plaintext, password)
		if err != nil {
			b.Fatalf("EncryptData failed: %v", err)
		}
	}
}

/**
 * 测试解密性能
 */
func BenchmarkDecrypt(b *testing.B) {
	plaintext := []byte(strings.Repeat("test data ", 100)) // ~1KB
	password := "test-password-123"

	ciphertext, err := EncryptData(plaintext, password)
	if err != nil {
		b.Fatalf("EncryptData failed: %v", err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := DecryptData(ciphertext, password)
		if err != nil {
			b.Fatalf("DecryptData failed: %v", err)
		}
	}
}
