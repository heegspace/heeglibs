package heeglibs

import (
	"fmt"
	"runtime"
	"testing"
)

func TestLoadTableFromFile(t *testing.T) {
	LoadTableFromFile("./table.json")
}

func TestTableName(t *testing.T) {
	funcName, file, line, ok := runtime.Caller(0)
	fmt.Println(funcName, file, line, ok)
	fmt.Println(TableName("USER"))
}
