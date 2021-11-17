package graph

import (
	"fmt"
	"pricesyn/util"
	"testing"
)

func TestBreadthFirstAlg_Find(t *testing.T) {

	links := [][]string{
		[]string{"cab", "car"},
		[]string{"cab", "cat"},
		[]string{"car", "bar"},
		[]string{"car", "bat"},
		[]string{"cat", "mat"},
		[]string{"cat", "bar"},
		[]string{"mat", "bat"},
		[]string{"car", "bar"},
		[]string{"bar", "bat"},
	}
	g := NewSimpleGraph(links, false)
	graph := BreadthFirstAlg{
		Target: g,
	}
	r, f := graph.Find("cab", "bat")
	if !f || len(r) == 0 {
		fmt.Println("not find result")
	} else {
		fmt.Println(util.JsonUtil.To2String(r))
	}

}
