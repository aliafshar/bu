package bu

import (
	"os"
	"sync"
	"time"

	"github.com/aliafshar/toylog"
)

type targetQueue struct {
	items []target
	done  map[string]bool
	sync.Mutex
}

func (q *targetQueue) build(s *script, t target) {
	q.items = append(q.items, t)
	for _, d := range t.Deps() {
		u := d.runnable(s)
		if u != nil {
			q.build(s, u)
		}
	}
}

func newTargetQueue(s *script, t target) *targetQueue {
	q := &targetQueue{items: []target{}, done: make(map[string]bool)}
	q.build(s, t)
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

func (q *targetQueue) hasDone(ts ...string) bool {
	q.Lock()
	defer q.Unlock()
	for _, t := range ts {
		if !q.done[t] {
			return false
		}
	}
	return true
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
		if w.canRun(t) {
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

func (w *worker) runTarget(t target) result {
	toylog.Infof("> [%v] %v", t.Name(), t.Desc())
  res := t.Run()
  if !res.Success() {
		toylog.Errorf("< [%v] fail %v", t.Name(), res.Desc())
    return res
  }
	toylog.Infof("< [%v] done %v", t.Name(), res.Desc())
  return res
}

func (w *worker) canRun(t target) bool {
	for _, d := range t.Deps() {
		if !d.isDone(w) {
			return false
		}
	}
	return true
}

func Run(s *script, t target) {
	q := newTargetQueue(s, t)
	for _, setvar := range s.setvars {
		os.Setenv(setvar.key, setvar.value)
	}
	p := &pool{Size: 4}
	p.start(q)
}
