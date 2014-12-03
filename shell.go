package bu

import (
	"fmt"
  "io"
	"github.com/aliafshar/toylog"
	"os"
	"os/exec"
)

type shellTarget struct {
	name      string
	body      string
	depsNames []string
	deps      []dependency
  pipe      []dependency
	shell     shell
	args      []string
	outfile   string
	infile    string
  toClose   []io.Closer
}

func (t *shellTarget) Name() string {
	return t.name
}

func (t *shellTarget) Deps() []dependency {
	return t.deps
}

func (t *shellTarget) Run(ctx *runContext) result {
  t.ctx(ctx)
	return t.shell.execute(ctx, t)
}

func (t *shellTarget) Desc() string {
	return fmt.Sprintf("!%v %q", t.shell.desc(), t.body)
}

func (t *shellTarget) ctx(ctx *runContext) error {
	if t.outfile != "" {
		f, err := os.Create(t.outfile)
		if err != nil {
			return err
		}
		ctx.out = io.MultiWriter(f, ctx.out)
    t.toClose = append(t.toClose, f)
	}
  if t.infile != "" {
		f, err := os.Open(t.infile)
		if err != nil {
			return err
		}
		defer f.Close()
		ctx.in = f
	}
  return nil
}

type shell interface {
	execute(*runContext, *shellTarget) *shellResult
	desc() string
}

type shlike struct {
	Cmd string
}

func (sh *shlike) desc() string {
	return sh.Cmd
}

var shells = map[string]*shlike{
	"sh": &shlike{Cmd: "bash"},
	"py": &shlike{Cmd: "python"},
}

func (sh *shlike) execute(ctx *runContext, t *shellTarget) *shellResult {
	if t.body == "" {
		toylog.Errorf("< [%v] nothing to run", t.Name())
		return nil
	}
	args := append([]string{"-c", t.body, t.Name()}, t.args...)
	cmd := exec.Command(sh.Cmd, args...)
  if ctx.in != nil {
    cmd.Stdin = ctx.in
  }
  cmd.Stdout = ctx.out
	err := cmd.Run()
  cmd.Stdout.(io.Closer).Close()
	if err != nil {
		return newErrorShellResult(err)
	}
	return &shellResult{state: cmd.ProcessState}
}

type shellResult struct {
	err   error
	state *os.ProcessState
}

func newErrorShellResult(err error) *shellResult {
	return &shellResult{err: err}
}

func (r *shellResult) Err() error {
	return r.err
}

func (r *shellResult) Success() bool {
	return r.Err() == nil && r.state.Success()
}

func (r *shellResult) Desc() string {
	var v string
	if r.Success() {
		v = "0"
	} else {
		v = fmt.Sprintf("%v", r.Err())
	}
	return v
}
