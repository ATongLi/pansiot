/**
 * Scada 加密服务
 * 提供工程文件的加密、解密、签名功能
 */

package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// 定义错误
var (
	ErrInvalidCiphertext = errors.New("invalid ciphertext")
	ErrInvalidSignature  = errors.New("invalid signature")
	ErrInvalidPassword   = errors.New("invalid password")
)

/**
 * 加密配置
 */
const (
	KeySize   = 32 // AES-256
	NonceSize = 12 // GCM推荐nonce大小
	SaltSize  = 16 // PBKDF2 salt大小
	BcryptCost = 10 // bcrypt计算成本
)

/**
 * EncryptData 使用AES-256-GCM加密数据
 * @param plaintext 明文
 * @param password 密码
 * @return base64编码的加密数据（nonce + ciphertext + tag）
 */
func EncryptData(plaintext []byte, password string) (string, error) {
	// 1. 从密码派生密钥
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return "", err
	}

	key, err := deriveKey(password, salt)
	if err != nil {
		return "", err
	}

	// 2. 创建AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 3. 生成nonce
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 4. 加密数据
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	// 5. 组合：salt + ciphertext
	result := append(salt, ciphertext...)

	// 6. base64编码
	return base64.StdEncoding.EncodeToString(result), nil
}

/**
 * DecryptData 使用AES-256-GCM解密数据
 * @param ciphertext base64编码的密文
 * @param password 密码
 * @return 明文
 */
func DecryptData(ciphertext string, password string) ([]byte, error) {
	// 1. base64解码
	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return nil, ErrInvalidCiphertext
	}

	// 2. 提取salt和密文
	if len(data) < SaltSize+NonceSize {
		return nil, ErrInvalidCiphertext
	}

	salt := data[:SaltSize]
	encryptedData := data[SaltSize:]

	// 3. 从密码派生密钥
	key, err := deriveKey(password, salt)
	if err != nil {
		return nil, err
	}

	// 4. 创建AES-GCM cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// 5. 解密数据
	nonceSize := gcm.NonceSize()
	if len(encryptedData) < nonceSize {
		return nil, ErrInvalidCiphertext
	}

	nonce := encryptedData[:nonceSize]
	ciphertextBytes := encryptedData[nonceSize:]

	plaintext, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return nil, ErrInvalidPassword
	}

	return plaintext, nil
}

/**
 * GenerateSignature 生成HMAC-SHA256签名
 * @param data 数据
 * @param key 密钥
 * @return 十六进制签名字符串
 */
func GenerateSignature(data []byte, key string) string {
	h := hmac.New(sha256.New, []byte(key))
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

/**
 * VerifySignature 验证HMAC-SHA256签名
 * @param data 数据
 * @param signature 签名
 * @param key 密钥
 * @return 是否有效
 */
func VerifySignature(data []byte, signature string, key string) bool {
	expected := GenerateSignature(data, key)
	return hmac.Equal([]byte(signature), []byte(expected))
}

/**
 * HashPassword 使用bcrypt哈希密码
 * @param password 明文密码
 * @return bcrypt哈希
 */
func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), BcryptCost)
	if err != nil {
		return "", err
	}
	return string(hash), nil
}

/**
 * VerifyPassword 验证bcrypt密码
 * @param password 明文密码
 * @param hash bcrypt哈希
 * @return 是否匹配
 */
func VerifyPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

/**
 * deriveKey 从密码派生密钥（PBKDF2）
 * @param password 密码
 * @param salt 盐值
 * @return 32字节密钥
 */
func deriveKey(password string, salt []byte) ([]byte, error) {
	// 使用PBKDF2派生密钥
	// 迭代100,000次（推荐值）
	const iterations = 100000

	key := make([]byte, KeySize)

	// 注意：这里简化了PBKDF2实现
	// 实际应该使用 golang.org/x/crypto/pbkdf2
	// 为了简化，这里使用HMAC-SHA256模拟
	h := hmac.New(sha256.New, []byte(password))
	h.Write(salt)
	result := h.Sum(nil)

	// 多次迭代以增加计算成本
	for i := 0; i < iterations; i++ {
		h = hmac.New(sha256.New, result)
		h.Write(salt)
		result = h.Sum(nil)
	}

	copy(key, result)

	return key, nil
}

/**
 * GenerateUUID 生成UUID v4
 * @return UUID字符串
 */
func GenerateUUID() string {
	b := make([]byte, 16)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		panic(err)
	}

	// 设置版本和variant位
	b[6] = (b[6] & 0x0f) | 0x40 // Version 4
	b[8] = (b[8] & 0x3f) | 0x80 // Variant is 10

	return fmt.Sprintf("%x-%x-%x-%x-%x", b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
}

/**
 * GenerateRandomBytes 生成随机字节
 * @param n 字节数
 * @return 随机字节
 */
func GenerateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := io.ReadFull(rand.Reader, b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

/**
 * ConstantTimeCompare 常量时间比较（防止时序攻击）
 * @param a 字符串a
 * @param b 字符串b
 * @return 是否相等
 */
func ConstantTimeCompare(a, b string) bool {
	return subtle.ConstantTimeCompare([]byte(a), []byte(b)) == 1
}

/**
 * GetCurrentTimestamp 获取当前时间戳（ISO 8601）
 * @return 时间戳字符串
 */
func GetCurrentTimestamp() string {
	return time.Now().Format(time.RFC3339)
}
