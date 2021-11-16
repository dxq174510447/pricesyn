package repository

import (
	"context"
	"fmt"
	"testing"
)

func TestDcScheduleRep_GetSchedule(t *testing.T) {
	r , err := DcScheduleRepImpl.GetSchedule(context.Background(),"GDR","2021-11-12",2)
	if err != nil {
		panic(err)
	}
	fmt.Println(r)
}