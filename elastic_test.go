package heeglibs

import (
	"fmt"
	"testing"
)

// 映射中的"enabled": false用于控制字段是否被elasticsearch分解
// 可以很好的控制某个字段仅仅是存储而已
func TestCreateIndex(t *testing.T) {
	body := `
		{
			"settings": {
				"number_of_shards": 3,
				"number_of_replicas": 2
			},
			"mappings": {
				"_source": {
					"enabled": true
				},
				"dynamic": true,
				"properties": {
					"message": {
						"type": "text",
						"analyzer": "standard"
					},
					"uid": {
						"type": "text"
					}
				}
				
			}
		}
	`

	elastis := NewElasticSearch("http://127.0.0.1:9200")
	r, err := elastis.CreateIndex("test", body)
	if nil != err {
		fmt.Println("TestCreateIndex error: ", err)
		return
	}

	fmt.Println("CreateIndex success! data: ", string(r))
	return
}

func TestCloseIndex(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.CloseIndex("test")
	if nil != err {
		fmt.Println("TestCloseIndex error: ", err)
		return
	}

	fmt.Println("TestCloseIndex success! data: ", string(response))
	return
}

func TestOpenIndex(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.OpenIndex("test")
	if nil != err {
		fmt.Println("TestOpenIndex error: ", err)
		return
	}

	fmt.Println("TestOpenIndex success! data: ", string(response))
	return
}

func TestUpdateIndex(t *testing.T) {
	body := `
		{
			"number_of_replicas": 2
		}
	`

	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.UpdateIndex("test", "_settings", body)
	if nil != err {
		fmt.Println("UpdateIndex error: ", err)
		return
	}

	fmt.Println("UpdateIndex success! data: ", string(response))
	return
}

func TestUpdateIndexMap(t *testing.T) {
	body := `
	{
		"properties": {
			"log": {
				"type": "text",
				"analyzer": "standard"
			}
		}
	}
	`

	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.UpdateIndex("test", "_mapping", body)
	if nil != err {
		fmt.Println("TestUpdateIndexMap error: ", err)
		return
	}

	fmt.Println("TestUpdateIndexMap success! data: ", string(response))
	return
}

func TestAddIndexAlias(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.AddIndexAlias("test", "myalias")
	if nil != err {
		fmt.Println("TestAddIndexAlias error: ", err)
		return
	}

	fmt.Println("TestAddIndexAlias success! data: ", string(response))
	return
}

func TestClearCache(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.ClearCache("test")
	if nil != err {
		fmt.Println("TestClearCache error: ", err)
		return
	}

	fmt.Println("TestClearCache success! data: ", string(response))
	return
}

func TestRefreshIndex(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.RefreshIndex("test")
	if nil != err {
		fmt.Println("TestRefreshIndex error: ", err)
		return
	}

	fmt.Println("TestRefreshIndex success! data: ", string(response))
	return
}

func TestFlushIndex(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.FlushIndex("test")
	if nil != err {
		fmt.Println("TestFlushIndex error: ", err)
		return
	}

	fmt.Println("TestFlushIndex success! data: ", string(response))
	return
}

func TestGetIndex(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.GetIndex("test")
	if nil != err {
		fmt.Println("TestGetIndex error: ", err)
		return
	}

	fmt.Println("TestGetIndex success! data: ", string(response))
	return
}

func TestPostDoc(t *testing.T) {
	doc := `
		{
			"log": "This is a test"
		}
	`

	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.PostDoc("test", doc)
	if nil != err {
		fmt.Println("TestPostDoc error: ", err)
		return
	}

	fmt.Println("TestPostDoc success! data: ", string(response))
	return
}

func TestGetDoc(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.GetDoc("test", "1")
	if nil != err {
		fmt.Println("TestGetDoc error: ", err)
		return
	}

	fmt.Println("TestGetDoc success! data: ", string(response))
	return
}

func TestDocSource(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.DocSource("index", "ktpDYXMBskf7Xeau3ZnL")
	if nil != err {
		fmt.Println("TestDocSource error: ", err)
		return
	}

	fmt.Println("TestDocSource success! data: ", string(response))
	return
}

func TestSearchByQuery(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.SearchByQuery("index", "1")
	if nil != err {
		fmt.Println("TestSearchByQuery error: ", err)
		return
	}

	fmt.Println("TestSearchByQuery success! data: ", string(response))
	return
}

// from / size 分页
// sort 对结果排序,字符串是不支持排序的
// _source 用于过滤
// query中的是模糊查找的，仅仅是包含
// highlight 高亮显示
func TestSearchByBody(t *testing.T) {
	body := `
		{
			"from": 0,
			"size": 20,
			"sort": [
				{
					"message": {
						"order": "asc"
					}
				}
			],
			"query": {
				"term": {
					"data.Timu": "搜索"
				}
			},
			"highlight": {
				"fields": {
				  "message": {} 
				}
			  },
			"_source": true
		}
	`
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.SearchByBody("topic", body)
	if nil != err {
		fmt.Println("TestSearchByBody error: ", err)
		return
	}

	fmt.Println("TestSearchByBody success! data: ", string(response))
	return
}

// 全文查找
// 其中match中的query跟分词器有关系
/*
单个字段匹配查找：
	"match": {
		"message": {
				"query": "this china",
				"operator": "or"
		}
	}

多个字段匹配查找：
	"multi_match": {
		"query": "资源管理 目录索引",
		"type": "most_fields",
		"fields": ["data.Timu","data.Daan.Name"]
	}
*/
func TestFullTestSearch(t *testing.T) {
	body := `
		{
			"from": 0,
			"size": 20,
			"query": {
			"multi_match": {
				"query": "地球和地球仪",
				"type": "most_fields",
				"operator": "or",
				"fields": ["chapter_name","chapter_gd","source_name", "ti_xing_name","data.Timu","data.Daan.Name"]
			}
			},
			"highlight": {
			"fields": {
				"data.Timu": {},
				"chapter_name": {},
				"chapter_gd": {},
				"ti_xing_name": {},
				"data.Daan.Name": {},
				"source_name": {}
			}
			},
			"_source": true
		}
	`
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.SearchByBody("topic", body)
	if nil != err {
		fmt.Println("TestFullTestSearch error: ", err)
		return
	}

	fmt.Println("TestFullTestSearch success! data: ", string(response))
	return
}

func TestSearchValidate(t *testing.T) {
	body := `
		{
			"query": {
				"term": {
					"message": "this"
				}
			}
		}
	`
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.SearchValidate("index", body)
	if nil != err {
		fmt.Println("TestSearchValidate error: ", err)
		return
	}

	fmt.Println("TestSearchValidate success! data: ", string(response))
	return
}

func TestCountByBody(t *testing.T) {
	body := `
		{
			"query": {
				"multi_match": {
					"query": "资源管理 目录索引",
					"type": "most_fields",
					"fields": ["data.Timu","data.Daan.Name"]
				}
			}
		}
	`
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.CountByBody("topic", body)
	if nil != err {
		fmt.Println("TestCountByBody error: ", err)
		return
	}

	fmt.Println("TestCountByBody success! data: ", string(response))
	return
}

func TestDeleteIndex(t *testing.T) {
	elastis := NewElasticSearch("http://127.0.0.1:9200")
	response, err := elastis.DeleteIndex("test")
	if nil != err {
		fmt.Println("TestDeleteIndex error: ", err)
		return
	}

	fmt.Println("TestDeleteIndex success! data: ", string(response))
	return
}
