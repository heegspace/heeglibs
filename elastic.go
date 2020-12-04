package heeglibs

import "fmt"

type ElasticSearch struct {
	Url string
}

// 创建elasticsearch操作对象
// @param url 	elasticsearch集群地址
//
func NewElasticSearch(url string) *ElasticSearch {
	obj := &ElasticSearch{
		Url: url,
	}

	return obj
}

// 创建索引
//
// @param index 	索引名
// @param  attr		索引参数[包含索引属性以及映射等信息]
// @body 	json
//
func (this *ElasticSearch) CreateIndex(index string, attr interface{}) (r []byte, err error) {
	url := this.Url + "/" + index

	request := NewHttpClient()
	err = request.NewRequest(url, attr)
	if nil != err {
		return
	}

	request.Header("Content-Type", "application/json")
	r, err = request.Put()
	if nil != err {
		return
	}

	return
}

// 更新索引，主要用于修改索引的属性以及添加索引的映射
//
// @param index 	索引名
// @param attr 		索引属性[包含索引属性以及映射等信息]
// @body json
// 注： 有些信息是不能被修改的，如:number_of_shards
//
func (this *ElasticSearch) UpdateIndex(index string, field string, attr interface{}) (r []byte, err error) {
	url := this.Url + "/" + index + "/" + field

	request := NewHttpClient()
	err = request.NewRequest(url, attr)
	if nil != err {
		return
	}

	request.Header("Content-Type", "application/json")
	r, err = request.Put()
	if nil != err {
		return
	}

	return
}

// 删除索引
//
// @param index 索引名
//
func (this *ElasticSearch) DeleteIndex(index string) (r []byte, err error) {
	url := this.Url + "/" + index

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Delete()
	if nil != err {
		return
	}

	return
}

// 打开索引
//
// @param index
//
func (this *ElasticSearch) OpenIndex(index string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_open"

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 添加索引别名
//
// @param 	index 	索引名
// @param 	alias 	别名
//
// @return []byte,err
//
func (this *ElasticSearch) AddIndexAlias(index string, alias string) (r []byte, err error) {
	url := this.Url + "/_aliases"

	body := `
		{
			"actions": [
				{
					"add": {
						"index": "%s",
						"alias": "%s"
					}
				}
			]
		}
	`
	body = fmt.Sprintf(body, index, alias)

	request := NewHttpClient()
	err = request.NewRequest(url, body)
	if nil != err {
		return
	}

	request.Header("Content-Type", "application/json")
	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 关闭索引
//
// @param index
//
func (this *ElasticSearch) CloseIndex(index string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_close"

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 获取索引信息
// @param 	index 	索引名
func (this *ElasticSearch) GetIndex(index string) (r []byte, err error) {
	url := this.Url + "/" + index

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Get()
	if nil != err {
		return
	}

	return
}

// 清除索引缓存
// @param 	index 	索引名
func (this *ElasticSearch) ClearCache(index string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_cache/clear"

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 刷新索引
// @param  	index 	索引名
func (this *ElasticSearch) RefreshIndex(index string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_refresh"

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 冲洗索引
// 索引主要是通过执行冲洗将数据保存到索引存储并清除内部事物日志
// @param  	index 	索引名
func (this *ElasticSearch) FlushIndex(index string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_flush"

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 添加文档{添加数据}
// 注： 数据格式和索引中的mapping格式样的
// @param  	index 		索引名
// @param 	doc[json] 	文档数据
func (this *ElasticSearch) PostDoc(index string, doc string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_doc"

	request := NewHttpClient()
	err = request.NewRequest(url, doc)
	if nil != err {
		return
	}

	request.Header("Content-Type", "application/json")
	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 通过id获取文档
// @param 	index 	索引名
// @param 	id 		文档id
func (this *ElasticSearch) GetDoc(index string, id string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_doc/" + id

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Get()
	if nil != err {
		return
	}

	return
}

// 仅仅是获取文档的内容
// @param 	index 	索引名
func (this *ElasticSearch) DocSource(index string, id string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_source/" + id
	fmt.Println(url)

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Get()
	if nil != err {
		return
	}

	return
}

// 通过url中添加搜索信息的方式搜索
// @param 	index 	索引名
// @param 	query 	索引参数
func (this *ElasticSearch) SearchByQuery(index string, query string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_search?" + query

	request := NewHttpClient()
	err = request.NewRequest(url, nil)
	if nil != err {
		return
	}

	r, err = request.Get()
	if nil != err {
		return
	}

	return
}

// 通过body进行搜索
// @param index		索引名
// @param body 		索引参数体
// 注： 可以进行任意的查询
//		包含全文查询、短语、前缀、多字段等等
func (this *ElasticSearch) SearchByBody(index string, body string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_search"

	request := NewHttpClient()
	err = request.NewRequest(url, body)
	if nil != err {
		return
	}

	request.Header("Content-Type", "application/json")
	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 验证查询语句是否正确
// @param 	index 	索引名字
// @param 	body 	查询体
// 注： 查询体只能包含query域
// 正确返回:{"_shards":{"total":1,"successful":1,"failed":0},"valid":true}
// 不能查询返回: {"valid": false}
func (this *ElasticSearch) SearchValidate(index string, body string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_validate/query/"

	request := NewHttpClient()
	err = request.NewRequest(url, body)
	if nil != err {
		return
	}

	request.Header("Content-Type", "application/json")
	r, err = request.Post()
	if nil != err {
		return
	}

	return
}

// 查询满足条件的数据条数
// @param index 	索引名
// @param body 		查询参数体
// 注： 查询体仅仅包含查询参数中的query域
func (this *ElasticSearch) CountByBody(index string, body string) (r []byte, err error) {
	url := this.Url + "/" + index + "/_count"

	request := NewHttpClient()
	err = request.NewRequest(url, body)
	if nil != err {
		return
	}

	request.Header("Content-Type", "application/json")
	r, err = request.Post()
	if nil != err {
		return
	}

	return
}
