package wsclient

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/hmac"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"time"
)

func (f *FnOsWsBase) GetReqId() string {
	// 1. 获取秒级时间戳，转为8位16进制字符串（不足左侧补0）
	secTimestamp := time.Now().Unix()
	n := fmt.Sprintf("%08x", secTimestamp)

	// 2. 对应Python中的e = f"{0:04x}"，固定为0的4位16进制
	e := fmt.Sprintf("%04x", 0)

	// 3. 获取毫秒级时间戳（十进制）
	milliTimestamp := time.Now().UnixMilli()

	// 4. 拼接生成req_id，保持和Python完全一致的格式
	reqID := fmt.Sprintf("%s000%d%s", n, milliTimestamp, e)

	return reqID
}

// PKCS7 填充
func (f *FnOsWsBase) pkcs7Pad(data []byte, blockSize int) []byte {
	padding := blockSize - len(data)%blockSize
	padText := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(data, padText...)
}

// AES-CBC 加密并返回 Base64（修复后）
func (f *FnOsWsBase) aesCbcEncryptBase64(plaintext []byte, key string, iv []byte) (string, error) {
	// 1. 校验 IV 长度
	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("IV 长度必须为 %d 字节（当前：%d）", aes.BlockSize, len(iv))
	}

	// 2. 转换 Key 类型并校验长度
	keyBytes := []byte(key)
	keyLen := len(keyBytes)
	if keyLen != 16 && keyLen != 24 && keyLen != 32 {
		return "", fmt.Errorf("AES Key 长度必须为 16/24/32 字节（当前：%d）", keyLen)
	}

	// 3. 创建 AES 加密块
	block, err := aes.NewCipher(keyBytes) // 修复：传入 []byte 类型的 key
	if err != nil {
		return "", fmt.Errorf("创建 AES 加密块失败：%w", err)
	}

	// 4. PKCS7 补位
	padded := f.pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))

	// 5. CBC 模式加密
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	// 6. 加密结果转 Base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSA PKCS1 v1.5 加密，返回 Base64（修复后）
func (f *FnOsWsBase) rsaEncrypt(plaintext string, publicKeyPEM []byte) (string, error) {
	// 1. 解码 PEM 格式公钥
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil {
		return "", errors.New("解析 PEM 公钥失败：空的 PEM 块")
	}

	// 2. 修复：公钥类型判断（&& 改为 ||）
	if block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY" { // 此处是 &&，因为要排除两种都不是的情况
		return "", fmt.Errorf("无效的公钥类型（当前：%s），仅支持 PUBLIC KEY/PKCS1 PUBLIC KEY", block.Type)
	}

	// 3. 解析公钥（兼容 PKCS#8 和 PKCS#1 格式）
	var pub interface{}
	var err error
	if block.Type == "PUBLIC KEY" {
		// PKCS#8 格式
		pub, err = x509.ParsePKIXPublicKey(block.Bytes)
	} else {
		// PKCS#1 格式
		pub, err = x509.ParsePKCS1PublicKey(block.Bytes)
	}
	if err != nil {
		return "", fmt.Errorf("解析公钥失败：%w", err)
	}

	// 4. 断言为 RSA 公钥
	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", errors.New("公钥类型错误：非 RSA 公钥")
	}

	// 5. 转换明文类型 + 校验长度（PKCS1 v1.5 明文长度 ≤ 模数长度-11）
	plaintextBytes := []byte(plaintext) // 修复：转为 []byte
	maxPlaintextLen := rsaPub.Size() - 11
	if len(plaintextBytes) > maxPlaintextLen {
		return "", fmt.Errorf("明文长度超出限制（最大：%d 字节，当前：%d）", maxPlaintextLen, len(plaintextBytes))
	}

	// 6. RSA PKCS1 v1.5 加密
	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, plaintextBytes)
	if err != nil {
		return "", fmt.Errorf("RSA 加密失败：%w", err)
	}

	// 7. 加密结果转 Base64
	return base64.StdEncoding.EncodeToString(cipherBytes), nil
}

// HMAC-SHA256 + Base64
func (f *FnOsWsBase) hmacSha256Base64(message string, base64Key string) (string, error) {
	keyBytes, err := base64.StdEncoding.DecodeString(base64Key)
	if err != nil {
		return "", err
	}

	mac := hmac.New(sha256.New, keyBytes)
	mac.Write([]byte(message))
	signature := mac.Sum(nil)

	return base64.StdEncoding.EncodeToString(signature), nil
}

func ConvertTo[T any](data any) (T, error) {
	var ret T

	b, err := json.Marshal(data)
	if err != nil {
		return ret, err
	}

	err = json.Unmarshal(b, &ret)
	return ret, err
}

func generateKey() string {
	// 1. 生成16字节安全随机数
	keyBytes := make([]byte, 16)
	_, err := rand.Read(keyBytes)
	if err != nil {
		panic(err)
	}
	// 2. 转为十六进制字符串（小写，和Python token_hex一致）
	return hex.EncodeToString(keyBytes)
}

func generateIv() []byte {
	ivBytes := make([]byte, 16)
	_, err := rand.Read(ivBytes)
	if err != nil {
		panic(err)
	}
	return ivBytes
}
