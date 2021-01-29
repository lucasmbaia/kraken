package core

/*import (
	"golang.org/x/net/context"
	"github.com/lucasmbaia/kraken/workflow"
)

func (c *Core) Rollback(ctx context.Context, wf workflow.Workflow, step int, ts chan<- workflow.TaskStatus, agr Results) (results Results, err error) {
	if step > len(wf.Tasks) {
		step = len(wf.Tasks)
	}

	step--
	if step < 0 {
		step = 0
	}

	wf.Tasks = wf.Tasks[step:]
	results, err = c.RunWorkflow(ctx, wf, ts, )

	return
}*/
