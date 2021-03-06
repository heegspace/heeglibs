package heeglibs

import (
	"fmt"
	"testing"
)

func Test_ExecAction(t *testing.T) {
	msql := NewSqlDB("127.0.0.1", "3306", "apilog", "root", "123456")

	// 插入数据
	msql.ExecAction("INSERT INTO log(id,level,message) VALUE(?,?,?)", func(rows int64, err error) {
		fmt.Println(err)
	}, 6, 5, "use data")

	// 更新数据
	msql.ExecAction("UPDATE log SET level=? WHERE id=5", func(rows int64, err error) {
		fmt.Println(err)
	}, 90)

	return
}

func Test_ExecRows(t *testing.T) {
	msql := NewSqlDB("127.0.0.1", "3306", "apilog", "root", "123456")

	// 查询数据
	var id float64
	var level float64
	var message string
	count, _ := msql.ExecRows("select id,level,message from log", func(rows [][]interface{}, err error) {
		if nil != err {
			fmt.Println(err)

			return
		}

		for _, v := range rows {
			id := v[0].(*float64)
			level := v[1].(*float64)
			message := v[2].(*string)

			fmt.Println(*id, *level, *message)
		}
	}, &id, &level, &message)

	fmt.Println("count: ", count)
}
