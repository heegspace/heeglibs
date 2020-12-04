package heeglibs

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"encoding/gob"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

// 用于设置cookie的结构
type TokenInfo struct {
	UID    string
	Time   int64
	Token  string
	Role   int64
	Expire int64
}

func (token TokenInfo) String() string {
	str := "UID: %s Time: %d Token: %s  Role: %d  Expire: %d"

	return fmt.Sprintf(str, token.UID, token.Time, token.Token, token.Role, token.Expire)
}

var ivspec = []byte("0000000000000000")

func pkCS5Padding(ciphertext []byte, blockSize int) []byte {
	padding := blockSize - len(ciphertext)%blockSize
	padtext := bytes.Repeat([]byte{byte(padding)}, padding)

	return append(ciphertext, padtext...)
}
func pkCS5Trimming(encrypt []byte) []byte {
	padding := encrypt[len(encrypt)-1]

	return encrypt[:len(encrypt)-int(padding)]
}

func aesEncode(src, key string) (value string, err error) {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
		return
	}
	if src == "" {
		fmt.Println("plain content empty")
		err = errors.New("plain content empty")
		return
	}
	ecb := cipher.NewCBCEncrypter(block, ivspec)
	content := []byte(src)
	content = pkCS5Padding(content, block.BlockSize())
	crypted := make([]byte, len(content))
	ecb.CryptBlocks(crypted, content)

	value = hex.EncodeToString(crypted)

	return
}

func aesDecode(token, key string) (subject string, err error) {
	crypted, err := hex.DecodeString(strings.ToLower(token))
	if err != nil || len(crypted) == 0 {
		fmt.Println("plain content empty")
		err = errors.New("plain content empty")

		return
	}

	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		fmt.Println("key error1", err)
		err = errors.New("key error1")

		return
	}

	ecb := cipher.NewCBCDecrypter(block, ivspec)
	decrypted := make([]byte, len(crypted))
	ecb.CryptBlocks(decrypted, crypted)

	subject = string(pkCS5Trimming(decrypted))

	return
}

// 使用密钥key对Token数据进行加密
//
// @param src	编码的结构数据
// @param key 	编码的秘钥
//
// @return string,error
//
func EnCookie(src TokenInfo, key string) (token string, err error) {
	if 0 == len(src.UID) {
		err = errors.New("Uid is empty")
		return
	}

	if 0 == len(src.Token) {
		err = errors.New("Token src is empty.")
		return
	}

	buf := &bytes.Buffer{}
	enc := gob.NewEncoder(buf)
	err = enc.Encode(src)
	if nil != err {
		fmt.Println("wirte err ", err)
		return
	}

	data := buf.Bytes()
	token, err = aesEncode(string(data), key)
	if nil != err {
		return
	}
	return
}

// 使用秘钥key从src中解码cookie
//
// @param src 	解码数据
// @param key 	秘钥key
// @return TokenInfo,error
//
func DeCookie(src string, key string) (token TokenInfo, err error) {
	if 0 == len(src) {
		err = errors.New("Token info format error.")
		return
	}

	obj, err := aesDecode(src, key)
	if nil != err {
		return
	}

	buf := bytes.NewBufferString(obj)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(&token)
	if nil != err {
		return
	}

	return
}
