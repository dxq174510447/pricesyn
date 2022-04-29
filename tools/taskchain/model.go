package taskchain

import (
	"context"
	"fmt"
	"gopkg.in/yaml.v3"
	"sync"
)

type TaskChainFactory struct {
	chainMap    map[string]*TaskChainDef //{name}-{version}
	latestChain map[string]*TaskChainDef //{name}
	chains      []*TaskChainDef

	tasks      map[string]Task
	exceptions map[string]Exception

	initLock sync.Once
}

func (t *TaskChainFactory) init(ctx context.Context) {
	t.initLock.Do(func() {
		t.chainMap = make(map[string]*TaskChainDef)
		t.latestChain = make(map[string]*TaskChainDef)
		t.tasks = make(map[string]Task)
		t.exceptions = make(map[string]Exception)
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
	t.chainMap[id] = chain
	if lastChain, ok := t.latestChain[chain.Name]; ok {
		if chain.Version >= lastChain.Version {
			t.latestChain[chain.Name] = chain
		}
	} else {
		t.latestChain[chain.Name] = chain
	}
	t.chains = append(t.chains, chain)
	return nil
}

func (t *TaskChainFactory) RegisterTask(ctx context.Context, task Task) {
	t.init(ctx)
	t.tasks[task.Name()] = task
}

func (t *TaskChainFactory) RegisterException(ctx context.Context, serviceName string, exception Exception) {
	t.init(ctx)
	t.exceptions[exception.Name()] = exception
}

type TaskCommand struct {
	task Task
	next Task
}

func (t *TaskCommand) Exec(ctx context.Context) {

}

type TaskChainDef struct {
	Name    string   `yaml:"name,omitempty"`
	Version int      `yaml:"version,omitempty"`
	Stages  []string `yaml:"stages,omitempty"`
	Failure []string `yaml:"failure,omitempty"`
}

func (t TaskChainDef) Validate(ctx context.Context) error {
	return nil
}
