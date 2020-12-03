package heeglibs
package heeglibs

import (
	"fmt"
	"testing"
)

func Test_ClickHouse(t *testing.T) {
	clchse := NewClickHouse("tcp://127.0.0.1:9000?debug=true")

	// 查询数据
	var id float64
	var level float64
	var message string
	count, _ := clchse.ExecRows("select id,level,message from log", func(rows [][]interface{}, err error) {
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

	// 插入数据
	clchse.ExecAction("INSERT INTO log(id,level,message) VALUE(?,?,?)", func(err error) {
		fmt.Println(err)
	}, 6, 5, "use data")

	// 更新数据
	clchse.ExecAction("UPDATE log SET level=? WHERE id=5", func(err error) {
		fmt.Println(err)
	}, 90)

	return
}
