package workflow

import (
	"testing"
	"fmt"
)

type Test struct {}

func (t *Test) Print(str string) {
	fmt.Println(str)
}

func Test_RegisterClients(t *testing.T) {
	var tt = &Test{}

	var i = []interface{}{tt}

	if tc, err := RegisterClients(i); err != nil {
		t.Fatal(err)
	} else {
		fmt.Println(tc)
	}
}
