package heeglibs

import (
	"context"

	"github.com/tsuna/gohbase"
	"github.com/tsuna/gohbase/hrpc"
)

type HBase struct {
	Client      gohbase.Client
	AdminClient gohbase.AdminClient
}

// HBase查询的值的对象
type HBaseResult struct {
	RowKey    string // rowkey
	Family    string // 族名
	Qualifier string // 列名
	Value     []byte // 值
}

func NewHBase() *HBase {
	obj := &HBase{}

	return obj
}

// 打开hbase普通用户端
//
// @param host 	集群的主机地址列表[192.168.0.119:2181,192.169.0.120:2181]
//
func (this *HBase) OpenClient(host string) {
	this.Client = gohbase.NewClient(host)

	return
}

// 通过rowkey查询一整行数据
//
// @param table_name 	表名
// @param rowkey 		表的rowkey
//
// @return []HBaseResult,err
//
func (this *HBase) GetRow(table_name string, rowkey string) (r []HBaseResult, err error) {
	getRequest, err := hrpc.NewGetStr(context.Background(), table_name, rowkey)
	if nil != err {
		return
	}

	getRsp, err := this.Client.Get(getRequest)
	if nil != err {
		return
	}

	r = make([]HBaseResult, 0)
	for _, v := range getRsp.Cells {
		var temp HBaseResult
		temp.RowKey = string(v.Row[:])
		temp.Family = string(v.Family[:])
		temp.Qualifier = string(v.Qualifier[:])
		temp.Value = make([]byte, len(v.Value[:]))
		copy(temp.Value, v.Value[:])

		r = append(r, temp)
	}

	return
}

// 插入或者更新一个值
//
// @param table_name 	表名
// @param value 			需要插入或者更新的值
//
// @return err
//
func (this *HBase) PutCell(table_name string, value HBaseResult) (err error) {
	values := map[string]map[string][]byte{
		value.Family: map[string][]byte{
			value.Qualifier: value.Value,
		},
	}

	putRequest, err := hrpc.NewPutStr(context.Background(), table_name, value.RowKey, values)
	if nil != err {
		return
	}

	_, err = this.Client.Put(putRequest)
	if nil != err {
		return
	}

	return
}

// 根据列名获取值
//
// @param  table_name 	表名
// @param  rowkey 		rowkey
// @param  family 		族名
// @param  columns 		列明
//
// @return map[string]HBaseResult,err
//
func (this *HBase) GetCell(table_name, rowkey, family string, columns []string) (r map[string]HBaseResult, err error) {
	faly := map[string][]string{family: columns}
	getRequest, err := hrpc.NewGetStr(context.Background(), table_name, rowkey,
		hrpc.Families(faly))
	if nil != err {
		return
	}

	getRsp, err := this.Client.Get(getRequest)
	if nil != err {
		return
	}

	ret := make(map[string]HBaseResult)
	for _, v := range getRsp.Cells {
		var temp HBaseResult
		temp.RowKey = string(v.Row[:])
		temp.Family = string(v.Family[:])
		temp.Qualifier = string(v.Qualifier[:])
		temp.Value = make([]byte, len(v.Value[:]))
		copy(temp.Value, v.Value[:])

		ret[temp.Qualifier] = temp
	}

	r = ret
	return
}

// 删除Hbase中对应的数据，可以删除族，或者删除某个族中的列值
//
// @param table_name 	表名
// @param value 			删除的值
//	注： 如果删除的是族，只需要将value中的Qualifier设置为""
// 	注： 如果删除的是某个列，只需要将value中的value设置为nil
//
// @return err
//
func (this *HBase) DelCell(table_name string, value HBaseResult) (err error) {
	if "" == value.Qualifier {
		// 删除族
		family := map[string]map[string][]byte{
			value.Family: nil,
		}

		_, err = hrpc.NewDelStr(context.Background(), table_name, value.RowKey,
			family)
		if nil != err {
			return
		}

		return
	}

	// 删除某个列
	family := map[string]map[string][]byte{
		value.Family: map[string][]byte{
			value.Qualifier: nil,
		},
	}

	_, err = hrpc.NewDelStr(context.Background(), table_name, value.RowKey,
		family)
	if nil != err {
		return
	}

	return
}

// 打开hbase管理员端
//
// @param host 	集群的主机列表[192.168.0.119:2181,192.169.0.120:2181]
//
func (this *HBase) OpenAdminClient(host string) {
	this.AdminClient = gohbase.NewAdminClient(host)

	return
}
