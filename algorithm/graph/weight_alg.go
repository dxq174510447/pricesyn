package graph

import (
	"pricesyn/util"
	"sync"
)

/**
有权图：狄克斯特拉算法

1. 每次从权重最小节点计算周边节点的权重，并更新。
2. 不支持负权重, 不支持有环图
*/

type WeightAlgResult struct {
	Weight int
	Nodes  []string
}

type WeightAlg struct {
	Target   *Graph
	initLock sync.Once
}

func (w *WeightAlg) init() {

}

func (n *WeightAlg) Find(source string, target string) (*WeightAlgResult, bool) {

	var costs map[string]int = make(map[string]int)         // 节点-->权重
	var parents map[string]string = make(map[string]string) // 节点-->父节点
	var checkedNode map[string]int = make(map[string]int)   // 节点检查次数 防止环

	if source == target {
		return &WeightAlgResult{
			Weight: 0,
			Nodes:  []string{source},
		}, true
	}

	sourceNode := n.Target.NodeMap[source]
	targetNode := n.Target.NodeMap[target]
	if sourceNode == nil || targetNode == nil {
		return nil, false
	}

	checkedNode[sourceNode.Name] = 1
	for _, out := range sourceNode.Target {
		parents[out.To.Name] = sourceNode.Name
		costs[out.To.Name] = out.Weight
	}
	found := n.doFind(target, costs, parents, checkedNode)

	if !found {
		return nil, false
	}

	var result []string
	result = append(result, target)
	var current string = target
	for {
		if parent, ok := parents[current]; ok {
			result = append(result, parent)
			current = parent
		} else {
			break
		}
	}

	util.ArrayUtil.Reverse(result)
	return &WeightAlgResult{
		Weight: costs[target],
		Nodes:  result,
	}, true

}

func (w *WeightAlg) doFind(target string, costs map[string]int, parents map[string]string, checkedNode map[string]int) bool {
	var minNodeName string
	var minWeight int
	for nodeName, weight := range costs {

		if _, ok := checkedNode[nodeName]; ok {
			continue
		}

		if minWeight == 0 {
			minNodeName = nodeName
			minWeight = weight
		} else if minWeight > weight {
			minNodeName = nodeName
			minWeight = weight
		}
	}

	if minNodeName == target {
		return true
	}

	if minNodeName == "" {
		// 都处理过了 但是没有
		return false
	}

	if flag, ok := checkedNode[minNodeName]; ok {
		checkedNode[minNodeName] = flag + 1
	} else {
		checkedNode[minNodeName] = 1
	}

	minNode := w.Target.NodeMap[minNodeName]
	for _, out := range minNode.Target {
		currentWeight := minWeight + out.Weight
		if weight, ok := costs[out.To.Name]; ok {
			if currentWeight < weight {
				costs[out.To.Name] = currentWeight
				parents[out.To.Name] = minNode.Name
			}
		} else {
			costs[out.To.Name] = currentWeight
			parents[out.To.Name] = minNode.Name
		}
	}
	return w.doFind(target, costs, parents, checkedNode)
}
