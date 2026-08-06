package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

type Aes256 struct{}

// Aes256GcmEncrypt AES-256-GCM 加密
// key: 32字节密钥
// plainText: 明文字节
// additionalData: 附加校验数据(可选，可传nil)
// 返回: nonce + cipherText + tag 拼接密文
func (aes256 *Aes256) Aes256GcmEncrypt(key, plainText, additionalData []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("aes-256 key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// GCM 模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()

	dst := make([]byte, nonceSize, nonceSize+len(plainText)+gcm.Overhead())

	// 生成12字节随机nonce
	// 注意：这里只读取到 dst 的前 nonceSize 个字节
	if _, err = io.ReadFull(rand.Reader, dst[:nonceSize]); err != nil {
		return nil, err
	}

	cipherWithTag := gcm.Seal(dst[:nonceSize], dst[:nonceSize], plainText, additionalData)
	return cipherWithTag, nil
}

// Aes256GcmDecrypt AES-256-GCM 解密
// key: 32字节密钥
// cipherWithTag: 加密返回的完整字节(nonce+密文+tag)
// additionalData: 加密时传入的附加数据，必须一致
func (aes256 *Aes256) Aes256GcmDecrypt(key, cipherWithTag, additionalData []byte) ([]byte, error) {
	if len(key) != 32 {
		return nil, errors.New("aes-256 key must be 32 bytes")
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonceSize := gcm.NonceSize()
	if len(cipherWithTag) < nonceSize {
		return nil, errors.New("cipher text too short")
	}

	// 拆分nonce和加密数据
	nonce, cipherTextTag := cipherWithTag[:nonceSize], cipherWithTag[nonceSize:]

	// 自动校验tag，篡改/密钥错误直接返回error
	// 💡 提示：Open 的第一个参数传 nil 表示分配新内存存放明文，这是完全正确的
	plainText, err := gcm.Open(nil, nonce, cipherTextTag, additionalData)
	if err != nil {
		return nil, errors.New("decrypt failed, key mismatch or data tampered")
	}
	return plainText, nil
}

// MD5String 计算字符串md5 32位小写
func (aes256 *Aes256) MD5String(s string) string {
	h := md5.New()
	io.WriteString(h, s)
	return fmt.Sprintf("%x", h.Sum(nil))
}

func (aes256 *Aes256) Encrypt(aesKey []byte, plainText string) (string, error) {
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}
	cipherText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return base64.URLEncoding.EncodeToString(cipherText), nil
}

func (aes256 *Aes256) Decrypt(aesKey []byte, cipherBase64 string) (string, error) {
	data, err := base64.URLEncoding.DecodeString(cipherBase64)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(aesKey)
	if err != nil {
		return "", err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("密文太短")
	}
	nonce, cipherText := data[:nonceSize], data[nonceSize:]
	plain, err := gcm.Open(nil, nonce, cipherText, nil)
	if err != nil {
		return "", err
	}
	return string(plain), nil
}
