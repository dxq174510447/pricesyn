package db

import (
	"context"
	"testing"
	"time"
)

func TestDbFactory_Get(t *testing.T) {

	factory := DbFactory{}
	ctx := context.Background()
	go func() {
		factory.GetDb(ctx)
	}()
	go func() {
		factory.GetDb(ctx)
	}()

	time.Sleep(time.Hour * 1)
}
