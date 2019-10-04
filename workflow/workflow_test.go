package workflow

import (
	"testing"
)

func Test_ReadWorkflows(t *testing.T) {
	var err error

	var tests = []struct{
		name	string
		c	Config
	}{
		{"success", Config{Path: "./"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if _, err = ReadWorkflows(tt.c); err != nil {
				if err.Error() != tt.name {
					t.Errorf(err.Error())
				}
			}
		})
	}
}
