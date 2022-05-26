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
