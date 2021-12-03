package graph

import (
	"fmt"
	"pricesyn/util"
	"testing"
)

func TestWeightAlgResult_Find(t *testing.T) {

	links := [][]string{
		[]string{"a", "b", "5"},
		[]string{"a", "c", "2"},
		[]string{"b", "d", "4"},
		[]string{"b", "e", "2"},
		[]string{"c", "b", "8"},
		//[]string{"c", "e","7"},
		[]string{"c", "e", "3"},
		[]string{"d", "f", "3"},
		[]string{"d", "e", "6"},
		[]string{"e", "f", "1"},
	}
	g := NewWeightGraph(links)
	graph := WeightAlg{
		Target: g,
	}
	r, f := graph.Find("a", "f")
	if !f || r == nil {
		fmt.Println("not find result")
	} else {
		fmt.Println(util.JsonUtil.To2String(r))
	}

}
