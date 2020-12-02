package heeglibs

import (
	"errors"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	NotFound = gorm.ErrRecordNotFound
)

type SqlDB struct {
	UserName string
	PassWord string
	Host     string
	Port     string
	DbName   string

	Db *gorm.DB
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

func (this *SqlDB) GetDB() *gorm.DB {
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
// @param statement 查询语句
// @param callback 查询回调函数   参数： 查询到的值 和  查询状态
// @param args 查询行的临时存储变量【主要用于查询的列,目前仅仅支持string和float64,也就是要查询的所有列的类型
// @return count,err
func (this *SqlDB) ExecRows(statement string, callback func([][]interface{}, error), args ...interface{}) (count int, err error) {
	db := this.GetDB()
	rows, err := db.Raw(statement).Rows()
	if nil != err {
		callback(nil, err)

		return
	}

	sum := 0
	value := make([][]interface{}, 0)
	for rows.Next() {
		err = rows.Scan(args...)
		if nil != err {
			continue
		}

		temp := make([]interface{}, 0)
		for _, v := range args {
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
			}
		}

		sum = sum + 1
		value = append(value, temp)
	}

	defer rows.Close()
	callback(value, err)

	count = sum

	return
}

// 执行数据操作动作，主要是插入数据和更新数据
// @param statement 	动作的语句
// @param callback 执行的回调函数
// @param args 动作的参数
// @return err
func (this *SqlDB) ExecAction(statement string, callback func(error), args ...interface{}) (err error) {
	db := this.GetDB()
	err = db.Exec(statement, args...).Error
	if nil != err {
		callback(err)

		return
	}

	callback(err)

	return
}

// 设置数据库为调试模式
func (this *SqlDB) LogMode(mode bool) {
	this.Db.LogMode(true)
}

// 设置数据库空闲连接数大小
func (this *SqlDB) SetMaxIdleConns(count int) {
	this.Db.DB().SetMaxIdleConns(count)
}

// 最大打开连接数
func (this *SqlDB) SetMaxOpenConns(count int) {
	this.Db.DB().SetMaxOpenConns(count)
}

func (this *SqlDB) enable() bool {
	if 0 == len(this.Host) || 0 == len(this.Port) ||
		0 == len(this.DbName) || 0 == len(this.UserName) ||
		0 == len(this.PassWord) {
		return false
	}

	return true
}

func (this *SqlDB) getdb() (*gorm.DB, error) {
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

	db, err := gorm.Open("mysql", server)
	if err != nil {
		return nil, err
	}

	this.Db = db
	this.Db.LogMode(true)
	this.Db.DB().SetMaxIdleConns(5)  //连接池的空闲数大小
	this.Db.DB().SetMaxOpenConns(15) //最大打开连接数
	this.Db.SingularTable(true)

	return this.Db, nil
}
