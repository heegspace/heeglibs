package heeglibs

import (
	"fmt"
	"testing"
	"time"
)

func TestEnCookie(t *testing.T) {
	var token TokenInfo
	token.UID = "12345"
	token.Time = 123123
	token.Token = "asdf2112341234"
	token.Role = 1
	token.Expire = time.Now().Unix()

	cookie, err := EnCookie(token, "12345678901234567890123456789098")
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println("cookie: ", cookie)
	obj, err := DeCookie(cookie, "12345678901234567890123456789098")
	if nil != err {
		fmt.Println(err)
		return
	}

	fmt.Println(obj)
	return
}
