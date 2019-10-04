package etcd

import (
	"testing"
	"fmt"
)

type Teste struct {
	Name		string	`json:",omitempty"`
	Sobrenome	string	`json:",omitempty"`
}

func Test_NewClient(t *testing.T) {
	if _, err := NewClient(Config{
		Endpoints:	[]string{"127.0.0.1:2379"},
		Timeout:	10,
	}); err != nil {
		t.Fatal(err)
	}
}

func Test_Get(t *testing.T) {
	c, err := NewClient(Config{
		Endpoints:	[]string{"http://127.0.0.1:2379"},
		Timeout:	10,
	})
	if err != nil {
		t.Fatal(err)
	}

	var teste Teste

	if err := c.Get("lucas", &teste); err != nil {
		t.Fatal(err)
	}

	fmt.Println(teste)
}

func Test_Watch(t *testing.T) {
	c, err := NewClient(Config{
		Endpoints:	[]string{"http://127.0.0.1:2379"},
		Timeout:	10,
	})
	if err != nil {
		t.Fatal(err)
	}

	var values = make(chan Response)

	go func() {
		for {
			select {
			case v := <-values:
				fmt.Println(v)
			}
		}
	}()

	c.Watch("lucas", values)
}
