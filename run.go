package bu

import (
	"fmt"
	"os"
  "os/exec"
	"sync"
	"time"

	"github.com/aliafshar/toylog"
	"github.com/mgutz/ansi"
	"golang.org/x/crypto/ssh/terminal"
)

const (
	feedbackChar = "●"
	launchChar   = "▶"
)

type queue struct {
	items []*target
	sync.Mutex
}

func (q *queue) reset() {
  q.items = nil
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

func feedback(char, color string) string {
	if terminal.IsTerminal(int(os.Stderr.Fd())) {
		return ansi.Color(char, color+"+b")
	} else {
		return feedbackChar + "(" + color + ")"
	}
}


func (w *worker) run(t *target) *result {
	p := newPipe(w.rt, t)
	toylog.Infof("%v [%v] %v", feedback(launchChar, "cyan"), targetDesc(t, w.rt), p.desc)
	res := p.run()
	if !res.success() {
		toylog.Errorf("%v %v [%v]",
			feedback(feedbackChar, "red"),
			res.desc,
			targetDesc(t, w.rt),
		)
		return res
	}
	toylog.Infof("%v %v [%v]",
		feedback(feedbackChar, "green"),
		res.desc,
		targetDesc(t, w.rt),
	)
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

func (h *history) reset() {
  h.log = make(map[string]bool)
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
  running map[int]*exec.Cmd
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
  for _, p := range t.pipe {
    u := p.resolve(r)
    for _, d := range u.deps {
      v := d.resolve(r)
      if v != nil {
        r.build(v)
      }
    }
  }
}

func (r *runtime) reset() {
  r.running = make(map[int]*exec.Cmd)
  r.queue.reset()
  r.history.reset()
}

func (r *runtime) stop() {
	for _, cmd := range r.running {
		cmd.Process.Kill()
	}
	for _, cmd := range r.running {
		cmd.Wait()
	}
  r.reset()
  r.wait.Wait()
}

func (r *runtime) run(t *target) {
  if t.watch != "" {
    r.runWatch(t)
    return
  }
  r.runOnce(t)
}

func (r *runtime) runOnce(t *target) {
  r.reset()
  r.build(t)
  r.start()
}

func (r *runtime) pollRestart(out chan bool) {
  for {
    select {
    case _ = <-out:
      r.stop()
    }
  }
}

func (r *runtime) runWatch(t *target) {
  out := make(chan bool)
  wt := &watcher{path: t.watch, out: out}
  toylog.Debugf("running a watch %+v", wt)
  go wt.watch()
  go r.pollRestart(out)
  for {
    r.runOnce(t)
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
    running: make(map[int]*exec.Cmd),
	}
}
