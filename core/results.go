package core

type Results map[string][]byte

type WError struct {
	Error error
	Task  string
}
