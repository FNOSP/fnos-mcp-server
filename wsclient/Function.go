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
	"encoding/json"
	"encoding/pem"
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

// AES-CBC 加密并返回 Base64
func (f *FnOsWsBase) aesCbcEncryptBase64(plaintext, key, iv []byte) (string, error) {
	if len(iv) != aes.BlockSize {
		return "", fmt.Errorf("IV 长度必须为 16 字节")
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	padded := f.pkcs7Pad(plaintext, aes.BlockSize)
	ciphertext := make([]byte, len(padded))

	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(ciphertext, padded)

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// RSA PKCS1 v1.5 加密，返回 Base64
func (f *FnOsWsBase) rsaEncrypt(plaintext []byte, publicKeyPEM []byte) (string, error) {
	block, _ := pem.Decode(publicKeyPEM)
	if block == nil || block.Type != "PUBLIC KEY" && block.Type != "RSA PUBLIC KEY" {
		return "", fmt.Errorf("无效的公钥 PEM")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		// 兼容 PKCS#1 公钥格式
		pubRsa, err2 := x509.ParsePKCS1PublicKey(block.Bytes)
		if err2 != nil {
			return "", fmt.Errorf("解析公钥失败: %v / %v", err, err2)
		}
		pub = pubRsa
	}

	rsaPub, ok := pub.(*rsa.PublicKey)
	if !ok {
		return "", fmt.Errorf("公钥类型错误")
	}

	cipherBytes, err := rsa.EncryptPKCS1v15(rand.Reader, rsaPub, plaintext)
	if err != nil {
		return "", err
	}

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
