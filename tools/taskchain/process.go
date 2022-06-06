package taskchain

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"sync"
)

type TaskChainFactory struct {
	chainMap       map[string]*TaskChainExecutor //{name}-{version}
	latestChainMap map[string]*TaskChainExecutor //{name}
	chains         []*TaskChainExecutor

	taskMap      map[string]Task
	exceptionMap map[string]ExceptionTask

	initLock sync.Once

	persistenceService TaskChainService
}

func (t *TaskChainFactory) init(ctx context.Context) {
	t.initLock.Do(func() {
		t.chainMap = make(map[string]*TaskChainExecutor)
		t.latestChainMap = make(map[string]*TaskChainExecutor)
		t.taskMap = make(map[string]Task)
		t.exceptionMap = make(map[string]ExceptionTask)
	})
}

func (t *TaskChainFactory) ParseYaml(ctx context.Context, yamlStr string) error {
	t.init(ctx)

	chain := &TaskChainDef{}
	err := yaml.Unmarshal([]byte(yamlStr), chain)
	if err != nil {
		return err
	}
	err = chain.Validate(ctx)
	if err != nil {
		return err
	}

	var id string = fmt.Sprintf("%s-%d", chain.Name, chain.Version)
	chainExecutor := &TaskChainExecutor{
		def:     chain,
		factory: t,
	}
	t.chainMap[id] = chainExecutor
	if lastChain, ok := t.latestChainMap[chain.Name]; ok {
		if chain.Version >= lastChain.def.Version {
			t.latestChainMap[chain.Name] = chainExecutor
		}
	} else {
		t.latestChainMap[chain.Name] = chainExecutor
	}
	t.chains = append(t.chains, chainExecutor)
	return nil
}

func (t *TaskChainFactory) RegisterPersistenceService(ctx context.Context, service TaskChainService) {
	t.persistenceService = service
}

func (t *TaskChainFactory) RegisterTask(ctx context.Context, task Task) {
	t.init(ctx)
	t.taskMap[task.Name()] = task
}

func (t *TaskChainFactory) RegisterException(ctx context.Context, exception ExceptionTask) {
	t.init(ctx)
	t.exceptionMap[exception.Name()] = exception
}

func (t *TaskChainFactory) StartByChainId(ctx context.Context, name string, serviceId string, param ...interface{}) (interface{}, error) {
	t.init(ctx)
	if _, ok := t.latestChainMap[name]; !ok {
		return nil, fmt.Errorf("chain name[%s] not found", name)
	}
	chain := t.latestChainMap[name]
	return chain.Start(ctx, serviceId, param)
}

type TaskChainDef struct {
	Name    string   `yaml:"name,omitempty"`
	Version int      `yaml:"version,omitempty"`
	Stage   []string `yaml:"stage,omitempty"`
	Failure []string `yaml:"failure,omitempty"`
}

func (t TaskChainDef) Validate(ctx context.Context) error {
	return nil
}

type TaskExecutor struct {
	task     Task
	id       string
	next     *TaskExecutor
	argument map[string]string
}

func (t *TaskExecutor) Invoke(ctx context.Context, result interface{}, param ...interface{}) (interface{}, error) {
	return t.task.Execute(ctx, result, t.argument, param)
}

type ExceptionExecutor struct {
	failure ExceptionTask
	id      string
	next    *ExceptionExecutor
}

func (t *ExceptionExecutor) Invoke(ctx context.Context, stageName string, err error, param ...interface{}) error {
	return t.failure.Callback(ctx, stageName, err, param)
}

type TaskChainExecutor struct {
	factory      *TaskChainFactory
	def          *TaskChainDef
	taskMap      map[string]*TaskExecutor
	failureMap   map[string]*ExceptionExecutor
	taskIdMap    map[string]*TaskExecutor
	failureIdMap map[string]*ExceptionExecutor
	firstTask    *TaskExecutor
	firstFailure *ExceptionExecutor
	initLock     sync.Mutex
}

func (t *TaskChainExecutor) init(ctx context.Context) error {
	if t.firstTask != nil {
		return nil
	}
	t.initLock.Lock()
	defer t.initLock.Unlock()
	if t.firstTask != nil {
		return nil
	}
	for _, name := range t.def.Stage {
		if _, ok := t.factory.taskMap[name]; !ok {
			return fmt.Errorf("task[%s] not definition", name)
		}
	}
	for _, name := range t.def.Failure {
		if _, ok := t.factory.exceptionMap[name]; !ok {
			return fmt.Errorf("failure[%s] not definition", name)
		}
	}

	t.taskMap = make(map[string]*TaskExecutor)
	t.taskIdMap = make(map[string]*TaskExecutor)
	t.failureMap = make(map[string]*ExceptionExecutor)
	t.failureIdMap = make(map[string]*ExceptionExecutor)
	var preTask *TaskExecutor
	var preFailure *ExceptionExecutor
	for index, name := range t.def.Stage {
		id := fmt.Sprintf("%s:%d", name, index)
		instance := t.factory.taskMap[name]
		current := &TaskExecutor{
			task: instance,
			id:   id,
		}
		if index == 0 {
			t.firstTask = current
		}
		if preTask != nil {
			preTask.next = current
		}
		t.taskMap[name] = current
		t.taskIdMap[id] = current
		preTask = current
	}

	for index, name := range t.def.Failure {
		id := fmt.Sprintf("%s:%d", name, index)
		instance := t.factory.exceptionMap[name]
		current := &ExceptionExecutor{
			failure: instance,
		}
		if index == 0 {
			t.firstFailure = current
		}
		if preFailure != nil {
			preFailure.next = current
		}
		t.failureMap[name] = current
		t.failureIdMap[id] = current
		preFailure = current
	}
	return nil
}

func (t *TaskChainExecutor) Start(ctx context.Context, serviceId string, param ...interface{}) (interface{}, error) {
	err := t.init(ctx)
	if err != nil {
		return nil, err
	}
	if t.firstTask == nil {
		return nil, fmt.Errorf("firstTask is nil")
	}
	current := t.firstTask
	if t.factory.persistenceService != nil {
		t.factory.persistenceService.SaveInstance(ctx, serviceId, t.def.Name, t.def.Version)
	}
	return t.processTask(ctx, current, serviceId, param)
}
func (t *TaskChainExecutor) processTask(ctx context.Context, task *TaskExecutor,
	serviceId string, param ...interface{}) (interface{}, error) {
	currentTask := task
	var result interface{}
	var err error
	for currentTask != nil {

		if t.factory.persistenceService != nil {
			t.factory.persistenceService.SaveTaskStage(ctx, serviceId, t.def.Name, t.def.Version, currentTask.id)
		}

		currentResult, currentError := currentTask.Invoke(ctx, result, param)
		if currentError != nil {
			err = currentError
			t.processFailure(ctx, serviceId, currentTask.task.Name(), err, param)
			break
		}
		if currentResult != nil && !reflect.ValueOf(currentResult).IsZero() {
			result = currentResult
		}
		currentTask = currentTask.next
	}
	if err != nil {
		return nil, err
	}
	return result, err
}

func (t *TaskChainExecutor) processFailure(ctx context.Context,
	serviceId string, stageName string,
	err error, param ...interface{}) {
	currentFailure := t.firstFailure
	if currentFailure == nil {
		return
	}
	for currentFailure != nil {
		err1 := currentFailure.Invoke(ctx, stageName, err, param)
		if err1 != nil {
			break
		}
		currentFailure = currentFailure.next
	}
}
