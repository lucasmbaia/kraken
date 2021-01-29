package server

import (
	"github.com/lucasmbaia/kraken/workflow"
	"github.com/lucasmbaia/kraken/proto"
	"testing"
	"golang.org/x/net/context"
	//"fmt"
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
		wf	[]workflow.Workflow
		clients	[]interface{}
	}{
		{"without rollback", []workflow.Workflow{{
			Name:	"kraken",
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
		}},[]interface{}{mc}},
		/*{"with rollback", workflow.Workflow{
			Name:	"kraken",
			Tasks:	[]workflow.Task{
				{
					Name:		"MockSum",
					TaskReference:	"Sum",
					Timeout:	50,
					Retry:		3,
					RetryDelay:	1500,
					Rollback:	workflow.Rollback{
						Name:	"MockSumRollback",
						Step:	1
					},
				},
				{
					Name:		"MockSplit",
					TaskReference:	"Split",
					Dependency:	[]string{"MockSum"},
					Timeout:	50,
					Rollback:	workflow.Rollback{
						Name:	"MockSumRollback",
						Step:	1
					},
				},
			},
		},[]interface{}{mc}},*/
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var ks *KrakenServer

			if ks, err = NewKrakenServer(KrakenServerConfig{
				Clients:	tt.clients,
			}); err != nil {
				t.Fatal(err)
			}

			ks.wf = append(ks.wf, tt.wf...)
			ks.Workflow(context.Background(), &orchestrator.Task{
				Name: "kraken",
				Parameters:	[]byte(`{"A": 1, "B": 2}`),
			})
		})
	}
}
