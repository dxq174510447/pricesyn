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

		stop := &WaitForSignalTask{}
		t.taskMap[stop.Name()] = stop
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

func (t *TaskChainFactory) Begin(ctx context.Context, chainName string, serviceId string, param ...interface{}) (interface{}, error) {
	t.init(ctx)
	if _, ok := t.latestChainMap[chainName]; !ok {
		return nil, fmt.Errorf("chain name[%s] not found", chainName)
	}
	chain := t.latestChainMap[chainName]
	return t.beginWithChan(ctx, chain, serviceId, param)
}

func (t *TaskChainFactory) BeginWithVersion(ctx context.Context, chainName string, chainVersion int, serviceId string, param ...interface{}) (interface{}, error) {
	t.init(ctx)
	chianId := fmt.Sprintf("%s-%d", chainName, chainVersion)
	if _, ok := t.chainMap[chianId]; !ok {
		return nil, fmt.Errorf("chain id[%s] not found", chianId)
	}
	chain := t.chainMap[chianId]
	return t.beginWithChan(ctx, chain, serviceId, param)
}

func (t *TaskChainFactory) beginWithChan(ctx context.Context, chain *TaskChainExecutor, serviceId string, param ...interface{}) (interface{}, error) {
	return chain.Begin(ctx, serviceId, "", param)
}

func (t *TaskChainFactory) StartTask(ctx context.Context, chainName string, serviceId string, param ...interface{}) (interface{}, error) {
	t.init(ctx)
	return t.StartTaskWithStageName(ctx, chainName, serviceId, "", param)
}

func (t *TaskChainFactory) StartTaskWithStageName(ctx context.Context,
	chainName string, serviceId string, stageName string, param ...interface{}) (interface{}, error) {
	t.init(ctx)

	if t.persistenceService == nil {
		return nil, fmt.Errorf("need persistenceService using Start* method")
	}

	stageId, chainVersion, end, err := t.persistenceService.GetStageId(ctx, serviceId, chainName)
	if err != nil {
		return nil, err
	}
	if end == 1 {
		fmt.Printf("[%s:%d] %s is finish\n", chainName, chainVersion, serviceId)
		return nil, nil
	}
	if len(stageId) == 0 {
		return nil, fmt.Errorf("no chain relation to %s serviceId %s", chainName, serviceId)
	}

	chainId := fmt.Sprintf("%s-%d", chainName, chainVersion)
	if _, ok := t.chainMap[chainId]; !ok {
		return nil, fmt.Errorf("chain id[%s] not found", chainId)
	}
	chain := t.chainMap[chainId]
	fmt.Printf("[%s:%d] %s start from %s\n", chain.def.Name, chain.def.Version, serviceId, stageId)
	return chain.StartTaskFromStageId(ctx, serviceId, stageId, stageName, param)
}

func (t *TaskChainFactory) ReStartTaskWithStageName(ctx context.Context,
	chainName string, serviceId string, stageName string, param ...interface{}) (interface{}, error) {
	t.init(ctx)

	if t.persistenceService == nil {
		return nil, fmt.Errorf("need persistenceService using ReStart* method")
	}

	stageId, chainVersion, _, err := t.persistenceService.GetStageId(ctx, serviceId, chainName)
	if err != nil {
		return nil, err
	}
	if len(stageId) == 0 {
		return nil, fmt.Errorf("no chain relation to %s serviceId %s", chainName, serviceId)
	}

	fmt.Printf("[%s:%d] %s current stage %s restart from stage %s\n", chainName, chainVersion,
		serviceId, stageId, stageName)

	chainId := fmt.Sprintf("%s-%d", chainName, chainVersion)
	if _, ok := t.chainMap[chainId]; !ok {
		return nil, fmt.Errorf("chain id[%s] not found", chainId)
	}
	chain := t.chainMap[chainId]
	fmt.Printf("[%s:%d] %s start from %s\n", chain.def.Name, chain.def.Version, serviceId, stageId)
	return chain.StartTaskFromStageName(ctx, serviceId, stageName, param)
}

type TaskChainDef struct {
	Name    string      `yaml:"name,omitempty"`
	Version int         `yaml:"version,omitempty"`
	Stage   []*StageDef `yaml:"stage,omitempty"`
	Failure []*StageDef `yaml:"failure,omitempty"`
}
type StageDef struct {
	Name string            `yaml:"name,omitempty"`
	Args map[string]string `yaml:"args,omitempty"`
}

func (t TaskChainDef) Validate(ctx context.Context) error {
	return nil
}

type TaskExecutor struct {
	task     Task
	taskDef  *StageDef
	id       string
	next     *TaskExecutor
	pre      *TaskExecutor
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
		if _, ok := t.factory.taskMap[name.Name]; !ok {
			return fmt.Errorf("task[%s] not definition", name.Name)
		}
	}
	for _, name := range t.def.Failure {
		if _, ok := t.factory.exceptionMap[name.Name]; !ok {
			return fmt.Errorf("failure[%s] not definition", name.Name)
		}
	}

	t.taskMap = make(map[string]*TaskExecutor)
	t.taskIdMap = make(map[string]*TaskExecutor)
	t.failureMap = make(map[string]*ExceptionExecutor)
	t.failureIdMap = make(map[string]*ExceptionExecutor)
	var preTask *TaskExecutor
	var preFailure *ExceptionExecutor
	for index, name := range t.def.Stage {
		id := fmt.Sprintf("%s:%d", name.Name, index)
		instance := t.factory.taskMap[name.Name]
		current := &TaskExecutor{
			task:     instance,
			id:       id,
			argument: name.Args,
			taskDef:  name,
		}
		if index == 0 {
			t.firstTask = current
		}
		if preTask != nil {
			preTask.next = current
			current.pre = preTask
		}
		t.taskMap[name.Name] = current
		t.taskIdMap[id] = current
		preTask = current
	}

	for index, name := range t.def.Failure {
		id := fmt.Sprintf("%s:%d", name.Name, index)
		instance := t.factory.exceptionMap[name.Name]
		current := &ExceptionExecutor{
			failure: instance,
			id:      id,
		}
		if index == 0 {
			t.firstFailure = current
		}
		if preFailure != nil {
			preFailure.next = current
		}
		t.failureMap[name.Name] = current
		t.failureIdMap[id] = current
		preFailure = current
	}
	return nil
}

func (t *TaskChainExecutor) Begin(ctx context.Context,
	serviceId string, stageName string, param ...interface{}) (interface{}, error) {
	err := t.init(ctx)
	if err != nil {
		return nil, err
	}
	if t.firstTask == nil {
		return nil, fmt.Errorf("firstTask is nil")
	}
	current := t.firstTask

	if len(stageName) > 0 && stageName != current.taskDef.Name {
		return nil, fmt.Errorf("serviceId %s current stage %s,not begin with %s", serviceId, current.taskDef.Name, stageName)
	}

	if t.factory.persistenceService != nil {

		stageId1, _, _, err := t.factory.persistenceService.GetStageId(ctx, serviceId, t.def.Name)
		if err != nil {
			return nil, err
		}
		if len(stageId1) > 0 {
			return nil, fmt.Errorf("chian[%s] alreay has serviceId %s", t.def.Name, serviceId)
		}

		pid, err := t.factory.persistenceService.SaveInstance(ctx, serviceId, t.def)
		if err != nil {
			return nil, err
		}
		fmt.Printf("[%s:%d] %s start chain %s \n", t.def.Name, t.def.Version, serviceId, pid)
	}
	return t.processTask(ctx, current, serviceId, param)
}

func (t *TaskChainExecutor) StartTaskFromStageId(ctx context.Context,
	serviceId string, stageId string, stageName string, param ...interface{}) (interface{}, error) {
	err := t.init(ctx)
	if err != nil {
		return nil, err
	}
	if _, ok := t.taskIdMap[stageId]; !ok {
		return nil, fmt.Errorf("%s:%d task id %s not exists", t.def.Name, t.def.Version, stageId)
	}
	current := t.taskIdMap[stageId]
	if len(stageName) > 0 && stageName != current.taskDef.Name {
		return nil, fmt.Errorf("serviceId %s current stage %s,not begin with %s", serviceId, current.taskDef.Name, stageName)
	}
	return t.processTask(ctx, current, serviceId, param)
}

func (t *TaskChainExecutor) StartTaskFromStageName(ctx context.Context,
	serviceId string, stageName string, param ...interface{}) (interface{}, error) {
	err := t.init(ctx)
	if err != nil {
		return nil, err
	}
	if _, ok := t.taskMap[stageName]; !ok {
		return nil, fmt.Errorf("%s:%d task name %s not exists", t.def.Name, t.def.Version, stageName)
	}
	current := t.taskMap[stageName]
	return t.processTask(ctx, current, serviceId, param)
}

func (t *TaskChainExecutor) processTask(ctx context.Context, task *TaskExecutor,
	serviceId string, param ...interface{}) (interface{}, error) {

	currentTask := task
	var result interface{}
	var err error

	for currentTask != nil {

		if t.factory.persistenceService != nil {
			t.factory.persistenceService.SaveTaskStage(ctx, serviceId, currentTask.id, currentTask.taskDef, t.def)
		}

		fmt.Printf("[%s:%d] %s begin invoke task %s \n", t.def.Name, t.def.Version, serviceId, currentTask.taskDef.Name)
		currentResult, currentError := currentTask.Invoke(ctx, result, param)
		fmt.Printf("[%s:%d] %s end invoke task %s \n", t.def.Name, t.def.Version, serviceId, currentTask.taskDef.Name)

		var stopError *WaitForSignalException
		if currentError != nil {
			switch v := currentError.(type) {
			case WaitForSignalException:
				stopError = &v
			case *WaitForSignalException:
				stopError = v
			}

			if stopError == nil {
				err = currentError
				t.processFailure(ctx, serviceId, currentTask.task.Name(), err, param)
				break
			}
		}
		if currentResult != nil && !reflect.ValueOf(currentResult).IsZero() {
			result = currentResult
		}

		if stopError != nil && t.factory.persistenceService != nil {
			if stopError.nextStage == 1 {
				if currentTask.next != nil {
					t.factory.persistenceService.SaveTaskStage(ctx, serviceId, currentTask.next.id, currentTask.next.taskDef, t.def)
				}
			} else if stopError.nextStage == 0 {

			} else if stopError.nextStage == -1 {
				if currentTask.pre != nil {
					t.factory.persistenceService.SaveTaskStage(ctx, serviceId, currentTask.pre.id, currentTask.pre.taskDef, t.def)
				}
			}
			break
		}
		currentTask = currentTask.next
	}
	if err != nil {
		return nil, err
	}
	if currentTask == nil && t.factory.persistenceService != nil {
		t.factory.persistenceService.EndInstance(ctx, serviceId, t.def)
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
