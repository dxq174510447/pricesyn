package file

import (
	"fmt"
	"testing"
)

func TestFileRowRead_ScanCsvRow(t *testing.T) {
	file := "/Users/klook/code/company/pricesyn/data/config_pkg.csv"
	r := &FileRowRead{}
	r.Parse(file)
	r.ScanCsvRow(func(rowIndex int, columns []string) {
		fmt.Println(rowIndex,len(columns),columns[0],columns[1],columns[2])
	})
}
