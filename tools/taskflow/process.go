package taskflow

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"reflect"
	"sync"
)

type TaskflowFactory struct {
	chainMap       map[string]*TaskflowExecutor //{name}-{version}
	latestChainMap map[string]*TaskflowExecutor //{name}
	chains         []*TaskflowExecutor

	taskMap      map[string]Task
	exceptionMap map[string]ExceptionTask

	initLock sync.Once

	persistenceService TaskChainService
}

func (t *TaskflowFactory) init(ctx context.Context) {
	t.initLock.Do(func() {
		t.chainMap = make(map[string]*TaskflowExecutor)
		t.latestChainMap = make(map[string]*TaskflowExecutor)
		t.taskMap = make(map[string]Task)
		t.exceptionMap = make(map[string]ExceptionTask)

		stop := &WaitForSignalTask{}
		t.taskMap[stop.Name()] = stop
	})
}

func (t *TaskflowFactory) ParseYaml(ctx context.Context, yamlStr string) error {
	t.init(ctx)

	chain := &TaskflowDef{}
	err := yaml.Unmarshal([]byte(yamlStr), chain)
	if err != nil {
		return err
	}
	err = chain.Validate(ctx)
	if err != nil {
		return err
	}

	var id string = fmt.Sprintf("%s-%d", chain.Name, chain.Version)
	chainExecutor := &TaskflowExecutor{
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

func (t *TaskflowFactory) RegisterPersistenceService(ctx context.Context, service TaskChainService) {
	t.persistenceService = service
}

func (t *TaskflowFactory) RegisterTask(ctx context.Context, task Task) {
	t.init(ctx)
	t.taskMap[task.Name()] = task
}

func (t *TaskflowFactory) RegisterException(ctx context.Context, exception ExceptionTask) {
	t.init(ctx)
	t.exceptionMap[exception.Name()] = exception
}

func (t *TaskflowFactory) Begin(ctx context.Context, chainName string, serviceId string, param ...interface{}) (interface{}, error) {
	t.init(ctx)
	if _, ok := t.latestChainMap[chainName]; !ok {
		return nil, fmt.Errorf("chain name[%s] not found", chainName)
	}
	chain := t.latestChainMap[chainName]
	return t.beginWithChan(ctx, chain, serviceId, param)
}

func (t *TaskflowFactory) BeginWithVersion(ctx context.Context, chainName string, chainVersion int, serviceId string, param ...interface{}) (interface{}, error) {
	t.init(ctx)
	chianId := fmt.Sprintf("%s-%d", chainName, chainVersion)
	if _, ok := t.chainMap[chianId]; !ok {
		return nil, fmt.Errorf("chain id[%s] not found", chianId)
	}
	chain := t.chainMap[chianId]
	return t.beginWithChan(ctx, chain, serviceId, param)
}

func (t *TaskflowFactory) beginWithChan(ctx context.Context, chain *TaskflowExecutor, serviceId string, param ...interface{}) (interface{}, error) {
	return chain.Begin(ctx, serviceId, "", param)
}

func (t *TaskflowFactory) StartTask(ctx context.Context, chainName string, serviceId string, param ...interface{}) (interface{}, error) {
	t.init(ctx)
	return t.StartTaskWithStageName(ctx, chainName, serviceId, "", param)
}

func (t *TaskflowFactory) StartTaskWithStageName(ctx context.Context,
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

func (t *TaskflowFactory) ReStartTaskWithStageName(ctx context.Context,
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

type TaskflowDef struct {
	Name    string      `yaml:"name,omitempty"`
	Version int         `yaml:"version,omitempty"`
	Stage   []*StageDef `yaml:"stage,omitempty"`
}
type StageDef struct {
	Name    string            `yaml:"name,omitempty"`
	Failure []string          `yaml:"failure,omitempty"`
	Args    map[string]string `yaml:"args,omitempty"`
}

func (t TaskflowDef) Validate(ctx context.Context) error {
	return nil
}

type Flow struct {
	t FlowType
}

func (f *Flow) Clear(ctx context.Context) {
	f.t = Normal
}
func (f *Flow) SetType(ctx context.Context, t FlowType) {
	f.t = t
}

type TaskExecutor struct {
	task      Task
	taskDef   *StageDef
	id        string
	next      *TaskExecutor
	pre       *TaskExecutor
	argument  map[string]string
	exception *ExceptionExecutor
}

func (t *TaskExecutor) Invoke(ctx context.Context, flow *Flow, result interface{}, param ...interface{}) (interface{}, error) {
	return t.task.Execute(ctx, flow, result, t.argument, param)
}

type ExceptionExecutor struct {
	failure ExceptionTask
	id      string
	next    *ExceptionExecutor
}

func (t *ExceptionExecutor) Invoke(ctx context.Context, stageName string, err error, param ...interface{}) error {
	return t.failure.Callback(ctx, stageName, err, param)
}

type TaskflowExecutor struct {
	factory   *TaskflowFactory
	def       *TaskflowDef
	taskMap   map[string]*TaskExecutor
	taskIdMap map[string]*TaskExecutor
	firstTask *TaskExecutor
	initLock  sync.Mutex
}

func (t *TaskflowExecutor) init(ctx context.Context) error {
	if t.firstTask != nil {
		return nil
	}
	t.initLock.Lock()
	defer t.initLock.Unlock()
	if t.firstTask != nil {
		return nil
	}
	for _, stage := range t.def.Stage {
		if _, ok := t.factory.taskMap[stage.Name]; !ok {
			return fmt.Errorf("task[%s] not definition", stage.Name)
		}

		for _, name := range stage.Failure {
			if _, ok := t.factory.exceptionMap[name]; !ok {
				return fmt.Errorf("failure[%s] not definition", name)
			}
		}
	}

	t.taskMap = make(map[string]*TaskExecutor)
	t.taskIdMap = make(map[string]*TaskExecutor)
	var preTask *TaskExecutor
	for index, stage := range t.def.Stage {
		id := fmt.Sprintf("%s:%d", stage.Name, index)
		instance := t.factory.taskMap[stage.Name]
		current := &TaskExecutor{
			task:     instance,
			id:       id,
			argument: stage.Args,
			taskDef:  stage,
		}
		if index == 0 {
			t.firstTask = current
		}
		if preTask != nil {
			preTask.next = current
			current.pre = preTask
		}
		t.taskMap[stage.Name] = current
		t.taskIdMap[id] = current

		if len(stage.Failure) > 0 {
			var preFailure *ExceptionExecutor
			for index1, failure1 := range stage.Failure {
				id1 := fmt.Sprintf("%s:%s:%d", stage.Name, failure1, index1)
				instance1 := t.factory.exceptionMap[failure1]
				current1 := &ExceptionExecutor{
					failure: instance1,
					id:      id1,
				}
				if index1 == 0 {
					current.exception = current1
				}
				if preFailure != nil {
					preFailure.next = current1
				}
				preFailure = current1
			}
		}
		preTask = current
	}
	return nil
}

func (t *TaskflowExecutor) Begin(ctx context.Context,
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

func (t *TaskflowExecutor) StartTaskFromStageId(ctx context.Context,
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

func (t *TaskflowExecutor) StartTaskFromStageName(ctx context.Context,
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

func (t *TaskflowExecutor) processTask(ctx context.Context, task *TaskExecutor,
	serviceId string, param ...interface{}) (interface{}, error) {

	currentTask := task
	var result interface{}
	var err error
	var flow *Flow = &Flow{t: Normal}
	for currentTask != nil {

		t.factory.persistenceService.SaveTaskStage(ctx, serviceId, currentTask.id, currentTask.taskDef, t.def)

		fmt.Printf("[%s:%d] %s begin invoke task %s \n", t.def.Name, t.def.Version, serviceId, currentTask.taskDef.Name)
		flow.Clear(ctx)
		currentResult, currentError := currentTask.Invoke(ctx, flow, result, param)
		fmt.Printf("[%s:%d] %s end invoke task %s \n", t.def.Name, t.def.Version, serviceId, currentTask.taskDef.Name)

		if currentResult != nil && !reflect.ValueOf(currentResult).IsZero() {
			result = currentResult
		}

		if flow.t == StopNext {
			if currentTask.next != nil {
				t.factory.persistenceService.SaveTaskStage(ctx, serviceId,
					currentTask.next.id,
					currentTask.next.taskDef,
					t.def)
			}
			if currentError == nil {
				break
			}
		} else if flow.t == StopCurrent {
			if currentError == nil {
				break
			}
		} else if flow.t == StopPre {
			if currentTask.pre != nil {
				t.factory.persistenceService.SaveTaskStage(ctx, serviceId,
					currentTask.pre.id,
					currentTask.pre.taskDef, t.def)
			}
			if currentError == nil {
				break
			}
		}

		if currentError != nil {
			err = currentError
			t.processFailure(ctx, serviceId, currentTask, err, param)
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

func (t *TaskflowExecutor) processFailure(ctx context.Context,
	serviceId string, task *TaskExecutor,
	err error, param ...interface{}) {
	currentFailure := task.exception
	if currentFailure == nil {
		return
	}
	for currentFailure != nil {
		err1 := currentFailure.Invoke(ctx, task.task.Name(), err, param)
		if err1 != nil {
			break
		}
		currentFailure = currentFailure.next
	}
}
