package heeglibs

// clickhouse客户端操作

import (
	"database/sql"
	"errors"
	"fmt"
	"reflect"

	"github.com/ClickHouse/clickhouse-go"
	_ "github.com/ClickHouse/clickhouse-go"
	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

type ClickHouse struct {
	Host string

	Db *sql.DB
}

// 创建clickhouse实例
//
// @param host 	clickhouse数据库连接地址
//
func NewClickHouse(host string) *ClickHouse {
	obj := &ClickHouse{
		Host: host,
	}

	return obj
}

// 获取对应的db对象
//
func (this *ClickHouse) GetDB() *sql.DB {
	if this.Db == nil {
		var err error

		err = this.getdb()
		if err != nil {
			panic("get Db err:" + err.Error())
		}

		return this.Db
	}

	return this.Db
}

// 查询接口
//
// @param statement 查询语句
// @param callback 查询回调函数   参数： 查询到的值 和  查询状态
// @param args 查询行的临时存储变量【主要用于查询的列,目前仅仅支持string和float64,也就是要查询的所有列的类型
// @return count,err
//
func (this *ClickHouse) ExecRows(statement string, callback func([][]interface{}, error), args ...interface{}) (count int, err error) {
	db := this.GetDB()
	rows, err := db.Query(statement)
	if nil != err {
		callback(nil, err)

		return
	}
	defer rows.Close()

	count = 0
	value := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(args...)
		if nil != err {
			continue
		}

		temp := make([]interface{}, 0)
		for _, v := range args {
			log.Println(reflect.TypeOf(v))
			switch v.(type) {
			case *string:
				tem := *v.(*string)
				temp = append(temp, &tem)
			case *float64:
				tem := *v.(*float64)
				temp = append(temp, &tem)
			case *int64:
				tem := *v.(*int64)
				temp = append(temp, &tem)
			case *[]int16:
				tem := *v.(*[]int16)
				temp = append(temp, &tem)
			}
		}

		count += 1
		value = append(value, temp)
	}

	callback(value, err)
	return
}

func (this *ClickHouse) ExecInsert(statement string, callback func(error), args ...interface{}) (err error) {
	tx, err := this.GetDB().Begin()
	if nil != err {
		callback(err)

		return
	}

	stmt, err := tx.Prepare(statement)
	if nil != err {
		callback(err)

		return
	}

	rows, err := stmt.Exec(args...)
	if nil != err {
		callback(err)

		return
	}

	tx.Commit()
	log.Println(statement, rows)
	callback(nil)
	return
}

// 执行数据操作动作，主要是插入数据和更新数据
//
// @param statement 	动作的语句
// @param callback 执行的回调函数
// @param args 动作的参数
// @return err
//
func (this *ClickHouse) ExecAction(statement string, callback func(error), args ...interface{}) (err error) {
	db := this.GetDB()
	result, err := db.Exec(statement, args...)
	if nil != err {
		callback(err)

		return
	}

	rows, err := result.RowsAffected()
	if nil != err {
		return
	}

	log.Println(statement, ", rows: ", rows)

	callback(err)
	return
}

// 设置数据库空闲连接数大小
//
func (this *ClickHouse) SetMaxIdleConns(count int) {
	this.Db.SetMaxIdleConns(count)
}

// 最大打开连接数
//
func (this *ClickHouse) SetMaxOpenConns(count int) {
	this.Db.SetMaxOpenConns(count)
}

func (this *ClickHouse) getdb() error {
	if 0 == len(this.Host) {
		return errors.New("mysql connect info error.")
	}

	connect, err := sql.Open("clickhouse", this.Host)
	if err != nil {
		return err
	}

	err = connect.Ping()
	if nil != err {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}

		return err
	}

	this.Db = connect
	return nil
}
