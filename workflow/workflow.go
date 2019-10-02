package workflow

import (
	"encoding/json"
	"context"
	"ioutil"
	"fmt"
	"os"
)

type Workflow struct {
	Name		string	`json:",omitempty"`
	Description	string	`json:",omitempty"`
	Version		int	`json:",omitempty"`
	Tasks		[]Task	`json:",omitempty"`
}

type WorklowManager struct {
	Cancel	context.CancelFunc
	Stop	chan struct{}
	Restart chan struct{}
}

type Congig struct {
	Path		string
	EtcdURL		string
	EtcdKeys	[]string
}

type WM map[string]*WorklowManager

func ReadWorkflows(c Config) (w []Workflow, err error) {
	if c.Path != "" {
		var (
			files	[]os.FileInfo
			body	[]byte
		)

		if c.Path[len(c.Path)-1:] != "/" {
			c.Path += "/"
		}

		if _, err = os.Stat(c.Path); os.IsNotExists(err) || err != nil {
			return
		}

		if files, err = ioutil.ReadDir(c.Path); err != nil {
			return
		}

		for _, file := range files {
			if path.Ext(file.Name()) == ".json" {
				var wf Workflow
				if body, err = ioutil.ReadFile(fmt.Sprintf("%s%s", c.Path, file.Name())); err != nil {
					return
				}

				if err = json.Unmarshal(body, &wf); err != nil {
					return
				}

				w = append(w, wf)
			}
		}
	}

	return
}
