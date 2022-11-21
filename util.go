package heeglibs

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"encoding/pem"
	"errors"
	"fmt"
	"io"
	"math/big"
	"net"
	"sort"
	"strconv"
	"strings"
	"time"

	pkcs7 "github.com/mergermarket/go-pkcs7"
)

// 根据提供的msg产生一个error对象
// msg string
func GenError(msg string) error {
	return errors.New(msg)
}

// 对字符串进行MD5签名并返回
// str string
func MD5(str string) string {
	m5 := md5.New()
	io.WriteString(m5, str)
	return fmt.Sprintf("%x", m5.Sum(nil))
}

func ObjConvert(src, dst interface{}) (err error) {
	b, err := json.Marshal(src)
	if err != nil {
		return
	}
	err = json.Unmarshal(b, &dst)
	return
}

// base64解码
//
// param 	buf
// return {string}
//
func DeBase64(buf string) string {
	basestring, _ := base64.StdEncoding.DecodeString(buf)

	return string(basestring)
}

// base64编码
//
// param 	buf 	数据缓存
// param 	n		数据长度
// return string
//
func EnBase64(buf []byte, n int) string {
	size := n
	if n > len(buf) {
		size = len(buf)
	}

	sourcestring := base64.StdEncoding.EncodeToString(buf[:size])

	return sourcestring
}

// 产生一个32位的UID字符串
func GenUUID32() (string, error) {
	uuid := make([]byte, 16)
	n, err := rand.Read(uuid)
	switch n != len(uuid) || err != nil {
	case true:
		return "", err
	}

	uuid[8] = 0x80
	uuid[4] = 0x40

	return hex.EncodeToString(uuid), nil
}

// 获取当天的日期,格式是"yyyy-mm-dd"
func GetNowYmd() string {
	return time.Now().Format("2006-01-02")
}

// 获取当前的日期和时间，格式是"yyyy-mm-dd hh:mm:SS"
func GetNowYmdHms() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

// 获取当前时间戳
func Timestamp() int64 {
	return time.Now().Unix()
}

// 格式化当前的时间将其中T和Z替换
func FormatDateString(str string) string {
	s := strings.Replace(str, "T", " ", -1)
	s = strings.Replace(s, "Z", "", -1)

	return s
}

// 获取n位的随机数字字符串
func GetRandomNumberString(n int) string {
	alphanum := "0123456789"
	var bytes = make([]byte, n)
	rand.Read(bytes)
	for i, b := range bytes {
		bytes[i] = alphanum[b%byte(len(alphanum))]
	}
	return string(bytes)
}

func mapMD5(secret string, param map[string]interface{}) string {
	keys := make([]string, 0)
	for k, _ := range param {
		keys = append(keys, k)
	}

	signStr := ""
	sort.Strings(keys)
	for _, value := range keys {
		signStr = signStr + value + param[value].(string)
	}

	signStr = secret + signStr + secret

	signStr = MD5(signStr)

	return strings.ToUpper(signStr)
}

// AES加密数据
//
// @param origInData	需要加密的数据
// @param keyIn		 	加密的key (32byte)
// @return encrypt 		加密以后的数据
// @return err
//
func AesEncode(origInData, keyIn string) (encrypt string, err error) {
	key := []byte(keyIn)
	plainText := []byte(origInData)
	plainText, err = pkcs7.Pad(plainText, aes.BlockSize)
	if err != nil {
		return
	}
	if len(plainText)%aes.BlockSize != 0 {
		err = fmt.Errorf(`plainText: "%s" has the wrong block size`, plainText)
		return
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]
	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
		return
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(cipherText[aes.BlockSize:], plainText)
	encrypt = hex.EncodeToString(cipherText)
	return
}

// AES解密数据
//
// @param crypted 		需要解密的数据
// @param keyIn 		解密数据的key(32byte)
// @return data 		解密以后的数据
// @return err
//
func AesDecode(crypted, keyIn string) (data string, err error) {
	key := []byte(keyIn)
	cipherText, _ := hex.DecodeString(crypted)
	block, err := aes.NewCipher(key)
	if err != nil {
		return
	}
	if len(cipherText) < aes.BlockSize {
		err = errors.New("cipherText too short")

		return
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]
	if len(cipherText)%aes.BlockSize != 0 {
		err = errors.New("cipherText is not a multiple of the block size")

		return
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(cipherText, cipherText)
	cipherText, _ = pkcs7.Unpad(cipherText, aes.BlockSize)
	data = fmt.Sprintf("%s", cipherText)

	return
}

// 使用公钥RSA签名
//
// @param origData	要签名的数据
// @param publicKey 公钥
// @return {sign},{err}
//
func RsaEncrypt(origData []byte, publicKey []byte) ([]byte, error) {
	block, _ := pem.Decode(publicKey)
	if block == nil {
		return nil, errors.New("public key error")
	}
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	pub := pubInterface.(*rsa.PublicKey)
	return rsa.EncryptPKCS1v15(rand.Reader, pub, origData)
}

// 使用对应的私钥解密
//
// @param ciphertext 	加密的数据
// @param privateKey 	私钥
// @return {data},{err}
//
func RsaDecrypt(ciphertext []byte, privateKey []byte) ([]byte, error) {
	block, _ := pem.Decode(privateKey)
	if block == nil {
		return nil, errors.New("private key error!")
	}
	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}
	return rsa.DecryptPKCS1v15(rand.Reader, priv, ciphertext)
}

// ip地址转到整数
//
func IpToInt64(ip string) int64 {
	IP := net.ParseIP(ip)
	if nil == IP {
		return 0
	}

	if nil != IP.To4() {
		bits := strings.Split(ip, ".")
		if 4 > len(bits) {
			return 0
		}

		b0, _ := strconv.Atoi(bits[0])
		b1, _ := strconv.Atoi(bits[1])
		b2, _ := strconv.Atoi(bits[2])
		b3, _ := strconv.Atoi(bits[3])

		b0 = b0 << 24
		b0 += b1 << 16
		b0 += b2 << 8
		b0 += b3

		return int64(b0)
	}

	if nil != IP.To16() {
		b0 := big.NewInt(0)
		b0.SetBytes(IP.To16())

		return b0.Int64()
	}

	return 0
}

// 发送短信验证码
//
// url 		请求的utl接口
// appkey 	请求发送短信的appkey
// template 请求发送短信的template
// secret 	请求发送短信的secret
// mobile 	接收短信的号码
//
// code,err
func SendSmsCode(url, appkey, template, secret, mobile string) (code string, err error) {
	code = GetRandomNumberString(6)

	sms_param := make(map[string]interface{})
	sms_param["code"] = code
	b, err := json.Marshal(sms_param)
	if err != nil {
		return
	}

	param := make(map[string]interface{})
	param["method"] = "alibaba.aliqin.fc.sms.num.send"
	param["app_key"] = appkey
	param["timestamp"] = GetNowYmdHms()
	param["format"] = "json"
	param["v"] = "2.0"
	param["partner_id"] = "apidoc"
	param["sign_method"] = "md5"
	param["sms_template_code"] = template
	param["extend"] = "123456"
	param["rec_num"] = mobile
	param["sms_type"] = "normal"
	param["sms_free_sign_name"] = "elearing注册"
	param["sms_param"] = string(b)

	param["sign"] = mapMD5(secret, param)

	body := ""
	for k, value := range param {
		body = body + k + "=" + value.(string) + "&"
	}

	client := NewHttpClient()
	err = client.NewRequest(url, body)
	if nil != err {
		fmt.Println("NewRequest Fail, ", err)

		return
	}

	client.Header("Content-Type", "application/x-www-form-urlencoded;charset=utf-8")
	res, err := client.Post()
	if nil != err {
		fmt.Println("Http Post Fail, ", err)

		return
	}

	var response map[string]interface{}
	err = json.Unmarshal(res, &response)
	if nil != err {
		return
	}

	if _, ok := response["alibaba_aliqin_fc_sms_num_send_response"]; ok {
		return
	}

	err_res := response["error_response"].(map[string]interface{})

	err = GenError(err_res["msg"].(string))
	code = ""

	return
}
