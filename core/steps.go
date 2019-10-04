package core

import (
	"reflect"
	"sync"

	"github.com/lucasmbaia/kraken/workflow"
)

type Tasks struct {
	HasDependency	bool
	Task		workflow.Task
	FN		reflect.Value
}

type Steps struct {
	Mutex		*sync.RWMutex

	ReadyTasks	chan Tasks
	Dependents	map[string][]*sync.WaitGroup
	ReadyToCheck	[]string
	UncheckDeps	map[string]int
}

func steps(w workflow.Workflow, err chan error) (s Steps) {
	s.ReadyTasks = make(chan Tasks, len(w.Tasks))
	s.Dependents = make(map[string][]*sync.WaitGroup)
	s.UncheckDeps = make(map[string]int)
	s.Mutex = &sync.RWMutex{}

	for _, t := range w.Tasks {
		if len(t.Dependency) > 0 {
			var wg = new(sync.WaitGroup)
			wg.Add(len(t.Dependency))

			for _, name := range t.Dependency {
				s.Dependents[name] = append(s.Dependents[name], wg)
			}

			s.UncheckDeps[t.Name] = len(t.Dependency)

			go func(wg *sync.WaitGroup, task workflow.Task) {
				wg.Wait()

				s.Mutex.Lock()
				select {
				case _ = <-err:
					s.Mutex.Unlock()
					return
				default:
					s.Mutex.Unlock()
				}

				s.ReadyTasks <- Tasks{HasDependency: true, Task: task}
			}(wg, t)
		} else {
			s.ReadyToCheck = append(s.ReadyToCheck, t.Name)
			s.ReadyTasks <- Tasks{Task: t}
		}
	}

	return
}

func closeTasks(s Steps) {
	for _, d := range s.Dependents {
		for _, wg := range d {
			wg.Done()
		}
	}

	close(s.ReadyTasks)
}
