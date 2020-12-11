package heeglibs

import (
	"database/sql"
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type SqlDB struct {
	UserName string
	PassWord string
	Host     string
	Port     string
	DbName   string

	Db *sql.DB
}

func NewSqlDB(host, port, dbname, username, password string) *SqlDB {
	sql := &SqlDB{
		Host:     host,
		Port:     port,
		DbName:   dbname,
		UserName: username,
		PassWord: password,
	}

	return sql
}

func (this *SqlDB) GetDB() *sql.DB {
	if this.Db == nil {
		var err error

		this.Db, err = this.getdb()
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
// @param args 查询行的临时存储变量
//
// @return count,err
//
func (this *SqlDB) ExecRows(statement string, callback func([][]interface{}, error), args ...interface{}) (count int, err error) {
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
			switch v.(type) {
			case *byte:
				tem := *v.(*byte)
				temp = append(temp, &tem)
			case *[]byte:
				tem := *v.(*[]byte)
				temp = append(temp, &tem)
			case *float32:
				tem := *v.(*float32)
				temp = append(temp, &tem)
			case *float64:
				tem := *v.(*float64)
				temp = append(temp, &tem)
			case *[]float32:
				tem := *v.(*[]float32)
				temp = append(temp, &tem)
			case *[]float64:
				tem := *v.(*[]float64)
				temp = append(temp, &tem)
			case *int:
				tem := *v.(*int)
				temp = append(temp, &tem)
			case *int8:
				tem := *v.(*int8)
				temp = append(temp, &tem)
			case *int16:
				tem := *v.(*int16)
				temp = append(temp, &tem)
			case *int32:
				tem := *v.(*int32)
				temp = append(temp, &tem)
			case *int64:
				tem := *v.(*int64)
				temp = append(temp, &tem)
			case *[]int:
				tem := *v.(*[]int)
				temp = append(temp, &tem)
			case *[]int8:
				tem := *v.(*[]int8)
				temp = append(temp, &tem)
			case *[]int16:
				tem := *v.(*[]int16)
				temp = append(temp, &tem)
			case *[]int32:
				tem := *v.(*[]int32)
				temp = append(temp, &tem)
			case *[]int64:
				tem := *v.(*[]int64)
				temp = append(temp, &tem)
			case *uint16:
				tem := *v.(*uint16)
				temp = append(temp, &tem)
			case *uint32:
				tem := *v.(*uint32)
				temp = append(temp, &tem)
			case *uint64:
				tem := *v.(*uint64)
				temp = append(temp, &tem)
			case *[]uint16:
				tem := *v.(*[]uint16)
				temp = append(temp, &tem)
			case *[]uint32:
				tem := *v.(*[]uint32)
				temp = append(temp, &tem)
			case *[]uint64:
				tem := *v.(*[]uint64)
				temp = append(temp, &tem)
			case *string:
				tem := *v.(*string)
				temp = append(temp, &tem)
			case *[]string:
				tem := *v.(*[]string)
				temp = append(temp, &tem)
			}
		}

		count += 1
		value = append(value, temp)
	}

	callback(value, err)
	return
}

// 执行数据操作动作，主要是插入数据和更新数据
//
// @param statement 	动作的语句
// @param callback 执行的回调函数
// @param args 动作的参数
//
// @return err
//
func (this *SqlDB) ExecAction(statement string, callback func(int64, error), args ...interface{}) (count int64, err error) {
	count = 0

	db := this.GetDB()
	result, err := db.Exec(statement, args...)
	if nil != err {
		callback(0, err)

		return
	}

	rows, _ := result.LastInsertId()
	callback(rows, err)
	return
}

// 插入数据
//
// @param statement 	动作的语句
// @param callback 执行的回调函数
// @param args 动作的参数
//
// @return err
//
func (this *SqlDB) ExecInsert(statement string, callback func(int64, error), args ...interface{}) (count int64, err error) {
	count = 0

	db := this.GetDB()
	result, err := db.Exec(statement, args...)
	if nil != err {
		callback(0, err)

		return
	}

	count, _ = result.RowsAffected()
	lastId, _ := result.LastInsertId()

	callback(lastId, err)
	return
}

// 设置数据库空闲连接数大小
func (this *SqlDB) SetMaxIdleConns(count int) {
	this.Db.SetMaxIdleConns(count)
}

// 最大打开连接数
func (this *SqlDB) SetMaxOpenConns(count int) {
	this.Db.SetMaxOpenConns(count)
}

func (this *SqlDB) enable() bool {
	if 0 == len(this.Host) || 0 == len(this.Port) ||
		0 == len(this.DbName) || 0 == len(this.UserName) ||
		0 == len(this.PassWord) {
		return false
	}

	return true
}

func (this *SqlDB) getdb() (*sql.DB, error) {
	if !this.enable() {
		return nil, errors.New("mysql connect info error.")
	}

	server := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True",
		this.UserName,
		this.PassWord,
		this.Host,
		this.Port,
		this.DbName,
	)

	db, err := sql.Open("mysql", server)
	if err != nil {
		return nil, err
	}
	err = db.Ping()
	if nil != err {
		return nil, err
	}

	this.Db = db
	this.Db.SetMaxIdleConns(20)  //连接池的空闲数大小
	this.Db.SetMaxOpenConns(100) //最大打开连接数

	return this.Db, nil
}
