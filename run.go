package bu

import (
	"github.com/aliafshar/toylog"
	"github.com/aliafshar/weezard"
	"os"
	"os/exec"
	"sync"
	"time"
)

type targetQueue struct {
	items []target
	done  map[string]bool
	sync.Mutex
}

func (q *targetQueue) build(t target) {
	q.items = append(q.items, t)
	for _, u := range t.Deps() {
		q.build(u)
	}
}

func newTargetQueue(t target) *targetQueue {
	q := &targetQueue{items: []target{}, done: make(map[string]bool)}
	q.build(t)
	return q
}

func (q *targetQueue) peek() target {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	return q.items[0]
}

func (q *targetQueue) rotate() {
	q.Lock()
	defer q.Unlock()
	if len(q.items) > 1 {
		q.items = append(q.items[1:], q.items[0])
	}
}

func (q *targetQueue) pop() target {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	t := q.items[0]
	q.items = q.items[1:]
	return t
}

func (q *targetQueue) markDone(t target) {
	q.Lock()
	defer q.Unlock()
	q.done[t.Name()] = true
}

func (q *targetQueue) hasDone(ts ...target) bool {
	q.Lock()
	defer q.Unlock()
	for _, t := range ts {
		if !q.done[t.Name()] {
			return false
		}
	}
	return true
}

func (q *targetQueue) canRun(t target) bool {
	return q.hasDone(t.Deps()...)
}

type pool struct {
	Size int
}

func (p *pool) start(q *targetQueue) {
	wg := &sync.WaitGroup{}
	for i := 0; i < p.Size; i++ {
		w := &worker{wg: wg, id: i, q: q}
		wg.Add(1)
		go w.run()
	}
	wg.Wait()
}

type worker struct {
	wg *sync.WaitGroup
	id int
	q  *targetQueue
}

func (w *worker) run() {
	for {
		t := w.q.peek()
		if t == nil {
			break
		}
		if w.q.canRun(t) {
			w.q.pop()
			w.runTarget(t)
			w.q.markDone(t)
		} else {
			w.q.rotate()
			time.Sleep(100 * time.Millisecond)
		}
	}
	w.wg.Done()
}

var runners = map[string]string{
	"sh": "bash",
	"py": "python",
	"":   "bash",
}

func (w *worker) runTarget(t target) {
	t.Run()
}

func (t *shellTarget) Run() {
	toylog.Infof("> [%v] %v:%#v", t.Name(), t.typ, t.body)
	shell := runners[t.typ]
	cmd := exec.Command(shell, "-c", t.body)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		toylog.Errorf("< [%v] failure, %v", t.name, err)
	} else {
		toylog.Infof("< [%v] success", t.name)
	}
}

func (t *questionTarget) Run() {
	toylog.Infof("> [%v] question: %#v", t.Name(), t.usage)
	q := &weezard.Question{Usage: t.usage, Default: t.dflt}
	v, err := weezard.AskQuestion(q)
	if err != nil {
		toylog.Errorf("< [%v] failure, %v", t.Name(), err)
	}
	os.Setenv(t.Name(), v)
	toylog.Infof("< [%v] success $%v=%q", t.Name(), t.Name(), v)
}

func Run(s *script, t target) {
	q := newTargetQueue(t)
	for _, setvar := range s.setvars {
    toylog.Errorf("%q %q", setvar.key, setvar.value)
		os.Setenv(setvar.key, setvar.value)
	}
	p := &pool{Size: 4}
	p.start(q)
}
