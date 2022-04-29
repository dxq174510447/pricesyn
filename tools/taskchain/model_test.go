package taskchain

import (
	"context"
	"fmt"
	"testing"
)

func TestTaskChainFactory_ParseYaml(t *testing.T) {

	ctx := context.Background()
	factory := &TaskChainFactory{}
	err := factory.ParseYaml(ctx, Eg1)
	if err != nil {
		t.Fatalf("%v", err)
	}
	fmt.Printf("%d", len(factory.chains))
}
