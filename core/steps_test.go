package core

import (
	"github.com/lucasmbaia/kraken/workflow"
	"testing"
	"fmt"
)

func Test_Steps(t *testing.T) {
	var tests = []struct{
		name	string
		wf	workflow.Workflow
		err	chan error
	}{
		{"without depedency", workflow.Workflow{
			Tasks: []workflow.Task{
				{
					Name:		"Name",
					TaskReference:	"TaskReference",
				},
			},
		}, make(chan error)},
		{"depedency", workflow.Workflow{
			Tasks: []workflow.Task{
				{
					Name:		"Name",
					TaskReference:	"TaskReference",
				},
				{
					Name:		"Dependecy",
					TaskReference:	"TaskDepedency",
					Dependency:	[]string{"Name"},
				},
			},
		}, make(chan error)},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s Steps

			s = steps(tt.wf, tt.err)
			fmt.Println(s)
		})
	}
}
