package util

import (
	"reflect"
)

type arrayUtil struct {
}

func (d *arrayUtil) Reverse(s interface{}) {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
}

var ArrayUtil arrayUtil = arrayUtil{}
