package bu

import (
  "os"
	"os/exec"
	"sync"
	"time"
  "github.com/aliafshar/toylog"
)

type targetQueue struct {
	items []*target
	done  map[string]bool
	sync.Mutex
}

func (q *targetQueue) build(t *target) {
	q.items = append(q.items, t)
	for _, u := range t.deps {
		q.build(u)
	}
}

func newTargetQueue(t *target) *targetQueue {
	q := &targetQueue{items: []*target{}, done: make(map[string]bool)}
	q.build(t)
	return q
}

func (q *targetQueue) peek() *target {
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

func (q *targetQueue) pop() *target {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	t := q.items[0]
	q.items = q.items[1:]
	return t
}

func (q *targetQueue) markDone(t *target) {
	q.Lock()
	defer q.Unlock()
	q.done[t.name] = true
}

func (q *targetQueue) hasDone(ts ...*target) bool {
	q.Lock()
	defer q.Unlock()
	for _, t := range ts {
		if !q.done[t.name] {
			return false
		}
	}
	return true
}

func (q *targetQueue) canRun(t *target) bool {
	return q.hasDone(t.deps...)
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

func (w *worker) runTarget(t *target) {
  toylog.Infof("< %q. [worker:%v]", t.body, w.id)
	cmd := runners[t.typ]
	c := exec.Command(cmd, "-c", t.body)
	o, err := c.CombinedOutput()
	if err != nil {
		toylog.Errorln("Run error.", err)
	}
  toylog.Infof("> %s", o)
}

type runner struct {
	out    <-chan []byte
	status <-chan int
}

func Newrunner() *runner {
	return &runner{out: make(chan []byte), status: make(chan int)}
}

func Run(s *script, t *target) {
	q := newTargetQueue(t)
  for _, setvar := range s.setvars {
    os.Setenv(setvar.key, setvar.value)
  }
	p := &pool{Size: 4}
	p.start(q)
}
