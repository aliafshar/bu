package bu

import (
	"io"
	"os/exec"

	"github.com/aliafshar/toylog"
)

type pipe struct {
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

func newPipe(r *runtime, t *target) *pipe {
	fst := t.cmd(r)
	fst.Stdin = t.redirect.in()
	p := &pipe{cmds: []*exec.Cmd{fst}}
	for _, d := range t.pipe {
		p.cmds = append(p.cmds, d.resolve(r).cmd(r))
	}
	p.connect(t.redirect.out())
	return p
}
