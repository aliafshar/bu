package bu

import (
	"fmt"
	"github.com/aliafshar/toylog"
	"io"
	"os/exec"
)

type pipe struct {
	desc    string
  runtime *runtime
  target *target
}

func (p *pipe) work(cmd *exec.Cmd, out chan *result) {
  err := cmd.Start()
  if err == nil {
    p.runtime.running[cmd.Process.Pid] = cmd
    err = cmd.Wait()
    delete(p.runtime.running, cmd.Process.Pid)
  }
	out <- &result{err: err}
}

func (p *pipe) connect(cmds []*exec.Cmd, final io.Writer) {
	if len(cmds) == 0 {
		return
	}
	last := len(cmds) - 1
	for i := 0; i < last; i++ {
		out, err := cmds[i].StdoutPipe()
		if err != nil {
			toylog.Errorln("unable to make a pipe", err)
		}
		cmds[i+1].Stdin = out
	}
	cmds[last].Stdout = final
}

func (p *pipe) run() *result {
  cmds := p.build()
	out := make(chan *result)
	for _, c := range cmds {
		go p.work(c, out)
	}
	var rs []*result
	for _ = range cmds {
		rs = append(rs, <-out)
	}
	return combinedResult(rs)
}

func (p *pipe) build() []*exec.Cmd {
	fst := p.target.cmd(p.runtime)
	fst.Stdin = p.target.redirect.in()
  cmds := []*exec.Cmd{fst}
	for _, d := range p.target.pipe {
		s := d.resolve(p.runtime)
		cmds = append(cmds, s.cmd(p.runtime))
	}
	p.connect(cmds, p.target.redirect.out())
  return cmds
}


func descTarget(t *target) string {
	return fmt.Sprintf("%q", t.body)
}

func combinedResult(rs []*result) *result {
	r := &result{}
	for i, rr := range rs {
		if i > 0 {
			r.desc = r.desc + " | "
		}
		if rr.err != nil {
			r.err = rr.err
			r.desc = r.desc + fmt.Sprintf("%v", r.err)
		} else {
			r.desc = r.desc + "0"
		}
	}
	return r
}

func newPipe(r *runtime, t *target) *pipe {
  p := &pipe{runtime: r, target: t, desc: descTarget(t)}
	for _, d := range p.target.pipe {
		s := d.resolve(p.runtime)
		p.desc = p.desc + " | " + descTarget(s)
	}
  return p
}
