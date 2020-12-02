package youyoulibs

import (
	"fmt"
	"testing"
)

func TestGenError(t *testing.T) {
	fmt.Println(GenError("Test Error"))
}

func TestMD5(t *testing.T) {
	fmt.Println(MD5("TestMD5"))
}

func TestGenUUID32(t *testing.T) {
	fmt.Println(GenUUID32())
}

func TestGetNowYmd(t *testing.T) {
	fmt.Println(GetNowYmd())
}

func TestTimestamp(t *testing.T) {
	fmt.Println(Timestamp())
}

func TestFormatDateString(t *testing.T) {
	fmt.Println(FormatDateString("2020-02-02T12:23:23Z"))
}

func TestGetRandomNumberString(t *testing.T) {
	fmt.Println(GetRandomNumberString(6))
}

func TestSendSmsCode(t *testing.T) {
	code,err := SendSmsCode(
		"http://gw.api.taobao.com/router/rest",
		"23285989",
		"SMS_186617096",
		"a8ab447e0ea153c1ba6b7915403e9280", 
		"15920955603",
	)

	if nil != err {
		fmt.Println("Request msm code err: ", err)

		return 
	}

	fmt.Println("Request sms code success, code: ", code)
	return 
}