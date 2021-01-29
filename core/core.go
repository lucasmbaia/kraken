package core

import (
	"time"
	"sync"

	"golang.org/x/net/context"
	"github.com/lucasmbaia/kraken/workflow"
)

const (
	defaultDeadlineContext = 50000000
)

type Core struct {
	sync.RWMutex

	clients	workflow.Clients
}

type CoreConfig struct {
	Clients	[]interface{}
}

func NewCore(cfg CoreConfig) (c Core, err error) {
	if c.clients, err = workflow.RegisterClients(cfg.Clients); err != nil {
		return
	}

	return
}

func (c *Core) RunWorkflow(ctx context.Context, wf workflow.Workflow, ts chan<- workflow.TaskStatus, agr Results) (results Results, we WError) {
	var (
		s		Steps
		errc		= make(chan error, len(wf.Tasks))
		size		= len(wf.Tasks)
		totalTasks	= len(wf.Tasks)
		err		error
	)

	results = make(Results, len(wf.Tasks))
	s = steps(wf, errc)

	for rt := range s.ReadyTasks {
		go func(t Tasks) {
			var (
				response	[]byte
				done		= make(chan struct{}, 1)
				errout		= make(chan struct{}, 1)
				call		= make(chan struct{}, 1)
				retry		= 1
				check		= make(chan struct{}, 1)
			)

			go func() {
				call <- struct{}{}
			}()

			if ts != nil {
				ts <- workflow.TaskStatus{Name: t.Task.Name}
			}

			go func() {
				for {
					select {
					case _ = <-call:
						if response, err = c.grpc(ctx, t.Task, wf.Body, results, agr); err != nil {
							check <- struct{}{}
						} else {
							done <- struct{}{}
						}
					case _ = <-check:
						if t.Task.Retry > 0 && retry < t.Task.Retry {
							retry++

							if t.Task.RetryDelay > 0 {
								time.Sleep(time.Duration(t.Task.RetryDelay) * time.Millisecond)
							}

							call <- struct{}{}
							break
						}

						errout <- struct{}{}
					}
				}
			}()

			select {
			case _ = <-done:
				results[t.Task.Name] = response
				err = nil

				s.Mutex.Lock()
				if wgs, ok := s.Dependents[t.Task.Name]; ok {
					for _, wg := range wgs {
						wg.Done()
					}

					delete(s.Dependents, t.Task.Name)
				}

				totalTasks--

				if totalTasks == 0 {
					close(s.ReadyTasks)
					close(errc)
				}
				s.Mutex.Unlock()
			case _ = <-errout:
				for i := 0; i < size; i++ {
					errc <- err
				}

				closeTasks(s)

				we = WError{
					Error:	err,
					Task:	t.Task.Name,
				}
			}
		}(rt)
	}

	return
}
