package heeglibs

import (
	"fmt"
	"testing"
	"time"

	log "github.com/sirupsen/logrus"
)

type ExampleTest struct {
	CountryCode string
	OsId        int64
	BrowserId   int64
	Categories  []int16
	ActionDay   string
	ActionTime  string
}

// 测试创建数据表
//
func Test_ExecTable(t *testing.T) {
	clchse := NewClickHouse("tcp://127.0.0.1:9000?username=default&password=576188&debug=true")

	clchse.ExecTable(`
		CREATE TABLE example_test (
			country_code FixedString(2),
			os_id        UInt8,
			browser_id   UInt8,
			categories   Array(Int16),
			action_day   Date,
			action_time  DateTime
		) engine=Memory
	`, func(err error) {
		fmt.Println("Create: ", err)
	})

	return
}

// 测试数据表操作
// 插入数据或更新
//
func Test_ExecAction(t *testing.T) {
	clchse := NewClickHouse("tcp://127.0.0.1:9000?username=default&password=576188&debug=true")

	// 插入数据
	clchse.ExecAction("INSERT INTO example_test(country_code, os_id, browser_id, categories, action_day, action_time) VALUES(?, ?, ?, ?, ?, ?)",
		func(rows int64, err error) {
			fmt.Println("Insert: ", err, rows)
		}, "CN", 99, 99, []int16{1, 2, 3}, time.Now(), time.Now(),
	)

	// 更新数据
	clchse.ExecAction("UPDATE example_test SET os_id=? WHERE country_code=99", func(rows int64, err error) {
		fmt.Println(rows, err)
	}, 200)

	return
}

// 测试数据查询
func Test_ExecRows(t *testing.T) {
	clchse := NewClickHouse("tcp://127.0.0.1:9000?username=default&password=576188&debug=true")

	// 查询数据
	var test ExampleTest
	count, _ := clchse.ExecRows("SELECT country_code, os_id, browser_id, categories, action_day, action_time FROM example_test", func(rows [][]interface{}, err error) {
		if nil != err {
			fmt.Println(err)

			return
		}

		// 将查询到的值进行转换
		for _, v := range rows {
			var item ExampleTest
			item.CountryCode = *(v[0].(*string))
			item.OsId = *(v[1].(*int64))
			item.BrowserId = *(v[2].(*int64))
			item.Categories = *(v[3].(*[]int16))
			item.ActionDay = *(v[4].(*string))
			item.ActionTime = *(v[4].(*string))

			log.Println(item)
		}
	}, &test.CountryCode, &test.OsId, &test.BrowserId, &test.Categories, &test.ActionDay, &test.ActionTime)

	fmt.Println("count: ", count)
	log.Println(test)

	return
}
