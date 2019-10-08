package server

import (
	"sync"
	"errors"
	"fmt"

	"github.com/lucasmbaia/kraken/workflow"
	"github.com/lucasmbaia/kraken/config"
	"github.com/lucasmbaia/kraken/proto"
	"github.com/lucasmbaia/kraken/core"
	"golang.org/x/net/context"
)

type KrakenServer struct {
	sync.RWMutex

	core	core.Core
	wf	[]workflow.Workflow
	tasks	map[string]workflow.TaskStatus
}

type KrakenServerConfig struct {
	Clients	[]interface{}
}

func NewKrakenServer(cfg KrakenServerConfig) (k *KrakenServer, err error) {
	k = &KrakenServer{tasks: make(map[string]workflow.TaskStatus)}

	if k.core, err = core.NewCore(core.CoreConfig{
		Clients:	cfg.Clients,
	}); err != nil {
		return
	}

	if k.wf, err = workflow.ReadWorkflows(config.EnvConfig.Workflow); err != nil {
		return
	}

	return
}

func (k *KrakenServer) Workflow(ctx context.Context, t *orchestrator.Task) (r *orchestrator.Response, err error) {
	var (
		wf	workflow.Workflow
		exists	bool
	)

	for _, w := range k.wf {
		if w.Name == t.Name && w.Version == t.Version {
			wf = w
			exists = true

			break
		}
	}

	if !exists {
		err = errors.New("Invalid workflow Name")
		return
	}

	wf.Body = t.Parameters

	k.Lock()
	k.tasks[t.Name] = workflow.TaskStatus{TotalSteps: int32(len(wf.Tasks))}
	k.Unlock()

	k.workflow(ctx, wf, t.Name)
	return
}

func (k *KrakenServer) workflow(ctx context.Context, wf workflow.Workflow, id string) {
	var (
		ts	= make(chan workflow.TaskStatus)
		results	core.Results
		err	error
	)

	go func() {
		for {
			select{
			case t := <-ts:
				fmt.Println(t)
			}
		}
	}()

	if results, err = k.core.RunWorkflow(ctx, wf, ts); err != nil {
		return
	}

	fmt.Println(results)
	return
}
