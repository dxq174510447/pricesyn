package taskchain

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
	Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error)
}

type ExceptionTask interface {
	Name() string
	Callback(ctx context.Context, stageName string, err error, param ...interface{}) error
}

type TaskChainService interface {
	SaveInstance(ctx context.Context, serviceId string, def *TaskChainDef) (string,error)
	SaveTaskStage(ctx context.Context, serviceId string,stageId string,stageDef *StageDef,def *TaskChainDef) error
	EndInstance(ctx context.Context, serviceId string, def *TaskChainDef) error
	GetStageId(ctx context.Context, serviceId string, chainName string) (string,int,error)
}


type waitForSignalException struct {
	//唤醒之后从（-1 上一个任务 0当前任务 1下一个任务） 开始执行
	nextStage int
}

func (w waitForSignalException) Error() string {
	return ""
}

type WaitForSignalTask struct {

}

func (w WaitForSignalTask) Name() string {
	return "stop"
}

func (w WaitForSignalTask) Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error) {
	return nil,
}
