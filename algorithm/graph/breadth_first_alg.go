package graph

import (
	"pricesyn/algorithm/mq"
	"sync"
)

/**
广度优先算法
1. 从节点A出发 是否有到b的路径
2. 从节点A出发，到b最短的路径

从起点开始向外延伸，先检查一度关系，在检查二度关系
*/

type TaskNode struct {
	CurrentNode *GraphNode
	Target      *GraphNode
	Links       []string
}

func (t *TaskNode) IdKey() string {
	return t.CurrentNode.Name
}

type BreadthFirstAlg struct {
	Target   *Graph
	taskPool *mq.LinkedBlockingQueue
	checked  map[string]*GraphNode
	initLock sync.Once
}

func (n *BreadthFirstAlg) init() {
	n.initLock.Do(func() {
		if n.taskPool == nil {
			n.taskPool = mq.NewLinkedBlockingQueue()
		}
		if n.checked == nil {
			n.checked = make(map[string]*GraphNode)
		}
	})
}

func (n *BreadthFirstAlg) Find(source string, target string) ([]string, bool) {
	n.init()

	if source == target {
		return []string{source}, true
	}

	sourceNode := n.Target.NodeMap[source]
	targetNode := n.Target.NodeMap[target]

	if sourceNode == nil || targetNode == nil {
		return nil, false
	}

	tn := &TaskNode{}
	tn.CurrentNode = sourceNode
	tn.Target = targetNode
	tn.Links = []string{sourceNode.Name}

	n.taskPool.Offer(tn)

	result := n.doFind()
	r1 := <-result

	if len(r1) == 0 {
		return nil, false
	}
	return r1, true
}

func (n *BreadthFirstAlg) doFind() chan []string {
	result := make(chan []string)
	go func() {
		for {
			n.taskPool.PrintLink()
			element := n.taskPool.Poll(0)
			if element == nil {
				result <- []string{}
				break
			}
			node := element.(*TaskNode)

			if _, ok := n.checked[node.CurrentNode.Name]; ok {
				continue
			}

			if node.CurrentNode.Name == node.Target.Name {
				result <- node.Links
				break
			}

			for _, next := range node.CurrentNode.Target {
				if _, ok := n.checked[next.To.Name]; ok {
					continue
				}

				tn := &TaskNode{}
				tn.CurrentNode = next.To
				tn.Target = node.Target

				tn.Links = make([]string, len(node.Links)+1)
				copy(tn.Links, node.Links)
				tn.Links[len(tn.Links)-1] = next.To.Name

				n.taskPool.Offer(tn)
			}

			for _, next := range node.CurrentNode.Source {
				if _, ok := n.checked[next.From.Name]; ok {
					continue
				}

				tn := &TaskNode{}
				tn.CurrentNode = next.From
				tn.Target = node.Target

				tn.Links = make([]string, len(node.Links)+1)
				copy(tn.Links, node.Links)
				tn.Links[len(tn.Links)-1] = next.From.Name

				n.taskPool.Offer(tn)
			}
			n.checked[node.CurrentNode.Name] = node.CurrentNode
		}
	}()
	return result
}
