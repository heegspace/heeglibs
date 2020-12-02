package youyoulibs

import (
	"errors"

	"gopkg.in/mgo.v2"
)

type MongoQuery = *mgo.Query
type Mongo struct {
	Url string

	Session *mgo.Session
}

// 创建新的mongodb实例
//
// @param url 	mongodb的连接地址
//
// @return *Mongo
func NewMongo(url string) *Mongo {
	mongo := &Mongo{
		Url:     url,
		Session: nil,
	}

	return mongo
}

// 打开mongodb连接
func (this *Mongo) Open() (err error) {
	if nil != this.Session {
		err = errors.New("Mongo session already open.")

		return
	}

	if "" == this.Url {
		err = errors.New("Url is error.")

		return
	}

	session, err := mgo.Dial(this.Url)
	if err != nil {
		panic(err)
	}

	// Optional. Switch the session to a monotonic behavior.
	session.SetMode(mgo.Monotonic, true)
	this.Session = session

	return
}

func (this *Mongo) Close() {
	if nil == this.Session {
		return
	}

	defer this.Session.Close()

	return
}

func (this *Mongo) IsOpen() bool {
	if nil == this.Session {
		return false
	}

	return true
}

// 插入数据到数据中
//
// @param db 		数据库
// @param collect	集合
// @param data 		插入的数据
//
// @return err
func (this *Mongo) Insert(db, collect string, data interface{}) (err error) {
	if !this.IsOpen() {
		err = errors.New("Mongo didn't connected.")
		return
	}

	if "" == db || "" == collect {
		err = errors.New("DB or Collect is nil.")
		return
	}

	if nil == data {
		err = errors.New("Insert's Data is nil.")
		return
	}

	err = this.Session.DB(db).C(collect).Insert(data)
	if nil == err {
		return
	}

	return
}

// 根据selector更新数据或插入数据
//
// @param db 		数据库
// @param collect	集合
// @param 	selector 	更新的条件
// @param 	update 		更新的数据
//
// @return err
func (this *Mongo) UpdateOrInsert(db, collect string, selector interface{}, update interface{}) (err error) {
	if !this.IsOpen() {
		err = errors.New("Mongo didn't connected.")
		return
	}

	if "" == db || "" == collect {
		err = errors.New("DB or Collect is nil.")
		return
	}

	if nil == update {
		err = errors.New("Insert's Data is nil.")
		return
	}

	err = this.Session.DB(db).C(collect).Update(selector, update)
	if err != nil {
		switch err {
		case mgo.ErrNotFound:
			err = this.Insert(db, collect, update)
			return
		}
	}

	return
}

// 根据条件查询数据
//
// @param db 		数据库
// @param collect	集合
// @param param 		查询的参数
//
// @return query		返回索引对象
// @return err 		查询状态
func (this *Mongo) Find(db, collect string, param interface{}) (query MongoQuery, err error) {
	if !this.IsOpen() {
		err = errors.New("Mongo didn't connected.")
		return
	}

	if "" == db || "" == collect {
		err = errors.New("DB or Collect is nil.")
		return
	}

	if nil == param {
		err = errors.New("Insert's Data is nil.")
		return
	}

	query = this.Session.DB(db).C(collect).Find(param)

	return
}
