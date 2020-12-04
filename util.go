package heeglibs

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"sort"
	"strings"
	"time"
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
