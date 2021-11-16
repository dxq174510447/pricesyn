package graph

type GraphNode struct {
	Name string
	Target []*GraphLink
	Source []*GraphLink
	targetMap map[string]*GraphLink
	sourceMap map[string]*GraphLink
}


type GraphLink struct {
	From *GraphNode
	To *GraphNode
}

type Graph struct {
	Node []*GraphNode
	NodeMap map[string]*GraphNode
}

func (g *Graph) AddGraphNode(name string){
	if _,ok := g.NodeMap[name];ok {
		return
	}
	n := NewSimpleGraphNode(name)
	g.Node = append(g.Node,n)
	g.NodeMap[name] = n
}

func NewGraphLink(form,to *GraphNode) *GraphLink{
	return &GraphLink{
		From: form,
		To: to,
	}
}

func NewSimpleGraphNode(name string) *GraphNode {
	return &GraphNode{
		Name: name,
		targetMap: make(map[string]*GraphLink),
		sourceMap: make(map[string]*GraphLink),
	}
}

func NewSimpleGraph(ts [][]string,direction bool) *Graph {
	var nodeRef map[string]*GraphNode = make(map[string]*GraphNode)
	var result *Graph = &Graph{}
	for _,t := range ts {
		n0 := t[0]
		n1 := t[1]

		if _,ok := nodeRef[n0];!ok {
			nn0 := NewSimpleGraphNode(n0)
			nodeRef[n0] = nn0
			result.Node = append(result.Node,nn0)
		}

		if _,ok := nodeRef[n1];!ok {
			nn1 := NewSimpleGraphNode(n1)
			nodeRef[n1] = nn1
			result.Node = append(result.Node,nn1)
		}
	}
	for _,t := range ts {
		sourceNode := nodeRef[t[0]]
		targetNode := nodeRef[t[1]]

		//link := NewGraphLink(sourceNode,targetNode)
		if direction {
			if _,ok := sourceNode.targetMap[targetNode.Name]; !ok {
				link := NewGraphLink(sourceNode,targetNode)
				sourceNode.Target = append(sourceNode.Target,link)
				sourceNode.targetMap[targetNode.Name] = link

				targetNode.Source = append(targetNode.Source,link)
				targetNode.sourceMap[sourceNode.Name] = link
			}
		}else{
			var hasAdd bool = false

			if _,ok := sourceNode.targetMap[targetNode.Name]; ok {
				hasAdd = true
			}

			if !hasAdd {
				if _,ok := sourceNode.sourceMap[targetNode.Name]; ok {
					hasAdd = true
				}
			}

			if !hasAdd {
				link := NewGraphLink(sourceNode,targetNode)
				sourceNode.Target = append(sourceNode.Target,link)
				sourceNode.targetMap[targetNode.Name] = link

				targetNode.Source = append(targetNode.Source,link)
				targetNode.sourceMap[sourceNode.Name] = link

			}
		}
	}
	return result
}
