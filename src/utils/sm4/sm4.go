package sm

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/tjfoc/gmsm/sm4"
)

type SM4 struct{}

func (sm4 *SM4) CreateLinkData(Appid string, data any) (string, error) {

	dMapBytes, _ := json.Marshal(data)
	requestBody := string(dMapBytes)
	secretMsg, err := encryptEcb(Appid, requestBody)
	if err != nil {
		return "", err
	}
	return secretMsg, nil
}

func (sm4 *SM4) SM_DecryptEcb(hexKey, cipherText string) (string, error) {
	// 16进制密钥转字节
	key, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("invalid hex key: %v", err)
	}

	// 16进制密文转字节
	cipherData, err1 := hex.DecodeString(cipherText)
	if err1 != nil {
		return "", fmt.Errorf("invalid cipher text: %v", err1)
	}

	// 执行SM4 ECB解密
	plainData, err2 := decryptEcbPadding(key, cipherData)
	if err2 != nil {
		return "", err2
	}

	// 返回UTF-8字符串
	return string(plainData), nil
}

func encryptEcb(hexKey, paramStr string) (string, error) {
	keyData, err := hex.DecodeString(hexKey)
	if err != nil {
		return "", fmt.Errorf("密钥解码失败: %w", err)
	}

	srcData := []byte(paramStr)

	cipherArray, err2 := encrypt_Ecb_Padding(keyData, srcData)
	if err2 != nil {
		return "", fmt.Errorf("加密失败: %w", err)
	}

	cipherText := hex.EncodeToString(cipherArray)

	return cipherText, nil
}

func encrypt_Ecb_Padding(key, src []byte) ([]byte, error) {
	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	blockSize := block.BlockSize()

	padded := pkcs7Pad(src, blockSize)

	ciphertext := make([]byte, len(padded))

	for i := 0; i < len(padded); i += blockSize {
		block.Encrypt(ciphertext[i:i+blockSize], padded[i:i+blockSize])
	}

	return ciphertext, nil
}

func pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

const SM_BlockSizeN = 16

func decryptEcbPadding(key, cipherData []byte) ([]byte, error) {
	if len(key) != SM_BlockSizeN {
		return nil, fmt.Errorf("invalid key length %d, expected %d", len(key), SM_BlockSizeN)
	}

	if len(cipherData) == 0 || len(cipherData)%SM_BlockSizeN != 0 {
		return nil, fmt.Errorf("ciphertext length %d is not multiple of block size", len(cipherData))
	}

	block, err := sm4.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// ECB模式解密 (逐块处理)
	plainData := make([]byte, len(cipherData))
	for i := 0; i < len(cipherData); i += SM_BlockSizeN {
		block.Decrypt(plainData[i:i+SM_BlockSizeN], cipherData[i:i+SM_BlockSizeN])
	}

	// 移除PKCS#7填充
	return removePKCS7Padding(plainData)
}

func removePKCS7Padding(data []byte) ([]byte, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("empty data")
	}

	// 获取填充长度
	paddingLen := int(data[len(data)-1])

	// 验证填充长度有效性
	if paddingLen == 0 || paddingLen > SM_BlockSizeN {
		return nil, fmt.Errorf("invalid padding size: %d", paddingLen)
	}

	// 验证填充内容
	paddingStart := len(data) - paddingLen
	for i := paddingStart; i < len(data); i++ {
		if int(data[i]) != paddingLen {
			return nil, fmt.Errorf("invalid padding byte at position %d", i)
		}
	}

	return data[:paddingStart], nil
}
