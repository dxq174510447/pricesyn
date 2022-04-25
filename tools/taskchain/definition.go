package taskchain

import "context"

type Process struct {
	startNode       Noder
	nodes           []Noder
	nodeMap         map[string]Noder //id对应的noder
	taskNodes       []TaskNoder
	taskNodeMap     map[string]TaskNoder //id对应的noder
	taskNodeTypeMap map[string][]TaskNoder
	directionMap    map[string]Direction
}

type Direction struct {
	Name   string
	Target Noder
}

type Noder interface {
	Id(ctx context.Context) string
	Execute(ctx context.Context, param interface{}) error
	OutDirections(ctx context.Context) ([]*Direction, error)
}
type TaskNoder interface {
	Noder
	TaskType(ctx context.Context) (string, error)
}

type Node struct {
	idKey string
	outs  []*Direction
}

type BeginNode struct {
	Node
}

func (b BeginNode) Id(ctx context.Context) string {
	return b.idKey
}

func (b BeginNode) Execute(ctx context.Context, param interface{}) error {
	return nil
}

func (b BeginNode) OutDirections(ctx context.Context) ([]*Direction, error) {
	return b.outs, nil
}

type EndNode struct {
	Node
}

func (b EndNode) Id(ctx context.Context) string {
	return b.idKey
}

func (b EndNode) Execute(ctx context.Context, param interface{}) error {
	return nil
}

func (b EndNode) OutDirections(ctx context.Context) ([]*Direction, error) {
	return b.outs, nil
}

type ExclusiveGatewayNode struct {
	Node
}

func (b ExclusiveGatewayNode) Id(ctx context.Context) string {
	return b.idKey
}

func (b ExclusiveGatewayNode) Execute(ctx context.Context, param interface{}) error {
	return nil
}

func (b ExclusiveGatewayNode) OutDirections(ctx context.Context) ([]*Direction, error) {
	return b.outs, nil
}

type TaskNode struct {
	Node
	taskType string
}

func (b TaskNode) Id(ctx context.Context) string {
	return b.idKey
}

func (b TaskNode) Execute(ctx context.Context, param interface{}) error {
	return nil
}

func (b TaskNode) OutDirections(ctx context.Context) ([]*Direction, error) {
	return b.outs, nil
}

func (b TaskNode) TaskType(ctx context.Context) (string, error) {
	return b.taskType, nil
}
