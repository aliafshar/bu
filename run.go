package bu

import (
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/aliafshar/toylog"
)

type queue struct {
	items []*target
	sync.Mutex
}

func (q *queue) peek() *target {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	return q.items[0]
}

func (q *queue) rotate() {
	q.Lock()
	defer q.Unlock()
	if len(q.items) > 1 {
		q.items = append(q.items[1:], q.items[0])
	}
}

func (q *queue) pop() *target {
	q.Lock()
	defer q.Unlock()
	if len(q.items) == 0 {
		return nil
	}
	t := q.items[0]
	q.items = q.items[1:]
	return t
}

func (q *queue) push(t *target) {
	q.Lock()
	defer q.Unlock()
	q.items = append(q.items, t)
}

type pool struct {
	size    int
	workers []*worker
}

type worker struct {
	id int
	rt *runtime
}

func (w *worker) start() {
	for {
		t := w.rt.queue.pop()
		if t == nil {
			break
		}
		if w.can(t) {
			w.run(t)
			w.rt.history.do(t.name)
		} else {
			w.rt.queue.push(t)
			time.Sleep(100 * time.Millisecond)
		}
	}
	w.rt.wait.Done()
}

func targetDesc(t *target, r *runtime) string {
	return fmt.Sprintf("%v:%v", r.script.filename, t.name)
}

func (w *worker) run(t *target) *result {
	p := newPipe(w.rt, t)
	toylog.Infof("> [%v] %v", targetDesc(t, w.rt), p.desc)
	res := p.run()
	if !res.success() {
		toylog.Errorf("< [%v] fail %v", targetDesc(t, w.rt), res.desc)
		return res
	}
	toylog.Infof("< [%v] done %v", targetDesc(t, w.rt), res.desc)
	return res
}

func (w *worker) can(t *target) bool {
	for _, d := range t.deps {
		if !d.can(w.rt) {
			return false
		}
	}
	return true
}

type history struct {
	log map[string]bool
	sync.Mutex
}

func (h *history) do(key string) {
	h.Lock()
	defer h.Unlock()
	h.log[key] = true
}

func (h *history) done(keys ...string) bool {
	h.Lock()
	defer h.Unlock()
	for _, t := range keys {
		if !h.log[t] {
			return false
		}
	}
	return true
}

type runtime struct {
	script  *script
	pool    *pool
	queue   *queue
	history *history
	wait    *sync.WaitGroup
	argv    []string
	env     []string
}

func (r *runtime) start() {
	for i := 0; i < r.pool.size; i++ {
		w := &worker{id: i, rt: r}
		r.pool.workers = append(r.pool.workers, w)
		r.wait.Add(1)
		go w.start()
	}
	r.wait.Wait()
}

func (r *runtime) build(t *target) {
	r.queue.items = append(r.queue.items, t)
	for _, d := range t.deps {
		u := d.resolve(r)
		if u != nil {
			r.build(u)
		}
	}
}

func newRuntime(script *script) *runtime {
	// TODO share an environment
	for _, e := range script.setvars {
		os.Setenv(e.key, e.value)
	}
	var env []string
	copy(os.Environ(), env)
	return &runtime{
		script:  script,
		pool:    &pool{size: 4},
		history: &history{log: make(map[string]bool)},
		wait:    &sync.WaitGroup{},
		queue:   &queue{},
		env:     env,
	}
}
