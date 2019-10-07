package workflow

import (
	"reflect"
)

type Task struct {
	Name		string		`json:",omitempty"`
	Description	string		`json:",omitempty"`
	TaskReference	string		`json:",omitempty"`
	Dependency	[]string	`json:",omitempty"`
	Timeout		int		`json:",omitempty"`
	Type		string		`json:",omitempty"`
	Drive		string		`json:",omitempty"`
	Retry		int		`json:",omitempty"`
	RetryDelay	int		`json:",omitempty"`
	Rollback	Rollback	`json:",omitempty"`
	GrpcService	GrpcService	`json:",omitempty"`
}

type GrpcService struct {
	Kind	string
}

type Rollback struct {
	Name		string	`json:",omitempty"`
	Description	string	`json:",omitempty"`
}

type Sign struct {
	FN	reflect.Value
	FT	reflect.Type
}

type Clients map[string]Sign

//var TaskClients = make(Clients)

func RegisterClients(clients []interface{}) (tc Clients, err error) {
	tc = make(Clients)

	for _, c := range clients {
		var fn = reflect.ValueOf(c)
		var ft = fn.Type()

		for i := 0; i < ft.NumMethod(); i++ {
			tc[ft.Method(i).Name] = Sign{
				FN:	fn,
				FT:	fn.MethodByName(ft.Method(i).Name).Type(),
				//FT: ft.Method(i).Func.Type(),
				//FT: ft.Method(i).Type,
			}
		}
	}

	return
}
