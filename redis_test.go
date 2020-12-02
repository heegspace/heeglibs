package youyoulibs

import (
	"fmt"
	"testing"
)

func TestSetValue(t *testing.T) {
	redis := NewRedis("127.0.0.1", "6379", "", 0, 10800)
	redis.Open()

	err := redis.SetValue("test", "success", 0)
	if nil != err {
		fmt.Println("SetValue fail,", err)

		return
	}

	fmt.Println("SetValue success.")

	return
}

func TestGetValue(t *testing.T) {
	redis := NewRedis("127.0.0.1", "6379", "", 0, 10800)
	redis.Open()

	value, err := redis.GetValue("test", false, 0)
	if nil != err {
		fmt.Println("GetValue fail,", err)

		return
	}

	fmt.Println("GetValue success. value: ", value)

	return
}

func TestGetKeys(t *testing.T) {
	redis := NewRedis("127.0.0.1", "6379", "", 0, 10800)
	redis.Open()

	value, err := redis.GetKeys()
	if nil != err {
		fmt.Println("GetKeys fail,", err)

		return
	}

	fmt.Println("GetKeys success. value: ", value)

	return
}
