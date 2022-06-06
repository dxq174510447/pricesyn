package taskchain

import "context"

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
	SaveInstance(ctx context.Context, serviceId string, chainName string, chainVersion int) error
	SaveTaskStage(ctx context.Context, serviceId string, chainName string, chainVersion int, stageId string) error
	GetTaskId(ctx context.Context, serviceId string, chainId string) error
}
