package graph

import "sync"

type GraphNode struct {
	Name      string
	Target    []*GraphLink
	Source    []*GraphLink
	targetMap map[string]*GraphLink
	sourceMap map[string]*GraphLink
}

type GraphLink struct {
	From *GraphNode
	To   *GraphNode
}

type Graph struct {
	Node     []*GraphNode
	NodeMap  map[string]*GraphNode
	initLock sync.Once
}

func (g *Graph) init() {
	g.initLock.Do(func() {
		if g.NodeMap == nil {
			g.NodeMap = make(map[string]*GraphNode)
		}
	})
}

func (g *Graph) AddGraphNode(name string) *GraphNode {
	g.init()
	if node, ok := g.NodeMap[name]; ok {
		return node
	}
	n := &GraphNode{
		Name:      name,
		targetMap: make(map[string]*GraphLink),
		sourceMap: make(map[string]*GraphLink),
	}
	g.Node = append(g.Node, n)
	g.NodeMap[name] = n
	return n
}
func (g *Graph) AddGraphLink(from, to *GraphNode, direction bool) *GraphLink {
	g.init()
	if direction {
		// 有向图
		if link, ok := from.targetMap[to.Name]; ok {
			return link
		} else {
			link = &GraphLink{
				From: from,
				To:   to,
			}
			from.Target = append(from.Target, link)
			from.targetMap[to.Name] = link

			to.Source = append(to.Source, link)
			to.sourceMap[from.Name] = link
			return link
		}
	} else {
		// 无向图

		if link, ok := from.targetMap[to.Name]; ok {
			return link
		}

		if link, ok := from.sourceMap[to.Name]; ok {
			return link
		}

		link := &GraphLink{
			From: from,
			To:   to,
		}

		from.Target = append(from.Target, link)
		from.targetMap[to.Name] = link

		to.Source = append(to.Source, link)
		to.sourceMap[from.Name] = link
		return link

	}
}

func NewSimpleGraph(ts [][]string, direction bool) *Graph {
	g := &Graph{}
	for _, link := range ts {
		fromNode := g.AddGraphNode(link[0])
		toNode := g.AddGraphNode(link[1])
		g.AddGraphLink(fromNode, toNode, direction)
	}
	return g
}
