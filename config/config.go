package config

import (
	"github.com/lucasmbaia/kraken/workflow"
)

var (
	EnvConfig	Config
)

type Config struct {
	Server	bool	`hcl:"server"`
	Client	bool	`hcl:"client"`
}

type Workflow struct {
	Path		string
	EtcdURL		string
	EtcdUsername	string
	EtcdPassword	string
}

type Singleton struct {
}
