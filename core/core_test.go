package core

import (
	"github.com/lucasmbaia/kraken/workflow"
	"testing"
	"context"
	"time"
	"fmt"
	//"errors"
)

type MockCore struct {}

type MockSum struct {
	A	int	`json:",omitempty"`
	B	int	`json:",omitempty"`
	C	int	`json:",omitempty"`
}

type MockSplit struct {
	A	int	`json:",omitempty"`
	B	int	`json:",omitempty"`
	C	int	`json:",omitempty"`
}

func (m MockCore) Sum(ctx context.Context, msint *MockSum) (msout *MockSum, err error) {
	msout = &MockSum{C: msint.A + msint.B}
	fmt.Println("SUM")
	//err = errors.New("TESTE")
	time.Sleep(1000 * time.Millisecond)
	return
}

func (m MockCore) Split(ctx context.Context, msint *MockSplit) (msout *MockSplit, err error) {
	msout = &MockSplit{C: (msint.A + msint.B) / msint.C}
	return
}

func Test_RunWorkflow(t *testing.T) {
	var err error
	var mc MockCore = MockCore{}

	var tests = []struct{
		name	string
		wf	workflow.Workflow
		clients	[]interface{}
	}{
		{"without depedency", workflow.Workflow{
			Body:	[]byte(`{"A": 1, "B": 2}`),
			Tasks:	[]workflow.Task{
				{
					Name:		"MockSum",
					TaskReference:	"Sum",
					Timeout:	50,
					Retry:		3,
					RetryDelay:	1500,
				},
				{
					Name:		"MockSplit",
					TaskReference:	"Split",
					Dependency:	[]string{"MockSum"},
					Timeout:	50,
				},
			},
		}, []interface{}{mc}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var c Core

			if c, err = NewCore(CoreConfig{
				Clients:	tt.clients,
			}); err != nil {
				t.Fatal(err)
			}

			c.RunWorkflow(context.Background(), tt.wf)
		})
	}
}
