package util

import (
	"fmt"
	"github.com/looplab/fsm"
)

type machineUtil struct {
}

func (c *machineUtil) NewMachine() error {

	fsm := fsm.NewFSM("closed", fsm.Events{
		{Name: "open1", Src: []string{"closed"}, Dst: "open"},
		{Name: "close1", Src: []string{"open"}, Dst: "closed"},
	},
		fsm.Callbacks{
			"before_open1": func(event *fsm.Event) {
				fmt.Println("enter_open")
				err := fmt.Errorf("err %s", "enter_open")
				event.Cancel(err)
			},
			"after_event": func(event *fsm.Event) {
				fmt.Println("after_event")
			},
		},
	)
	fmt.Println(fsm.Current())

	err := fsm.Event("open1", "aa", "bbb")
	if err != nil {
		fmt.Println(err)
	}

	return nil
}

var MachineUtil machineUtil = machineUtil{}
