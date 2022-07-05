package taskflow

import (
	"context"
)

type Task interface {
	Name() string
	/*
		Execute

		result 上一次非空返回结果

		argument 配置文件中的参数

		param 方法参数
	*/
	Execute(ctx context.Context, flow *Flow, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error)
}

type ExceptionTask interface {
	Name() string
	Callback(ctx context.Context, stageName string, err error, param ...interface{}) error
}

type TaskChainService interface {
	SaveInstance(ctx context.Context, serviceId string, def *TaskflowDef) (string, error)
	SaveTaskStage(ctx context.Context, serviceId string, stageId string, stageDef *StageDef, def *TaskflowDef) error
	EndInstance(ctx context.Context, serviceId string, def *TaskflowDef) error
	GetStageId(ctx context.Context, serviceId string, chainName string) (string, int, int, error)
}

type WaitForSignalTask struct {
}

func (w WaitForSignalTask) Name() string {
	return "_stop"
}

func (w WaitForSignalTask) Execute(ctx context.Context, flow *Flow, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	flow.SetType(ctx, StopNext)
	return nil, nil
}
