package aes_128_cdc

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
)

func DecodeBase64(src string) []byte {
	dst, err := base64.StdEncoding.DecodeString(src)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
	}

	return dst
}

func EncodeBase64(src []byte) string {
	dst := base64.StdEncoding.EncodeToString(src)
	return dst
}

//enc
func PswEncrypt(src string, ivParameter string, sKey string) string {
	key := []byte(sKey)
	iv := []byte(ivParameter)
	fmt.Printf("start enc：\nstring: %s\nSkey：%s\niv：%s\n", src, key, iv)

	result, err := Aes128Encrypt([]byte(src), key, iv)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
	}
	return base64.RawStdEncoding.EncodeToString(result)
}

//dec
func PswDecrypt(src string, ivParameter string, sKey string) string {
	fmt.Fprintf(gin.DefaultWriter, "[GIN-debug] start dec\nencstring：%s\nSkey：%s\niv：%s\n", src, sKey, ivParameter)
	key := DecodeBase64(sKey)
	iv := DecodeBase64(ivParameter)

	var result []byte
	var err error

	result, err = base64.StdEncoding.DecodeString(src)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
	}
	origData, err := Aes128Decrypt(result, key, iv)
	if err != nil {
		fmt.Fprintf(gin.DefaultErrorWriter, "[GIN-debug] [ERROR] %v\n", err)
	}
	return string(origData)

}
func Aes128Encrypt(origData, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	origData = PKCS5Padding(origData, blockSize)
	blockMode := cipher.NewCBCEncrypter(block, IV[:blockSize])
	crypted := make([]byte, len(origData))

	blockMode.CryptBlocks(crypted, origData)
	return crypted, nil
}

func Aes128Decrypt(crypted, key []byte, IV []byte) ([]byte, error) {
	if key == nil || len(key) != 16 {
		return nil, nil
	}
	if IV != nil && len(IV) != 16 {
		return nil, nil
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	blockSize := block.BlockSize()
	blockMode := cipher.NewCBCDecrypter(block, IV[:blockSize])
	origData := make([]byte, len(crypted))
	blockMode.CryptBlocks(origData, crypted)
	origData = PKCS5UnPadding(origData)
	return origData, nil
}

func PKCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)
	return append(ciphertext, padtext...)
}

func PKCS5UnPadding(origData []byte) []byte {
	length := len(origData)

	unpadding := int(origData[length-1])
	return origData[:(length - unpadding)]
}
