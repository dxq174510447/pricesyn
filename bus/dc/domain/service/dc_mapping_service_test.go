package service

import (
	"context"
	"fmt"
	"testing"
)

func TestDcMappingService_FindTop(t *testing.T) {

	r,err := DcMappingServiceImp.FindTop(context.Background(),3)
	if err != nil {
		panic(err)
	}

	fmt.Println(len(r))

}