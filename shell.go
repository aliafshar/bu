package bu

import(
  "fmt"
  "os"
  "os/exec"
	"github.com/aliafshar/toylog"
)

type shellTarget struct {
	name      string
	body      string
	depsNames []string
	deps      []dependency
	shell     shell
	args      []string
  outfile   string
  infile string
}

func (t *shellTarget) Name() string {
	return t.name
}

func (t *shellTarget) Deps() []dependency {
	return t.deps
}

func (t *shellTarget) Run() result {
  return t.shell.execute(t)
}

func (t *shellTarget) Desc() string {
  return fmt.Sprintf("!%v %q", t.shell.desc(), t.body)
}

type shell interface {
	execute(t *shellTarget) *shellResult
  desc() string
}

type shlike struct{
  Cmd string
}

func (sh *shlike) desc() string {
  return sh.Cmd
}

var shells = map[string]*shlike {
    "sh": &shlike{Cmd: "bash"},
    "py": &shlike{Cmd: "python"},
}

func (sh *shlike) execute(t *shellTarget) *shellResult {
	if t.body == "" {
		toylog.Errorf("< [%v] nothing to run", t.Name())
    return nil
	}
	args := append([]string{"-c", t.body, t.Name()}, t.args...)
	cmd := exec.Command(sh.Cmd, args...)
  if t.outfile != "" {
    f, err := os.Create(t.outfile)
    if err != nil {
      return newErrorShellResult(err)
    }
    defer f.Close()
    cmd.Stdout = f
    cmd.Stderr = f
  } else {
	  cmd.Stdout = os.Stdout
	  cmd.Stderr = os.Stderr
  }
  if t.infile != "" {
    f, err := os.Open(t.infile)
    if err != nil {
      return newErrorShellResult(err)
    }
    defer f.Close()
    cmd.Stdin = f
  } else {
    cmd.Stdin = os.Stdin
  }
	err := cmd.Run()
	if err != nil {
    return newErrorShellResult(err)
  }
  return &shellResult{state: cmd.ProcessState}
}

type shellResult struct {
  err error
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
