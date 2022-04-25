package taskchain

import "strings"

type NodeType string

func (n NodeType) String() string {
	return string(n)
}

func GetNodeType(name string) NodeType {
	n := strings.ToLower(name)
	switch n {
	case "begin":
		return BEGIN
	case "end":
		return END
	case "task":
		return TASK
	}
	return NodeType("unknown")
}

const (
	BEGIN NodeType = NodeType("begin")
	END   NodeType = NodeType("end")
	TASK  NodeType = NodeType("task")
)
