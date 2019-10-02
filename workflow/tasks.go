package workflow

type Task struct {
	Name		string		`json:",omitempty"`
	Description	string		`json:",omitempty"`
	TaskReference	string		`json:",omitempty"`
	Dependency	[]string	`json:",omitempty"`
	Timeout		string		`json:",omitempty"`
	Type		string		`json:",omitempty"`
	Drive		string		`json:",omitempty"`
	Retry		int		`json:",omitempty"`
	RetryDelay	string		`json:",omitempty"`
	Rollback	Rollback	`json:",omitempty"`
}

type Rollback struct {
	Name		string	`json:",omitempty"`
	Description	string	`json:",omitempty"`
}
