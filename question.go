package bu

import (
	"fmt"
	"github.com/aliafshar/weezard"
	"os"
)

type questionTarget struct {
	name  string
	dflt  string
	usage string
}

func (t *questionTarget) Name() string {
	return t.name
}

func (t *questionTarget) Deps() []dependency {
	return nil
}

func (t *questionTarget) Run() result {
	q := &weezard.Question{Usage: t.usage, Default: t.dflt}
	v, err := weezard.AskQuestion(q)
	if err != nil {
		return &questionResult{err: err}
	}
	os.Setenv(t.Name(), v)
	return &questionResult{key: t.Name(), value: v}
}

func (t *questionTarget) Desc() string {
	return "question"
}

type questionResult struct {
	err   error
	key   string
	value string
}

func (r *questionResult) Success() bool {
	return r.err == nil
}

func (r *questionResult) Desc() string {
	if r.Success() {
		return fmt.Sprintf("$%v=%q", r.key, r.value)
	} else {
		return fmt.Sprintf("%v", r.Err())
	}
}

func (r *questionResult) Err() error {
	return r.err
}
