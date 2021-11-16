package graph


/**
广度优先算法
广度优先搜索不仅查找从A到B的路径，而且找到的是最短的路径
从起点开始向外延伸
*/

type NonDirectionGraph struct {
	Target *Graph
}

func (n *NonDirectionGraph) Find(source string,target string) []string {
	//var checkNode map[string]*GraphNode = make(map[string]*GraphNode)
	//if source == target {
	//	return []string{source}
	//}
//	checkNode[]
//	for _,node := range n.Target.Node {
//
//	}
	return []string{}
}

