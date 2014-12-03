package bu

import (
  "fmt"
  "io"
  "os"
  "os/exec"
  "strings"
  "sync"
)

func work(wg *sync.WaitGroup, cmd *exec.Cmd) {
  cmd.Run()
  cmd.Stdout.(io.Closer).Close()
  wg.Done()
}

func connect(cs []*exec.Cmd, final io.Writer) {
  if len(cs) == 0 {
    return
  }
  last := len(cs) - 1
  for i := 0; i < last; i++ {
    cs[i + 1].Stdin, cs[i].Stdout = io.Pipe()
  }
  cs[last].Stdout = final
}

func run(cs []*exec.Cmd) {
  wg := &sync.WaitGroup{}
  for _, c := range cs {
    wg.Add(1)
    go work(wg, c)
  }
  wg.Wait()
}

func pipe(final io.Writer, cs ...*exec.Cmd) {
  connect(cs, final)
  run(cs)
}

func main() {
  pipe(
    os.Stdout,
    exec.Command("ls"),
    exec.Command("wc", "-l"),
    exec.Command("wcalc", "-h"),
  )
}

type pipeResult struct {
  err error
  rs []result
}

func (r *pipeResult) Err() error {
  return r.err
}

func (r *pipeResult) Success() bool {
  return r.err == nil
}

func (r *pipeResult) Desc() string {
  return fmt.Sprintf("total: %v\n", len(r.rs))
}

type pipeTarget struct {
  name string
  deps []dependency
}

func (t *pipeTarget) Name() string {
  return t.name
}

func (t *pipeTarget) Deps() []dependency {
  return nil
}

func (t *pipeTarget) Run(ctx *runContext) result {
  ctxs := []*runContext{}
  last := len(t.deps) - 1
  for i := 0; i <= last; i++ {
    ctxs = append(ctxs, &runContext{script: ctx.script, worker: ctx.worker})
  }
  for i := 0; i < last; i++ {
    pr, pw := io.Pipe()
    ctxs[i + 1].in = pr
    ctxs[i].out = pw
  }
  ctxs[last].out = ctx.out
  wg := &sync.WaitGroup{}
  r := &pipeResult{}
  for i, d := range t.deps {
    wg.Add(1)
    go func(cx *runContext, d dependency) {
      t := d.runnable(cx.script)
      sr := t.Run(cx)
      r.rs = append(r.rs, sr)
      cx.out.(io.Closer).Close()
      wg.Done()
    }(ctxs[i], d)
  }
  wg.Wait()
  return r
}

func (t *pipeTarget) Desc() string {
  ns := []string{}
  for _, d := range t.deps {
    ns = append(ns, d.(*targetDependency).name)
  }
  return fmt.Sprintf("pipe: %v", strings.Join(ns, " | "))
}


