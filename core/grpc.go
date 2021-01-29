package core

import (
	"encoding/json"
	"reflect"
	"errors"
	"time"
	"fmt"

	"golang.org/x/net/context"
	"github.com/lucasmbaia/kraken/workflow"
)

func (c *Core) grpc(ctx context.Context, t workflow.Task, body []byte, results, agr Results) (response []byte, err error) {
	var (
		fn      reflect.Value
		ft      reflect.Type
		args    []reflect.Value
		arg     reflect.Value
		output  []reflect.Value
		ok      bool
		ct      context.Context
		cancel  context.CancelFunc
		done    = make(chan struct{})
	)

	c.Lock()
	fn = c.clients[t.TaskReference].FN
	ft = c.clients[t.TaskReference].FT
	c.Unlock()

	if ft.NumIn() > 0 {
		for i := 0; i < ft.NumIn(); i++ {
			if arg, err = c.setGrpcParameters(ctx, body, results, agr, t.Dependency, ft.In(i)); err != nil {
				break
			}

			args = append(args, arg)
		}
	}

	if t.Timeout == 0 {
		t.Timeout = defaultDeadlineContext
	}

	ct, cancel = context.WithTimeout(context.Background(), time.Duration(t.Timeout) * time.Millisecond)
	defer cancel()

	go func() {
		output = fn.MethodByName(t.TaskReference).Call(args)
		done <- struct{}{}
	}()

	select {
	case _ = <-done:
	case <-ct.Done():
		err = errors.New(fmt.Sprintf("Timeout reached to task name \"%s\"", t.Name))
		return
	}

	if err == nil {
		for out := range output {
			if ft.Out(out).Kind() != reflect.Interface {
				if response, err = json.Marshal(output[out].Interface()); err != nil {
					break
				}
			} else {
				if err, ok = output[out].Interface().(error); ok {
					break
				}
			}
		}
	}

	return
}

func (c *Core) setGrpcParameters(ctx context.Context, body []byte, results, agr Results, dep []string, ft reflect.Type) (arg reflect.Value, err error) {
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

		if agr != nil {
			for _, r := range agr {
				if err = json.Unmarshal(r, arg.Interface()); err != nil {
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

		if agr != nil {
			for _, r := range agr {
				if err = json.Unmarshal(r, &arg); err != nil {
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

