package taskflow

const (
	name111 = "123"
)

type FlowType string

const (
	Normal      FlowType = "normal"
	StopCurrent FlowType = "stop_current"
	StopPre     FlowType = "stop_pre"
	StopNext    FlowType = "stop_next"
)
