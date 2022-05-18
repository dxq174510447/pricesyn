package taskchain

import "context"

type Task interface {
	Name() string
	Execute(ctx context.Context, result interface{}, argument map[string]string, param ...interface{}) (interface{}, error)
}

type Exception interface {
	Name() string
	Callback(ctx context.Context, err error, args ...interface{})
}
