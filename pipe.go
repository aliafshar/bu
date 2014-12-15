package bu

import (
	"io"
	"os/exec"
	"sync"

  _ "github.com/aliafshar/toylog"
)

type pipe struct {
	cmds []*exec.Cmd
	sync.WaitGroup
}

func (p *pipe) work(cmd *exec.Cmd) {
	cmd.Run()
	p.Done()
}

func (p *pipe) connect(final io.Writer) {
	if len(p.cmds) == 0 {
		return
	}
	last := len(p.cmds) - 1
	for i := 0; i < last; i++ {
		p.cmds[i+1].Stdin, p.cmds[i].Stdout = io.Pipe()
	}
	p.cmds[last].Stdout = final
}

func (p *pipe) run() *result {
	for _, c := range p.cmds {
		p.Add(1)
		go p.work(c)
	}
	p.Wait()
	return &result{}
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
