package core

import (
	"encoding/json"
	"reflect"
	"context"
	"errors"
	"time"
	"sync"
	"fmt"

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

func (c *Core) RunWorkflow(ctx context.Context, wf workflow.Workflow) {
	var (
		s		Steps
		errc		= make(chan error, len(wf.Tasks))
		results		= make(Results, len(wf.Tasks))
		err		error
		size		= len(wf.Tasks)
		totalTasks	= len(wf.Tasks)
	)

	s = steps(wf, errc)

	for rt := range s.ReadyTasks {
		go func(t Tasks) {
			var (
				fn	reflect.Value
				ft	reflect.Type
				args	[]reflect.Value
				arg	reflect.Value
				output	[]reflect.Value
				body	[]byte
				ok	bool
				ct	context.Context
				cancel	context.CancelFunc
				done	= make(chan struct{})
			)

			defer s.Mutex.Unlock()

			c.Lock()
			fn = c.clients[t.Task.TaskReference].FN
			ft = c.clients[t.Task.TaskReference].FT
			c.Unlock()

			if ft.NumIn() > 0 {
				for i := 0; i < ft.NumIn(); i++ {
					if arg, err = c.setParameters(ctx, wf.Body, results, t.Task.Dependency, ft.In(i)); err != nil {
						break
					}

					args = append(args, arg)
				}
			}

			if t.Task.Timeout > 0 {
				t.Task.Timeout = defaultDeadlineContext
			}

			ct, cancel = context.WithTimeout(context.Background(), time.Duration(t.Task.Timeout) * time.Millisecond)
			defer cancel()

			go func() {
				output = fn.MethodByName(t.Task.TaskReference).Call(args)
				done <- struct{}{}
			}()

			select {
			case _ = <-done:
			case <-ct.Done():
				err = errors.New(fmt.Sprintf("Timeout reached to task name \"%s\"", t.Task.Name))
			}

			s.Mutex.Lock()
			if err == nil {
				for out := range output {
					if ft.Out(out).Kind() != reflect.Interface {
						if body, err = json.Marshal(output[out].Interface()); err != nil {
							break
						}

						results[t.Task.Name] = body
					} else {
						if err, ok = output[out].Interface().(error); ok {
							break
						}
					}
				}

				if err != nil {
					for i := 0; i < size; i++ {
						errc <- err
					}

					closeTasks(s)
					return
				}

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
			} else {
				for i := 0; i < size; i++ {
					errc <- err
				}

				closeTasks(s)
				return
			}
		}(rt)
	}

	for _, key := range results {
		fmt.Println(string(key))
	}
}

func (c *Core) setParameters(ctx context.Context, body []byte, results Results, dep []string, ft reflect.Type) (arg reflect.Value, err error) {
	switch ft.Kind() {
	case reflect.Ptr:
		arg = reflect.New(ft.Elem())

		if err = json.Unmarshal(body, arg.Interface()); err != nil {
			return
		}

		for _, d := range dep {
			if _, ok := results[d]; ok {
				if err = json.Unmarshal(results[d], arg.Interface()); err != nil {
					return
				}
			}
		}
	case reflect.Struct:
		arg = reflect.New(ft).Elem()

		if err = json.Unmarshal(body, &arg); err != nil {
			return
		}

		for _, d := range dep {
			if _, ok := results[d]; ok {
				if err = json.Unmarshal(results[d], &arg); err != nil {
					return
				}
			}
		}
	case reflect.Interface:
		arg = reflect.New(ft).Elem()
		arg.Set(reflect.ValueOf(ctx))
	}

	return
}
