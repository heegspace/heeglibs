package heeglibs

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
)

var TABLE_NAME map[string]interface{}

// 加载数据表名配置文件
// @param file
func LoadTableFromFile(file string) {
	if nil != TABLE_NAME {
		return
	}

	f, err := os.Open(file)
	if err != nil {
		fmt.Println(err)

		return
	}
	defer f.Close()

	data, _ := ioutil.ReadAll(f)

	err = json.Unmarshal([]byte(data), &TABLE_NAME)
	if nil != err {
		fmt.Println(err)
		return
	}

	return
}

// 通过key获取表的真实名字
// @param key 	表的key（主要是所需要表名的全大写）
// @return string
func TableName(key string) string {
	if _, ok := TABLE_NAME[key]; !ok {
		return ""
	}

	return TABLE_NAME[key].(string)
}
