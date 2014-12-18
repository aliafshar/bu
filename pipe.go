package bu

import (
	"fmt"
	"github.com/aliafshar/toylog"
	"io"
	"os/exec"
)

type pipe struct {
	desc string
	cmds []*exec.Cmd
}

func (p *pipe) work(cmd *exec.Cmd, out chan *result) {
	out <- &result{err: cmd.Run()}
}

func (p *pipe) connect(final io.Writer) {
	if len(p.cmds) == 0 {
		return
	}
	last := len(p.cmds) - 1
	for i := 0; i < last; i++ {
		out, err := p.cmds[i].StdoutPipe()
		if err != nil {
			toylog.Errorln("unable to make a pipe", err)
		}
		p.cmds[i+1].Stdin = out
	}
	p.cmds[last].Stdout = final
}

func (p *pipe) run() *result {
	out := make(chan *result)
	for _, c := range p.cmds {
		go p.work(c, out)
	}
	var rs []*result
	for _ = range p.cmds {
		rs = append(rs, <-out)
	}
	return combinedResult(rs)
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
	fst := t.cmd(r)
	fst.Stdin = t.redirect.in()
	p := &pipe{cmds: []*exec.Cmd{fst}}
	p.desc = descTarget(t)
	for _, d := range t.pipe {
		s := d.resolve(r)
		p.cmds = append(p.cmds, s.cmd(r))
		p.desc = p.desc + " | " + descTarget(s)
	}
	p.connect(t.redirect.out())
	return p
}
