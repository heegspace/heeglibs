package youyoulibs

import (
	"fmt"
	"testing"
)

func TestHBasePut(t *testing.T) {
	hbase := NewHBase()
	hbase.OpenClient("192.168.0.119,192.168.0.120,192.168.0.121")
	var value HBaseResult
	value.RowKey = "rowkey2"
	value.Family = "schoolInfo"
	value.Qualifier = "name"
	value.Value = make([]byte, len("Univesty of Kunming"))
	copy(value.Value, []byte("Univesty of Guangzhou"))

	err := hbase.PutCell("Student", value)
	if nil != err {
		fmt.Println("GetRow error: ", err)

		return
	}

	return
}

func TestHBaseRow(t *testing.T) {
	hbase := NewHBase()
	hbase.OpenClient("192.168.0.119,192.168.0.120,192.168.0.121")
	row, err := hbase.GetRow("Student", "rowkey2")
	if nil != err {
		fmt.Println("GetRow error: ", err)

		return
	}

	for i, v := range row {
		str := fmt.Sprintf("index:%d		Rowkey: %s		family: %s		qualifies:%s		value:%s", i, v.RowKey, v.Family, v.Qualifier, string(v.Value))
		fmt.Println(str)
	}

	return
}

func TestHBaseCell(t *testing.T) {
	hbase := NewHBase()
	hbase.OpenClient("192.168.0.119,192.168.0.120,192.168.0.121")
	cell, err := hbase.GetCell("Student", "rowkey1", "baseInfo", []string{"age", "name"})
	if nil != err {
		fmt.Println("GetCell error: ", err)

		return
	}

	for k, v := range cell {
		str := fmt.Sprintf("index:%s		Rowkey: %s		family: %s		qualifies:%s		value:%s", k, v.RowKey, v.Family, v.Qualifier, string(v.Value))
		fmt.Println(str)
	}
}

func TestHBaseDel(t *testing.T) {
	hbase := NewHBase()
	hbase.OpenClient("192.168.0.119,192.168.0.120,192.168.0.121")

	var value HBaseResult
	value.RowKey = "rowkey1"
	value.Family = "schoolInfo"
	value.Qualifier = "name"
	value.Value = nil

	// 删除列
	err := hbase.DelCell("Student", value)
	if nil != err {
		fmt.Println("DelCell cell error: ", err)

		return
	}

	value.Qualifier = ""
	err = hbase.DelCell("Student", value)
	if nil != err {
		fmt.Println("DelCell family error: ", err)

		return
	}

	return
}
